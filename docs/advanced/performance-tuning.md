# Performance Tuning

This guide covers performance optimization techniques for Synapse caches in production environments.

## Understanding Performance Characteristics

### Operation Complexity

| Operation    | Time Complexity | Space Complexity | Lock Type            |
| ------------ | --------------- | ---------------- | -------------------- |
| `Get`        | O(1)            | O(1)             | RLock (single shard) |
| `Set`        | O(1) amortized  | O(1)             | Lock (single shard)  |
| `Delete`     | O(n) per shard  | O(1)             | Lock (single shard)  |
| `GetSimilar` | O(n × s)        | O(1)             | RLock (all shards)   |
| `Len`        | O(s)            | O(1)             | RLock (all shards)   |

Where:

- n = entries per shard
- s = number of shards

### Memory Layout

Each entry consumes:

```
Entry[K, V]:
├── Key K                    (varies)
├── Value V                  (varies)
├── CreatedAt time.Time      (24 bytes)
├── AccessedAt time.Time     (24 bytes)
├── AccessCount uint64       (8 bytes)
├── ExpiresAt time.Time      (24 bytes)
├── Metadata map[string]any  (8 bytes + contents)
└── Namespace string         (16 bytes + length)

Base overhead: ~104 bytes + key + value + metadata
```

## Sharding Optimization

### Shard Count Selection

The number of shards affects both concurrency and memory:

```go
// Formula: shards ≈ expectedConcurrency / 4
// Minimum: 4, Maximum: 256

// Low concurrency (1-10 goroutines)
synapse.WithShards(4)

// Medium concurrency (10-50 goroutines)
synapse.WithShards(16)  // Default

// High concurrency (50-200 goroutines)
synapse.WithShards(32)

// Very high concurrency (200+ goroutines)
synapse.WithShards(64)
```

### Shard Count Trade-offs

| More Shards               | Fewer Shards                  |
| ------------------------- | ----------------------------- |
| ✅ Less lock contention   | ✅ Lower memory overhead      |
| ✅ Better parallel writes | ✅ Faster similarity search   |
| ❌ More memory overhead   | ❌ More lock contention       |
| ❌ Slower `GetSimilar`    | ❌ Worse parallel performance |

### Benchmarking Shard Count

```go
func BenchmarkShardCount(b *testing.B) {
    shardCounts := []int{4, 8, 16, 32, 64}

    for _, shards := range shardCounts {
        b.Run(fmt.Sprintf("shards=%d", shards), func(b *testing.B) {
            cache := synapse.New[string, string](
                synapse.WithShards(shards),
                synapse.WithMaxSize(10000),
            )

            ctx := context.Background()

            // Populate cache
            for i := 0; i < 10000; i++ {
                cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
            }

            b.ResetTimer()
            b.RunParallel(func(pb *testing.PB) {
                i := 0
                for pb.Next() {
                    cache.Get(ctx, fmt.Sprintf("key-%d", i%10000))
                    i++
                }
            })
        })
    }
}
```

## Similarity Search Optimization

### Limit Search Scope

Similarity search is O(n) per shard. Reduce n with namespaces:

```go
// Instead of one large namespace
cache.Set(ctx, "query:electronics:laptop", result)
cache.Set(ctx, "query:clothing:shirt", result)

// Use namespaces to partition
electronicsCtx := synapse.WithNamespace(ctx, "electronics")
clothingCtx := synapse.WithNamespace(ctx, "clothing")

cache.Set(electronicsCtx, "query:laptop", result)
cache.Set(clothingCtx, "query:shirt", result)

// Similarity search only scans relevant entries
cache.GetSimilar(electronicsCtx, "query:laptap")  // Only scans electronics
```

### Use Timeouts

Prevent slow searches from blocking:

```go
func searchWithTimeout(cache *synapse.Cache[string, string], query string) (string, bool) {
    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
    defer cancel()

    value, _, _, found := cache.GetSimilar(ctx, query)
    return value, found
}
```

### Optimize Similarity Function

The similarity function runs for every entry during search:

```go
// Slow: Complex computation
func slowSimilarity(a, b string) float64 {
    // Expensive regex or external call
    return expensiveComputation(a, b)
}

// Fast: Simple comparison
func fastSimilarity(a, b string) float64 {
    if a == b {
        return 1.0
    }
    // Quick early exit for obviously different strings
    if len(a) == 0 || len(b) == 0 {
        return 0.0
    }
    if abs(len(a)-len(b)) > 5 {
        return 0.0  // Too different in length
    }
    return algorithms.Levenshtein(a, b)
}
```

### Pre-filter with Length Check

Add early exits in custom similarity functions:

```go
func optimizedLevenshtein(a, b string) float64 {
    if a == b {
        return 1.0
    }

    lenA, lenB := len(a), len(b)
    if lenA == 0 || lenB == 0 {
        return 0.0
    }

    // Quick rejection: if lengths differ too much, similarity will be low
    maxLen := lenA
    if lenB > maxLen {
        maxLen = lenB
    }
    minLen := lenA
    if lenB < minLen {
        minLen = lenB
    }

    // If more than 50% length difference, max similarity is < 0.5
    if float64(minLen)/float64(maxLen) < 0.5 {
        return 0.0
    }

    return algorithms.Levenshtein(a, b)
}
```

