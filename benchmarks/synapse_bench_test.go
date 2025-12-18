package benchmarks

import (
	"context"
	"fmt"
	"testing"

	"github.com/kolosys/synapse"
	"github.com/kolosys/synapse/algorithms"
	"github.com/kolosys/synapse/eviction"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := synapse.New[string, string](
		synapse.WithShards(16),
		synapse.WithMaxSize(10000),
	)
	ctx := context.Background()

	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(ctx, key, "value")
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := synapse.New[string, string](
		synapse.WithShards(16),
		synapse.WithMaxSize(10000),
	)
	ctx := context.Background()

	// Pre-populate cache
	for i := range 1000 {
		key := fmt.Sprintf("key%d", i)
		cache.Set(ctx, key, "value")
	}

	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Get(ctx, key)
	}
}

func BenchmarkCacheGetSimilar(b *testing.B) {
	cache := synapse.New[string, string](
		synapse.WithShards(16),
		synapse.WithMaxSize(1000),
		synapse.WithThreshold(0.7),
	)
	cache.WithSimilarity(algorithms.Levenshtein)
	ctx := context.Background()

	// Pre-populate cache
	words := []string{"hello", "world", "test", "example", "cache", "benchmark"}
	for _, word := range words {
		cache.Set(ctx, word, word)
	}

	for b.Loop() {
		cache.GetSimilar(ctx, "helo")
	}
}

func BenchmarkCacheConcurrentSet(b *testing.B) {
	cache := synapse.New[string, string](
		synapse.WithShards(32),
		synapse.WithMaxSize(100000),
	)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			cache.Set(ctx, key, "value")
			i++
		}
	})
}

func BenchmarkCacheConcurrentGet(b *testing.B) {
	cache := synapse.New[string, string](
		synapse.WithShards(32),
		synapse.WithMaxSize(10000),
	)
	ctx := context.Background()

	// Pre-populate cache
	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		cache.Set(ctx, key, "value")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%10000)
			cache.Get(ctx, key)
			i++
		}
	})
}

func BenchmarkCacheWithLRU(b *testing.B) {
	policy := eviction.NewLRU(1000)
	cache := synapse.New[int, string](
		synapse.WithMaxSize(1000),
		synapse.WithEviction(policy),
	)
	ctx := context.Background()

	for i := 0; b.Loop(); i++ {
		cache.Set(ctx, i, "value")
		if i > 0 {
			cache.Get(ctx, i-1)
		}
	}
}

func BenchmarkCacheWithNamespace(b *testing.B) {
	cache := synapse.New[string, string](
		synapse.WithShards(16),
	)
	ctx := synapse.WithNamespace(context.Background(), "test")

	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(ctx, key, "value")
		cache.Get(ctx, key)
	}
}

func BenchmarkCacheSharding_16(b *testing.B) {
	benchmarkSharding(b, 16)
}

func BenchmarkCacheSharding_32(b *testing.B) {
	benchmarkSharding(b, 32)
}

func BenchmarkCacheSharding_64(b *testing.B) {
	benchmarkSharding(b, 64)
}

func BenchmarkCacheSharding_128(b *testing.B) {
	benchmarkSharding(b, 128)
}

func benchmarkSharding(b *testing.B, numShards int) {
	cache := synapse.New[string, string](
		synapse.WithShards(numShards),
		synapse.WithMaxSize(100000),
	)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			cache.Set(ctx, key, "value")
			cache.Get(ctx, key)
			i++
		}
	})
}
