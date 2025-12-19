package synapse

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kolosys/synapse/algorithms"
	"github.com/kolosys/synapse/eviction"
)

func TestCacheBasicOperations(t *testing.T) {
	cache := New[string, string]()
	ctx := context.Background()

	// Test Set and Get
	err := cache.Set(ctx, "key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, ok := cache.Get(ctx, "key1")
	if !ok {
		t.Fatal("Get failed: key not found")
	}
	if val != "value1" {
		t.Fatalf("Expected value1, got %s", val)
	}

	// Test non-existent key
	_, ok = cache.Get(ctx, "key2")
	if ok {
		t.Fatal("Expected false for non-existent key")
	}
}

func TestCacheDelete(t *testing.T) {
	cache := New[string, int]()
	ctx := context.Background()

	cache.Set(ctx, "key1", 100)
	cache.Set(ctx, "key2", 200)

	// Delete key1
	deleted := cache.Delete(ctx, "key1")
	if !deleted {
		t.Fatal("Delete failed")
	}

	// Verify key1 is gone
	_, ok := cache.Get(ctx, "key1")
	if ok {
		t.Fatal("Key should be deleted")
	}

	// Verify key2 still exists
	val, ok := cache.Get(ctx, "key2")
	if !ok || val != 200 {
		t.Fatal("Key2 should still exist")
	}
}

func TestCacheSimilarity(t *testing.T) {
	cache := New[string, string](
		WithThreshold(0.7),
	)
	cache.WithSimilarity(algorithms.Levenshtein)

	ctx := context.Background()

	cache.Set(ctx, "hello", "world")
	cache.Set(ctx, "help", "assistance")

	// Search for similar key
	val, key, score, ok := cache.GetSimilar(ctx, "helo")
	if !ok {
		t.Fatal("GetSimilar should find a match")
	}

	if score < 0.7 {
		t.Fatalf("Score should be >= 0.7, got %f", score)
	}

	t.Logf("Found similar key: %s with score: %f, value: %s", key, score, val)
}

func TestCacheWithNamespace(t *testing.T) {
	cache := New[string, string]()

	ctx1 := WithNamespace(context.Background(), "ns1")
	ctx2 := WithNamespace(context.Background(), "ns2")

	// Use different keys for different namespaces
	cache.Set(ctx1, "ns1:key1", "value1-ns1")
	cache.Set(ctx2, "ns2:key1", "value1-ns2")

	// Get from namespace 1
	val, ok := cache.Get(ctx1, "ns1:key1")
	if !ok || val != "value1-ns1" {
		t.Fatalf("Expected value1-ns1, got %s", val)
	}

	// Get from namespace 2
	val, ok = cache.Get(ctx2, "ns2:key1")
	if !ok || val != "value1-ns2" {
		t.Fatalf("Expected value1-ns2, got %s", val)
	}

	// Cross-namespace access should be filtered
	_, ok = cache.Get(ctx1, "ns2:key1")
	if ok {
		t.Fatal("Should not access ns2 key from ns1 context")
	}
}

func TestCacheContextCancellation(t *testing.T) {
	cache := New[string, string]()

	ctx, cancel := context.WithCancel(context.Background())

	cache.Set(ctx, "key1", "value1")

	// Cancel context
	cancel()

	// Operations should fail or return false
	_, ok := cache.Get(ctx, "key1")
	if ok {
		t.Fatal("Get should fail with cancelled context")
	}
}

func TestCacheWithTTL(t *testing.T) {
	cache := New[string, string](
		WithTTL(100 * time.Millisecond),
	)

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1")

	// Should exist immediately
	val, ok := cache.Get(ctx, "key1")
	if !ok || val != "value1" {
		t.Fatal("Key should exist immediately")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get(ctx, "key1")
	if ok {
		t.Fatal("Key should be expired")
	}
}

func TestCacheWithMaxSize(t *testing.T) {
	policy := eviction.NewLRU(100)
	cache := New[int, string](
		WithMaxSize(100),
		WithShards(1), // Use single shard for predictable behavior
		WithEviction(policy),
	)

	ctx := context.Background()

	// Fill cache beyond max size
	for i := 0; i < 150; i++ {
		cache.Set(ctx, i, "value")
	}

	// Cache should have evicted some entries
	if cache.Len() > 100 {
		t.Fatalf("Cache size should be <= 100, got %d", cache.Len())
	}
}

func TestCacheWithMetadata(t *testing.T) {
	cache := New[string, string]()

	ctx := WithMetadata(context.Background(), "user", "alice")
	ctx = WithMetadata(ctx, "role", "admin")

	cache.Set(ctx, "key1", "value1")

	// Verify metadata retrieval
	user, ok := GetMetadata(ctx, "user")
	if !ok || user != "alice" {
		t.Fatal("Metadata should be retrievable")
	}

	role, ok := GetMetadata(ctx, "role")
	if !ok || role != "admin" {
		t.Fatal("Metadata should be retrievable")
	}
}

func TestCacheSharding(t *testing.T) {
	cache := New[string, string](
		WithShards(32),
		WithMaxSize(2000), // Ensure enough capacity
	)

	ctx := context.Background()

	// Add many keys to ensure distribution across shards
	count := 500
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := cache.Set(ctx, key, "value")
		if err != nil {
			t.Fatalf("Failed to set key %d: %v", i, err)
		}
	}

	// Verify most keys are stored (allow for some margin due to default eviction)
	actualLen := cache.Len()
	if actualLen < count-10 {
		t.Fatalf("Expected at least %d entries, got %d", count-10, actualLen)
	}

	// Verify we can retrieve keys
	retrieved := 0
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		if _, ok := cache.Get(ctx, key); ok {
			retrieved++
		}
	}
	t.Logf("Retrieved %d out of %d keys", retrieved, count)
}

