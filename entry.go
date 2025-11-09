package synapse

import (
	"time"
)

// Entry represents a cache entry with metadata
type Entry[K comparable, V any] struct {
	Key         K
	Value       V
	CreatedAt   time.Time
	AccessedAt  time.Time
	AccessCount uint64
	ExpiresAt   time.Time
	Metadata    map[string]any
	Namespace   string
}

// newEntry creates a new cache entry
func newEntry[K comparable, V any](key K, value V, ttl time.Duration, namespace string) *Entry[K, V] {
	now := time.Now()
	entry := &Entry[K, V]{
		Key:         key,
		Value:       value,
		CreatedAt:   now,
		AccessedAt:  now,
		AccessCount: 0,
		Namespace:   namespace,
		Metadata:    make(map[string]any),
	}

	if ttl > 0 {
		entry.ExpiresAt = now.Add(ttl)
	}

	return entry
}

// IsExpired checks if the entry has expired
func (e *Entry[K, V]) IsExpired() bool {
	if e.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// Touch updates the access time and increments access count
func (e *Entry[K, V]) Touch() {
	e.AccessedAt = time.Now()
	e.AccessCount++
}
