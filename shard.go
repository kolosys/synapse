package synapse

import (
	"context"
	"sync"
	"time"

	"github.com/kolosys/synapse/eviction"
)

// Shard represents a single shard of the cache
type Shard[K comparable, V any] struct {
	mu             sync.RWMutex
	data           map[K]*Entry[K, V]
	keys           []K // For similarity search iteration
	evictionPolicy eviction.EvictionPolicy
	maxSize        int
	similarity     SimilarityFunc[K]
	threshold      float64
	ttl            time.Duration
	stats          *shardStats
	enableStats    bool
}

// newShard creates a new cache shard
func newShard[K comparable, V any](maxSize int, similarity SimilarityFunc[K], threshold float64, ttl time.Duration, policy eviction.EvictionPolicy, enableStats bool) *Shard[K, V] {
	s := &Shard[K, V]{
		data:           make(map[K]*Entry[K, V]),
		keys:           make([]K, 0),
		evictionPolicy: policy,
		maxSize:        maxSize,
		similarity:     similarity,
		threshold:      threshold,
		ttl:            ttl,
		enableStats:    enableStats,
	}
	if enableStats {
		s.stats = newShardStats()
	}
	return s
}

// get retrieves a value by exact key match
func (s *Shard[K, V]) get(ctx context.Context, key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		var zero V
		return zero, false
	default:
	}

	namespace := GetNamespace(ctx)

	entry, ok := s.data[key]
	if !ok {
		if s.enableStats {
			s.stats.recordMiss()
		}
		var zero V
		return zero, false
	}

	// Check namespace match
	if namespace != "" && entry.Namespace != namespace {
		if s.enableStats {
			s.stats.recordMiss()
		}
		var zero V
		return zero, false
	}

	// Check expiration
	if entry.IsExpired() {
		if s.enableStats {
			s.stats.recordExpired()
			s.stats.recordMiss()
		}
		var zero V
		return zero, false
	}

	// Update access tracking
	entry.Touch()
	if s.evictionPolicy != nil {
		s.evictionPolicy.OnAccess(key)
	}

	if s.enableStats {
		s.stats.recordHit()
	}

	return entry.Value, true
}

// getSimilar finds the most similar key above the threshold
func (s *Shard[K, V]) getSimilar(ctx context.Context, key K) (V, K, float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		var zeroV V
		var zeroK K
		return zeroV, zeroK, 0, false
	default:
	}

	if s.enableStats {
		s.stats.recordSimilarSearch()
	}

	namespace := GetNamespace(ctx)

	var bestKey K
	var bestValue V
	bestScore := 0.0
	found := false

	for _, k := range s.keys {
		entry := s.data[k]

		// Check namespace match
		if namespace != "" && entry.Namespace != namespace {
			continue
		}

		// Check expiration
		if entry.IsExpired() {
			continue
		}

		// Check context cancellation periodically
		select {
		case <-ctx.Done():
			var zeroV V
			var zeroK K
			return zeroV, zeroK, 0, false
		default:
		}

		// Compute similarity
		if s.similarity != nil {
			score := s.similarity(key, k)
			if score >= s.threshold && score > bestScore {
				bestKey = k
				bestValue = entry.Value
				bestScore = score
				found = true
			}
		}
	}

	if found {
		// Update access tracking
		entry := s.data[bestKey]
		entry.Touch()
		if s.evictionPolicy != nil {
			s.evictionPolicy.OnAccess(bestKey)
		}
		if s.enableStats {
			s.stats.recordSimilarHit()
		}
	}

	return bestValue, bestKey, bestScore, found
}

// set stores a value
func (s *Shard[K, V]) set(ctx context.Context, key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	namespace := GetNamespace(ctx)

	// Check if key already exists
	if entry, ok := s.data[key]; ok {
		entry.Value = value
		entry.Touch()
		if s.evictionPolicy != nil {
			s.evictionPolicy.OnAccess(key)
		}
		if s.enableStats {
			s.stats.recordSet()
		}
		return nil
	}

	// Evict if necessary
	if s.maxSize > 0 && len(s.data) >= s.maxSize {
		if err := s.evict(); err != nil {
			return err
		}
	}

	// Create new entry
	entry := newEntry(key, value, s.ttl, namespace)
	s.data[key] = entry
	s.keys = append(s.keys, key)

	if s.evictionPolicy != nil {
		s.evictionPolicy.OnAdd(key, entry.AccessCount, entry.CreatedAt, entry.AccessedAt)
	}

	if s.enableStats {
		s.stats.recordSet()
	}

	return nil
}

// delete removes a key from the shard
func (s *Shard[K, V]) delete(ctx context.Context, key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return false
	default:
	}

	if _, ok := s.data[key]; !ok {
		return false
	}

	delete(s.data, key)

	// Remove from keys slice
	for i, k := range s.keys {
		if k == key {
			s.keys = append(s.keys[:i], s.keys[i+1:]...)
			break
		}
	}

	if s.evictionPolicy != nil {
		s.evictionPolicy.OnRemove(key)
	}

	if s.enableStats {
		s.stats.recordDelete()
	}

	return true
}

// evict removes an entry based on the eviction policy
func (s *Shard[K, V]) evict() error {
	if s.evictionPolicy == nil {
		// No eviction policy, just remove the first key
		if len(s.keys) > 0 {
			key := s.keys[0]
			delete(s.data, key)
			s.keys = s.keys[1:]
			if s.enableStats {
				s.stats.recordEviction()
			}
		}
		return nil
	}

	victim, ok := s.evictionPolicy.SelectVictim()
	if !ok {
		return nil
	}

	key, ok := victim.(K)
	if !ok {
		return nil
	}

	delete(s.data, key)

	// Remove from keys slice
	for i, k := range s.keys {
		if k == key {
			s.keys = append(s.keys[:i], s.keys[i+1:]...)
			break
		}
	}

	s.evictionPolicy.OnRemove(key)

	if s.enableStats {
		s.stats.recordEviction()
	}

	return nil
}

// len returns the number of entries in the shard
func (s *Shard[K, V]) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