## Memory Optimization

### Choose Compact Key Types

```go
// Less efficient: long string keys
cache := synapse.New[string, Data]()
cache.Set(ctx, "very-long-key-that-uses-lots-of-memory", data)

// More efficient: shorter keys
cache.Set(ctx, "k:12345", data)

// Most efficient: numeric keys (when similarity isn't needed)
numCache := synapse.New[int64, Data]()
numCache.Set(ctx, 12345, data)
```

### Limit Metadata Usage

Metadata adds overhead per entry:

```go
// Expensive: large metadata maps
entry.Metadata = map[string]any{
    "field1": "value1",
    "field2": "value2",
    // ... many fields
}

// Better: store metadata in the value
type CachedData struct {
    Data     ActualData
    Metadata DataMetadata
}
cache := synapse.New[string, CachedData]()
```

### Use TTL to Bound Memory

```go
cache := synapse.New[string, string](
    synapse.WithMaxSize(100000),
    synapse.WithTTL(15 * time.Minute),
    synapse.WithEviction(eviction.NewLRU(100000)),
)
```

## Eviction Policy Optimization

### LRU Performance

LRU operations are O(1) but have lock overhead:

```go
// LRU tracks every access
cache.Get(ctx, key)  // Updates LRU position

// For read-heavy workloads, consider if eviction tracking is needed
```

### Eviction Sizing

Match eviction policy size to cache size:

```go
maxSize := 10000
shards := 16

// Policy tracks all entries
lru := eviction.NewLRU(maxSize)

cache := synapse.New[string, string](
    synapse.WithMaxSize(maxSize),
    synapse.WithShards(shards),
    synapse.WithEviction(lru),
)
```

## Benchmarking

### Standard Benchmarks

```go
func BenchmarkGet(b *testing.B) {
    cache := synapse.New[string, string](synapse.WithMaxSize(10000))
    ctx := context.Background()
    cache.Set(ctx, "key", "value")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Get(ctx, "key")
    }
}

func BenchmarkSet(b *testing.B) {
    cache := synapse.New[string, string](synapse.WithMaxSize(b.N))
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
    }
}

func BenchmarkGetSimilar(b *testing.B) {
    cache := synapse.New[string, string](
        synapse.WithMaxSize(1000),
        synapse.WithThreshold(0.8),
    )
    cache.WithSimilarity(algorithms.Levenshtein)

    ctx := context.Background()
    for i := 0; i < 1000; i++ {
        cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.GetSimilar(ctx, "key-500")
    }
}
```

### Parallel Benchmarks

```go
func BenchmarkParallelGet(b *testing.B) {
    cache := synapse.New[string, string](
        synapse.WithMaxSize(10000),
        synapse.WithShards(32),
    )

    ctx := context.Background()
    for i := 0; i < 10000; i++ {
        cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            cache.Get(ctx, fmt.Sprintf("key-%d", i%10000))
            i++
        }
    })
}
```

### Memory Benchmarks

```go
func BenchmarkMemory(b *testing.B) {
    b.ReportAllocs()

    cache := synapse.New[string, string](synapse.WithMaxSize(b.N))
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
    }
}
```

### Running Benchmarks

```bash
# Basic benchmarks
go test -bench=. -benchmem ./...

# With CPU profiling
go test -bench=BenchmarkGetSimilar -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# With memory profiling
go test -bench=BenchmarkMemory -memprofile=mem.prof ./...
go tool pprof mem.prof

# Longer benchmark runs for accuracy
go test -bench=. -benchtime=5s ./...
```

## Profiling

### CPU Profiling

```go
import (
    "os"
    "runtime/pprof"
)

func profileSimilaritySearch() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Run similarity searches
    for i := 0; i < 1000; i++ {
        cache.GetSimilar(ctx, "query")
    }
}
```

### Memory Profiling

```go
func profileMemory() {
    // Force GC before profiling
    runtime.GC()

    f, _ := os.Create("mem.prof")
    pprof.WriteHeapProfile(f)
    f.Close()
}
```

## Performance Targets

Typical performance on modern hardware:

| Operation    | Target  | Notes                       |
| ------------ | ------- | --------------------------- |
| `Get` (hit)  | < 100ns | Single shard, no contention |
| `Get` (miss) | < 50ns  | Single shard, no contention |
| `Set`        | < 200ns | Without eviction            |
| `Set`        | < 500ns | With LRU eviction           |
| `GetSimilar` | < 1ms   | 1000 entries, Levenshtein   |
| `GetSimilar` | < 10ms  | 10000 entries, Levenshtein  |

## Quick Wins Checklist

- [ ] Set shard count based on concurrency level
- [ ] Use namespaces to partition data
- [ ] Add timeouts to similarity searches
- [ ] Optimize similarity function with early exits
- [ ] Use compact key types
- [ ] Align eviction policy size with cache size
- [ ] Run benchmarks with `-race` flag
- [ ] Profile before optimizing
- [ ] Test under realistic load

## Further Reading

- [Best Practices](best-practices.md) - Production recommendations
- [Core Concepts](../core-concepts/synapse.md) - Architecture details
- [Algorithms](../core-concepts/algorithms.md) - Algorithm performance