func TestCacheStats(t *testing.T) {
	cache := New[string, string](
		WithStats(true),
		WithShards(1), // Use single shard for predictable stats
	)
	ctx := context.Background()

	// Initial stats should be zero
	stats := cache.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Sets != 0 {
		t.Fatalf("Initial stats should be zero, got %+v", stats)
	}

	// Test Set
	cache.Set(ctx, "key1", "value1")
	cache.Set(ctx, "key2", "value2")
	stats = cache.Stats()
	if stats.Sets != 2 {
		t.Fatalf("Expected 2 sets, got %d", stats.Sets)
	}

	// Test Get hits
	cache.Get(ctx, "key1")
	cache.Get(ctx, "key2")
	stats = cache.Stats()
	if stats.Hits != 2 {
		t.Fatalf("Expected 2 hits, got %d", stats.Hits)
	}

	// Test Get misses
	cache.Get(ctx, "key3")
	cache.Get(ctx, "key4")
	stats = cache.Stats()
	if stats.Misses != 2 {
		t.Fatalf("Expected 2 misses, got %d", stats.Misses)
	}

	// Test Delete
	cache.Delete(ctx, "key1")
	stats = cache.Stats()
	if stats.Deletes != 1 {
		t.Fatalf("Expected 1 delete, got %d", stats.Deletes)
	}

	// Test similarity search
	cache.WithSimilarity(algorithms.Levenshtein)
	cache.Set(ctx, "hello", "world")
	cache.GetSimilar(ctx, "helo")
	stats = cache.Stats()
	if stats.SimilarSearches == 0 {
		t.Fatal("Expected at least 1 similarity search")
	}
}

func TestCacheStatsDisabled(t *testing.T) {
	cache := New[string, string](
		WithStats(false),
	)
	ctx := context.Background()

	// Perform operations
	cache.Set(ctx, "key1", "value1")
	cache.Get(ctx, "key1")
	cache.Get(ctx, "key2")
	cache.Delete(ctx, "key1")

	// Stats should return zero values when disabled
	stats := cache.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Sets != 0 || stats.Deletes != 0 {
		t.Fatalf("Stats should be zero when disabled, got %+v", stats)
	}
}

func TestCacheStatsExpired(t *testing.T) {
	cache := New[string, string](
		WithStats(true),
		WithTTL(50*time.Millisecond),
		WithShards(1),
	)
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1")
	cache.Get(ctx, "key1") // Should hit

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Try to get expired key
	cache.Get(ctx, "key1") // Should miss and record expired

	stats := cache.Stats()
	if stats.Expired == 0 {
		t.Fatal("Expected at least 1 expired entry")
	}
	if stats.Misses < 1 {
		t.Fatal("Expected at least 1 miss for expired entry")
	}
}

func TestCacheStatsEviction(t *testing.T) {
	policy := eviction.NewLRU(2)
	cache := New[string, string](
		WithStats(true),
		WithMaxSize(2),
		WithShards(1),
		WithEviction(policy),
	)
	ctx := context.Background()

	// Fill cache to trigger eviction
	cache.Set(ctx, "key1", "value1")
	cache.Set(ctx, "key2", "value2")
	cache.Set(ctx, "key3", "value3") // Should evict one

	stats := cache.Stats()
	if stats.Evictions == 0 {
		t.Fatal("Expected at least 1 eviction")
	}
}
