# Synapse Cache

The core `synapse` package provides a high-performance, generic similarity-based cache with intelligent sharding.

**Import Path:** `github.com/kolosys/synapse`

## Overview

Synapse is more than a traditional key-value cache. It supports **similarity-based lookups**, allowing you to find entries with keys that are "close enough" to your query when exact matches don't exist.

## Architecture

### Sharded Design

Synapse distributes entries across multiple shards to reduce lock contention and improve concurrent performance:

```
┌──────────────────────────────────────────────────────────┐
│                        Cache[K, V]                       │
├──────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ Shard 0  │  │ Shard 1  │  │ Shard 2  │  │ Shard N  │  │
│  ├──────────┤  ├──────────┤  ├──────────┤  ├──────────┤  │
│  │ RWMutex  │  │ RWMutex  │  │ RWMutex  │  │ RWMutex  │  │
│  │ data map │  │ data map │  │ data map │  │ data map │  │
│  │ keys []K │  │ keys []K │  │ keys []K │  │ keys []K │  │
│  │ eviction │  │ eviction │  │ eviction │  │ eviction │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
└──────────────────────────────────────────────────────────┘
```

Each shard operates independently with:

- `sync.RWMutex` for thread-safe access
- Hash map for O(1) exact lookups
- Key slice for similarity search iteration
- Individual eviction policy tracking

### Key Distribution

Keys are distributed to shards using FNV-1a hashing:

```go
// Simplified key-to-shard mapping
hash := fnv.New64a()
hash.Write([]byte(keyToString(key)))
shardIndex := hash.Sum64() % uint64(numShards)
```

This ensures even distribution across shards for most key patterns.

## Core Types

### Cache[K, V]

The main cache type, generic over key type `K` (must be `comparable`) and value type `V` (any type):

```go
type Cache[K comparable, V any] struct {
    shards     []*Shard[K, V]
    similarity SimilarityFunc[K]
    threshold  float64
    options    *Options
}
```

Create a cache with functional options:

```go
cache := synapse.New[string, string](
    synapse.WithShards(16),
    synapse.WithMaxSize(10000),
    synapse.WithThreshold(0.8),
    synapse.WithTTL(5 * time.Minute),
    synapse.WithEviction(eviction.NewLRU(10000)),
)
```

### Entry[K, V]

Each cache entry contains the value and metadata:

```go
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
```

The cache automatically manages:

- `CreatedAt` - Set when the entry is created
- `AccessedAt` - Updated on each access
- `AccessCount` - Incremented on each access
- `ExpiresAt` - Set based on TTL configuration

### SimilarityFunc[K]

A function type for computing similarity between keys:

```go
type SimilarityFunc[K comparable] func(a, b K) float64
```

Must return a score between `0.0` (completely different) and `1.0` (identical).

## Core Operations

### Set

Store a key-value pair:

```go
err := cache.Set(ctx, "key", "value")
if err != nil {
    // Handle context cancellation
}
```

**Behavior:**

1. Check context cancellation
2. Hash key to determine shard
3. If key exists, update value and touch entry
4. If cache is full, trigger eviction
5. Create new entry with metadata

### Get

Retrieve a value by exact key match:

```go
value, found := cache.Get(ctx, "key")
if found {
    fmt.Println(value)
}
```

**Behavior:**

1. Check context cancellation
2. Hash key to determine shard
3. Look up key in shard's data map
4. Verify namespace match (if set)
5. Check TTL expiration
6. Update access tracking

### GetSimilar

Find the most similar key above the threshold:

```go
value, matchedKey, score, found := cache.GetSimilar(ctx, "query")
if found {
    fmt.Printf("Found %v at key %v (score: %.2f)\n", value, matchedKey, score)
}
```

**Behavior:**

1. Search ALL shards (not just the target shard)
2. For each entry, compute similarity score
3. Track the best match above threshold
4. Respect context cancellation between iterations
5. Return best match across all shards

### Delete

Remove an entry:

```go
deleted := cache.Delete(ctx, "key")
if deleted {
    fmt.Println("Entry removed")
}
```

### Len

Get total entry count:

```go
count := cache.Len()
```

Sums the count across all shards.

## Configuration Options

### WithShards

Set the number of shards (1-256):

```go
synapse.WithShards(32)
```

More shards reduce lock contention but increase memory overhead. Default: 16.

### WithMaxSize

Set the maximum total entries:

```go
synapse.WithMaxSize(10000)
```

Distributed across shards: `maxSizePerShard = maxSize / numShards`. Default: 1000.

### WithThreshold

Set the minimum similarity score for matches:

```go
synapse.WithThreshold(0.8)  // 80% similarity required
```

Must be between 0.0 and 1.0. Default: 0.8.

### WithEviction

Set the eviction policy:

```go
synapse.WithEviction(eviction.NewLRU(10000))
```

Default: nil (FIFO fallback).

### WithTTL

Set time-to-live for entries:

```go
synapse.WithTTL(5 * time.Minute)
```

Default: 0 (no expiration).

### WithStats

Enable statistics tracking:

```go
synapse.WithStats(true)
```

Default: false.

## Context Integration

### Namespace Isolation

Partition entries by namespace:

```go
tenant1Ctx := synapse.WithNamespace(ctx, "tenant-1")
tenant2Ctx := synapse.WithNamespace(ctx, "tenant-2")

// Same key, different namespaces
cache.Set(tenant1Ctx, "config", "tenant1-config")
cache.Set(tenant2Ctx, "config", "tenant2-config")

// Each tenant sees only their data
value, _ := cache.Get(tenant1Ctx, "config")  // "tenant1-config"
```

### Metadata

Attach metadata to the context:

```go
ctx := synapse.WithMetadata(ctx, "request_id", "abc123")
requestID, ok := synapse.GetMetadata(ctx, "request_id")
```

### Cancellation

All operations respect context cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// Long similarity search will be cancelled
_, _, _, found := cache.GetSimilar(ctx, "query")
```

## Lookup Strategies

### Exact Match First

For best performance, try exact match before similarity search:

```go
// Fast O(1) exact lookup
if value, found := cache.Get(ctx, key); found {
    return value, nil
}

// Slower O(n) similarity search as fallback
if value, _, _, found := cache.GetSimilar(ctx, key); found {
    return value, nil
}

return nil, ErrNotFound
```

### Similarity Only

When you always want the best match:

```go
value, matchedKey, score, found := cache.GetSimilar(ctx, query)
if !found {
    return nil, ErrNotFound
}

// Optionally check if it was an exact match
if score == 1.0 {
    // Exact match
}
```

## Performance Characteristics

| Operation    | Complexity | Notes                                |
| ------------ | ---------- | ------------------------------------ |
| `Get`        | O(1)       | Single shard lookup                  |
| `Set`        | O(1)       | Amortized, may trigger eviction      |
| `Delete`     | O(n)       | O(1) map delete + O(n) slice removal |
| `GetSimilar` | O(n×s)     | n = entries per shard, s = shards    |
| `Len`        | O(s)       | s = number of shards                 |

## Thread Safety

- **Read operations** (`Get`, `GetSimilar`, `Len`) use `RLock`
- **Write operations** (`Set`, `Delete`) use `Lock`
- Each shard has its own mutex, reducing contention
- Context cancellation is checked without holding locks

## Further Reading

- [Algorithms](algorithms.md) - Similarity algorithms
- [Eviction Policies](eviction.md) - Cache eviction strategies
- [Best Practices](../advanced/best-practices.md) - Production recommendations
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization guide
- [API Reference](../api-reference/synapse.md) - Complete API documentation
