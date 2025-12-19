package eviction

import (
	"sync"
	"time"
)

// TTL implements a Time-To-Live eviction policy
type TTL struct {
	mu          sync.RWMutex
	items       map[any]time.Time
	ttl         time.Duration
	cleanupDone chan struct{}
}

// NewTTL creates a new TTL eviction policy
func NewTTL(ttl time.Duration) *TTL {
	t := &TTL{
		items:       make(map[any]time.Time),
		ttl:         ttl,
		cleanupDone: make(chan struct{}),
	}

	// Start background cleanup goroutine
	if ttl > 0 {
		go t.cleanupLoop()
	}

	return t
}

// cleanupLoop runs periodic cleanup of expired entries
func (t *TTL) cleanupLoop() {
	ticker := time.NewTicker(t.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.cleanup()
		case <-t.cleanupDone:
			return
		}
	}
}

// cleanup removes expired entries
func (t *TTL) cleanup() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	for key, expiry := range t.items {
		if now.After(expiry) {
			delete(t.items, key)
		}
	}
}

// OnAccess implements EvictionPolicy
func (t *TTL) OnAccess(key any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if expiry, ok := t.items[key]; ok {
		if time.Now().After(expiry) {
			delete(t.items, key)
		}
	}
}

// OnAdd implements EvictionPolicy
func (t *TTL) OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.items[key] = createdAt.Add(t.ttl)
}

// OnRemove implements EvictionPolicy
func (t *TTL) OnRemove(key any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.items, key)
}

// SelectVictim implements EvictionPolicy
func (t *TTL) SelectVictim() (any, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	now := time.Now()
	for key, expiry := range t.items {
		if now.After(expiry) {
			return key, true
		}
	}

	return nil, false
}

// Len implements EvictionPolicy
func (t *TTL) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.items)
}

// Close stops the cleanup goroutine
func (t *TTL) Close() {
	close(t.cleanupDone)
}
