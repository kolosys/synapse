# Synapse 🧠

A high-performance, generic similarity-based cache for Go with intelligent sharding and pluggable eviction policies.

![GoVersion](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Zero Dependencies](https://img.shields.io/badge/Zero-Dependencies-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/synapse.svg)](https://pkg.go.dev/github.com/kolosys/synapse)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/synapse)](https://goreportcard.com/report/github.com/kolosys/synapse)

## Overview

Synapse is a thread-safe, context-aware cache that goes beyond traditional key-value storage by supporting **similarity-based lookups**. When an exact key match isn't found, Synapse can find the "closest" matching key based on a configurable similarity function and threshold.

### Key Features

- **🎯 Similarity-Based Lookups**: Find approximate matches when exact keys don't exist
- **🔧 Generic Types**: Fully type-safe with Go generics (1.18+)
- **⚡ High Performance**: Automatic sharding distributes load across multiple concurrent-safe partitions
- **🧩 Pluggable Similarity Functions**: Define custom similarity algorithms for your use case
- **♻️ Eviction Policies**: Currently supports LRU with more policies coming soon
- **⏰ TTL Support**: Automatic expiration of cache entries
- **🏷️ Namespace Isolation**: Partition cache entries by namespace via context
- **🔒 Thread-Safe**: Lock-free reads and efficient write locking per shard
- **📊 Metadata Support**: Attach custom metadata to cache entries
- **🔌 Context-Aware**: Full context.Context integration for cancellation and values

## Installation

```bash
go get github.com/kolosys/synapse
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/eviction"
)

func stringSimilarity(a, b string) float64 {
    a, b = strings.ToLower(a), strings.ToLower(b)
    if a == b {
        return 1.0
    }
    minLen := min(len(a), len(b))
    matches := 0
    for i := 0; i < minLen; i++ {
        if a[i] == b[i] {
            matches++
        } else {
            break
        }
    }
    return float64(matches) / float64(max(len(a), len(b)))
}

func main() {
    ctx := context.Background()

    cache := synapse.New[string, string](
        synapse.WithMaxSize(1000),
        synapse.WithShards(16),
        synapse.WithThreshold(0.7),
        synapse.WithEviction(eviction.NewLRU(1000)),
    )

    cache.WithSimilarity(stringSimilarity)

    cache.Set(ctx, "user:alice", "Alice's data")
    cache.Set(ctx, "user:bob", "Bob's data")

    if value, found := cache.Get(ctx, "user:alice"); found {
        fmt.Println("Exact match:", value)
    }

    value, key, score, found := cache.GetSimilar(ctx, "user:ali")
    if found {
        fmt.Printf("Similar match: %s (key: %s, score: %.2f)\n", value, key, score)
    }
}
```

## API Reference

### Core Methods

- `New[K, V](opts ...Option) *Cache[K, V]` - Create a new cache instance
- `Get(ctx context.Context, key K) (V, bool)` - Retrieve value by exact key match
- `Set(ctx context.Context, key K, value V) error` - Store a key-value pair
- `GetSimilar(ctx context.Context, key K) (V, K, float64, bool)` - Find most similar key above threshold
- `Delete(ctx context.Context, key K) bool` - Remove a key from the cache
- `Len() int` - Get total number of entries across all shards
- `WithSimilarity(fn SimilarityFunc[K]) *Cache[K, V]` - Set similarity function

### Configuration Options

| Option                 | Description                    | Default           |
| ---------------------- | ------------------------------ | ----------------- |
| `WithShards(n)`        | Number of shards (1-256)       | 16                |
| `WithMaxSize(size)`    | Maximum number of entries      | 1000              |
| `WithThreshold(t)`     | Similarity threshold (0.0-1.0) | 0.8               |
| `WithEviction(policy)` | Eviction policy                | nil               |
| `WithTTL(duration)`    | Time-to-live for entries       | 0 (no expiration) |
| `WithStats(enable)`    | Enable statistics tracking     | false             |

### Context Functions

- `WithNamespace(ctx context.Context, namespace string) context.Context` - Add namespace to context
- `GetNamespace(ctx context.Context) string` - Retrieve namespace from context
- `WithMetadata(ctx context.Context, key string, value any) context.Context` - Add metadata to context
- `GetMetadata(ctx context.Context, key string) (any, bool)` - Retrieve metadata from context

## Examples

### Basic Cache Operations

```go
ctx := context.Background()
cache := synapse.New[string, int]()

cache.Set(ctx, "key1", 42)
if value, found := cache.Get(ctx, "key1"); found {
    fmt.Println("Value:", value)
}

cache.Delete(ctx, "key1")
size := cache.Len()
```

### Similarity Search

```go
cache := synapse.New[string, string]()
cache.WithSimilarity(func(a, b string) float64 {
    // Return similarity score between 0.0 and 1.0
    return computeSimilarity(a, b)
})

cache.Set(ctx, "apple", "A fruit")
cache.Set(ctx, "application", "A software program")

value, matchedKey, score, found := cache.GetSimilar(ctx, "app")
if found {
    fmt.Printf("Found: %s (matched: %s, similarity: %.2f)\n", value, matchedKey, score)
}
```

### Namespace Isolation

```go
cache := synapse.New[string, string]()

ctx1 := synapse.WithNamespace(context.Background(), "tenant1")
ctx2 := synapse.WithNamespace(context.Background(), "tenant2")

cache.Set(ctx1, "config", "tenant1's config")
cache.Set(ctx2, "config", "tenant2's config")

if value, found := cache.Get(ctx1, "config"); found {
    fmt.Println(value) // Output: tenant1's config
}
```

### TTL Expiration

```go
cache := synapse.New[string, string](
    synapse.WithTTL(5 * time.Minute),
)

cache.Set(ctx, "temp-key", "temporary value")
// Entry expires after 5 minutes
```

### Eviction Policy

```go
import "github.com/kolosys/synapse/eviction"

lru := eviction.NewLRU(1000)
cache := synapse.New[string, string](
    synapse.WithMaxSize(1000),
    synapse.WithEviction(lru),
)
```

## Architecture

Synapse uses sharding to distribute keys across multiple partitions, reducing lock contention and improving concurrent performance. Each shard operates independently with its own:

- `sync.RWMutex` for thread-safe access
- Hash map for O(1) exact lookups
- Key slice for similarity searches
- Eviction policy tracker

Exact lookups (`Get`) route to a single shard using FNV-1a hashing. Similarity searches (`GetSimilar`) search across all shards sequentially, respecting context cancellation.

## Performance

- **Exact lookups**: O(1) average case per shard
- **Similarity search**: O(n) per shard where n is the number of keys
- **Sharding**: More shards improve concurrency but increase overhead
- **Recommendation**: Start with 16 shards, adjust based on workload

Each cache entry stores the key, value, timestamps, access count, expiration time, metadata map, and namespace string. Plan capacity based on your key/value sizes.

## License

MIT License - see [LICENSE](LICENSE) file for details.
