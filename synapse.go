package synapse

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/kolosys/synapse/eviction"
)

// EvictionPolicy is re-exported from the eviction package
type EvictionPolicy = eviction.EvictionPolicy

// Cache is a generic similarity-based cache with sharding
type Cache[K comparable, V any] struct {
	shards     []*Shard[K, V]
	similarity SimilarityFunc[K]
	threshold  float64
	options    *Options
}

// New creates a new cache with the given options
func New[K comparable, V any](opts ...Option) *Cache[K, V] {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	c := &Cache[K, V]{
		shards:    make([]*Shard[K, V], options.NumShards),
		threshold: options.SimilarityThreshold,
		options:   options,
	}

	// Initialize shards
	maxSizePerShard := options.MaxSize / options.NumShards
	if maxSizePerShard == 0 {
		maxSizePerShard = 1
	}

	for i := 0; i < options.NumShards; i++ {
		var policy eviction.EvictionPolicy
		if options.EvictionPolicy != nil {
			policy = options.EvictionPolicy
		}

		c.shards[i] = newShard[K, V](
			maxSizePerShard,
			c.similarity,
			c.threshold,
			options.TTL,
			policy,
			options.EnableStats,
		)
	}

	return c
}

// WithSimilarity sets the similarity function for the cache
func (c *Cache[K, V]) WithSimilarity(fn SimilarityFunc[K]) *Cache[K, V] {
	c.similarity = fn
	for _, shard := range c.shards {
		shard.similarity = fn
	}
	return c
}

// getShard returns the shard for a given key
func (c *Cache[K, V]) getShard(key K) *Shard[K, V] {
	h := fnv.New64a()
	// Use string representation of key for hashing
	// This is a simple approach; for production, you might want a more sophisticated method
	h.Write([]byte(keyToString(key)))
	hash := h.Sum64()
	return c.shards[int(hash%uint64(len(c.shards)))]
}

// Get retrieves a value by exact key match
func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, bool) {
	shard := c.getShard(key)
	return shard.get(ctx, key)
}

// Set stores a value
func (c *Cache[K, V]) Set(ctx context.Context, key K, value V) error {
	shard := c.getShard(key)
	return shard.set(ctx, key, value)
}

// GetSimilar finds the most similar key above the threshold
func (c *Cache[K, V]) GetSimilar(ctx context.Context, key K) (V, K, float64, bool) {
	// For similarity search, we need to search across all shards
	// In a production implementation, you might want to use LSH or other indexing

	var bestValue V
	var bestKey K
	bestScore := 0.0
	found := false

	for _, shard := range c.shards {
		v, k, score, ok := shard.getSimilar(ctx, key)
		if ok && score > bestScore {
			bestValue = v
			bestKey = k
			bestScore = score
			found = true
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			var zeroV V
			var zeroK K
			return zeroV, zeroK, 0, false
		default:
		}
	}

	return bestValue, bestKey, bestScore, found
}

// Delete removes a key from the cache
func (c *Cache[K, V]) Delete(ctx context.Context, key K) bool {
	shard := c.getShard(key)
	return shard.delete(ctx, key)
}

// Len returns the total number of entries in the cache
func (c *Cache[K, V]) Len() int {
	total := 0
	for _, shard := range c.shards {
		total += shard.len()
	}
	return total
}

// Stats returns aggregated statistics from all shards
// Returns zero values if stats are not enabled
func (c *Cache[K, V]) Stats() Stats {
	if !c.options.EnableStats {
		return Stats{}
	}

	var stats Stats
	for _, shard := range c.shards {
		if shard.stats != nil {
			shardStats := shard.stats.snapshot()
			stats.Hits += shardStats.Hits
			stats.Misses += shardStats.Misses
			stats.Sets += shardStats.Sets
			stats.Deletes += shardStats.Deletes
			stats.SimilarSearches += shardStats.SimilarSearches
			stats.SimilarHits += shardStats.SimilarHits
			stats.Evictions += shardStats.Evictions
			stats.Expired += shardStats.Expired
		}
	}
	return stats
}

// keyToString converts a key to a string for hashing
// This is a simple implementation; for production use, consider a more robust approach
func keyToString[K comparable](key K) string {
	// Use fmt.Sprintf for a simple string representation
	return fmt.Sprintf("%v", key)
}
