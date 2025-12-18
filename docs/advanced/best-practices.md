# Best Practices

This guide covers recommended patterns and practices for using Synapse effectively in production environments.

## Cache Design

### Choose the Right Key Type

Use meaningful, consistent key formats:

```go
// Good: Structured, predictable keys
cache.Set(ctx, "user:12345:profile", userData)
cache.Set(ctx, "product:sku-abc:details", productData)

// Avoid: Inconsistent or opaque keys
cache.Set(ctx, "u12345", userData)
cache.Set(ctx, "abc_product_data", productData)
```

For similarity search, keys should have semantic meaning:

```go
// Good: Similar queries have similar keys
cache.Set(ctx, "search:red running shoes", results)
cache.Set(ctx, "search:blue running shoes", results)

// These will match with similarity search
value, _, _, _ := cache.GetSimilar(ctx, "search:red runing shoes")
```

### Set Appropriate Thresholds

The similarity threshold should match your use case:

```go
// High threshold (0.9+): Strict matching, fewer false positives
cache := synapse.New[string, string](
    synapse.WithThreshold(0.9),
)
// Good for: Typo correction, near-exact matching

// Medium threshold (0.7-0.8): Balanced
cache := synapse.New[string, string](
    synapse.WithThreshold(0.8),
)
// Good for: General fuzzy matching

// Low threshold (0.5-0.6): Loose matching, more results
cache := synapse.New[string, string](
    synapse.WithThreshold(0.6),
)
// Good for: Semantic similarity, broad matching
```

### Size Your Cache Appropriately

Consider memory usage when sizing:

```go
// Each entry stores:
// - Key (size depends on type)
// - Value (size depends on type)
// - Metadata: timestamps, counters, namespace, metadata map
// - Overhead: map entries, slice elements

// For string keys/values averaging 100 bytes each:
// ~300-400 bytes per entry including overhead
// 10,000 entries ≈ 3-4 MB

cache := synapse.New[string, string](
    synapse.WithMaxSize(10000),
)
```

## Performance Optimization

### Try Exact Match First

Exact lookups are O(1), similarity search is O(n):

```go
func getCached(ctx context.Context, cache *synapse.Cache[string, string], key string) (string, bool) {
    // Fast path: exact match
    if value, found := cache.Get(ctx, key); found {
        return value, true
    }

    // Slow path: similarity search
    if value, _, score, found := cache.GetSimilar(ctx, key); found && score > 0.7 {
        return value, true
    }

    return "", false
}
```

### Use Appropriate Shard Count

Balance concurrency vs. overhead:

```go
// Low concurrency (< 10 goroutines)
synapse.WithShards(4)

// Medium concurrency (10-100 goroutines)
synapse.WithShards(16)  // Default

// High concurrency (100+ goroutines)
synapse.WithShards(64)

// Very high concurrency (1000+ goroutines)
synapse.WithShards(128)
```

More shards reduce lock contention but increase memory and similarity search overhead.

### Limit Similarity Search Scope

For large caches, consider partitioning:

```go
// Use namespaces to limit similarity search scope
userCtx := synapse.WithNamespace(ctx, "user-queries")
productCtx := synapse.WithNamespace(ctx, "product-queries")

// Searches only scan entries in the same namespace
cache.GetSimilar(userCtx, "search term")
```

### Choose the Right Similarity Algorithm

Match the algorithm to your data:

```go
// For text with typos: Levenshtein or Damerau-Levenshtein
cache.WithSimilarity(algorithms.Levenshtein)

// For fixed-length identifiers: Hamming
cache.WithSimilarity(algorithms.Hamming)

// For numeric vectors: Euclidean or Manhattan
cache.WithSimilarity(func(a, b []float64) float64 {
    return algorithms.Euclidean(a, b)
})
```

## Context Usage

### Always Use Context

Pass context for cancellation and tracing:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Request context with deadline

    value, found := cache.Get(ctx, key)
    if !found {
        // Context cancellation handled internally
        value, _, _, found = cache.GetSimilar(ctx, key)
    }
}
```

### Set Timeouts for Similarity Search

Similarity search can be slow for large caches:

```go
ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
defer cancel()

value, key, score, found := cache.GetSimilar(ctx, query)
if ctx.Err() == context.DeadlineExceeded {
    // Search took too long, use fallback
}
```

### Use Namespaces for Multi-Tenancy

Isolate data between tenants:

```go
func getTenantCache(ctx context.Context, tenantID string) context.Context {
    return synapse.WithNamespace(ctx, fmt.Sprintf("tenant:%s", tenantID))
}

// Each tenant's operations are isolated
tenant1Ctx := getTenantCache(ctx, "tenant-1")
cache.Set(tenant1Ctx, "config", config)

// Similarity search only sees tenant-1's entries
cache.GetSimilar(tenant1Ctx, "confg")
```

## Error Handling

### Handle Context Errors

```go
err := cache.Set(ctx, key, value)
if err != nil {
    if errors.Is(err, context.Canceled) {
        return fmt.Errorf("operation cancelled: %w", err)
    }
    if errors.Is(err, context.DeadlineExceeded) {
        return fmt.Errorf("operation timeout: %w", err)
    }
    return fmt.Errorf("cache set failed: %w", err)
}
```

### Don't Rely on Cache Availability

Design for cache misses:

```go
func getData(ctx context.Context, id string) (Data, error) {
    // Try cache first
    if cached, found := cache.Get(ctx, id); found {
        return cached, nil
    }

    // Fall back to source
    data, err := fetchFromSource(ctx, id)
    if err != nil {
        return Data{}, err
    }

    // Best-effort cache population
    _ = cache.Set(ctx, id, data)

    return data, nil
}
```

## TTL and Eviction

### Use TTL for Time-Sensitive Data

```go
// Short TTL for frequently changing data
cache := synapse.New[string, MarketData](
    synapse.WithTTL(5 * time.Second),
)

// Longer TTL for stable data
cache := synapse.New[string, UserProfile](
    synapse.WithTTL(15 * time.Minute),
)
```

### Configure Eviction for Your Workload

```go
// LRU for general workloads with temporal locality
cache := synapse.New[string, string](
    synapse.WithMaxSize(10000),
    synapse.WithEviction(eviction.NewLRU(10000)),
)
```

### Align Eviction Policy Size with Cache Size

```go
maxSize := 10000

cache := synapse.New[string, string](
    synapse.WithMaxSize(maxSize),
    synapse.WithEviction(eviction.NewLRU(maxSize)),  // Same size
)
```

## Testing

### Test with Race Detector

```bash
go test -race ./...
```

### Test Context Cancellation

```go
func TestCacheContextCancellation(t *testing.T) {
    cache := synapse.New[string, string]()
    cache.WithSimilarity(algorithms.Levenshtein)

    // Populate cache
    ctx := context.Background()
    for i := 0; i < 10000; i++ {
        cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
    }

    // Test cancellation
    cancelCtx, cancel := context.WithCancel(ctx)
    cancel()  // Cancel immediately

    _, _, _, found := cache.GetSimilar(cancelCtx, "search")
    if found {
        t.Error("expected no result after cancellation")
    }
}
```

### Benchmark Your Use Case

```go
func BenchmarkCacheGet(b *testing.B) {
    cache := synapse.New[string, string](
        synapse.WithMaxSize(10000),
    )

    ctx := context.Background()
    cache.Set(ctx, "key", "value")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Get(ctx, "key")
    }
}

func BenchmarkCacheGetSimilar(b *testing.B) {
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

## Production Checklist

- [ ] Set appropriate `MaxSize` based on memory constraints
- [ ] Configure `Shards` based on expected concurrency
- [ ] Set `Threshold` appropriate for your similarity requirements
- [ ] Configure `TTL` for time-sensitive data
- [ ] Use `Eviction` policy (LRU recommended for most cases)
- [ ] Use namespaces for multi-tenant isolation
- [ ] Always pass `context.Context` to operations
- [ ] Set timeouts for similarity searches
- [ ] Test with `-race` flag
- [ ] Benchmark critical paths
- [ ] Monitor cache hit rates in production

## Further Reading

- [Performance Tuning](performance-tuning.md) - Detailed optimization guide
- [Core Concepts](../core-concepts/synapse.md) - Architecture details
- [API Reference](../api-reference/synapse.md) - Complete API documentation
