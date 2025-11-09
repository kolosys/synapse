# Synapse

A high-performance, generic similarity-based cache for Go with intelligent sharding and pluggable eviction policies.

[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/doc/go1.24)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

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

// Simple string similarity function (case-insensitive prefix match)
func stringSimilarity(a, b string) float64 {
    a, b = strings.ToLower(a), strings.ToLower(b)
    if a == b {
        return 1.0
    }
    // Simple prefix matching
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
    
    // Create a cache with custom options
    cache := synapse.New[string, string](
        synapse.WithMaxSize(1000),
        synapse.WithShards(16),
        synapse.WithThreshold(0.7),
        synapse.WithEviction(eviction.NewLRU(1000)),
    )
    
    // Set the similarity function
    cache.WithSimilarity(stringSimilarity)
    
    // Store some values
    cache.Set(ctx, "user:alice", "Alice's data")
    cache.Set(ctx, "user:bob", "Bob's data")
    cache.Set(ctx, "user:charlie", "Charlie's data")
    
    // Exact match
    if value, found := cache.Get(ctx, "user:alice"); found {
        fmt.Println("Exact match:", value)
    }
    
    // Similarity-based match
    value, key, score, found := cache.GetSimilar(ctx, "user:ali")
    if found {
        fmt.Printf("Similar match: %s (key: %s, score: %.2f)\n", value, key, score)
        // Output: Similar match: Alice's data (key: user:alice, score: 0.82)
    }
}
```

## Usage Examples

### Basic Operations

```go
ctx := context.Background()
cache := synapse.New[string, int]()

// Set a value
cache.Set(ctx, "key1", 42)

// Get a value
if value, found := cache.Get(ctx, "key1"); found {
    fmt.Println("Value:", value)
}

// Delete a value
cache.Delete(ctx, "key1")

// Get cache size
size := cache.Len()
```

### Similarity-Based Search

```go
// Create cache with similarity function
cache := synapse.New[string, string]()
cache.WithSimilarity(func(a, b string) float64 {
    // Your custom similarity algorithm
    // Return 0.0 (completely different) to 1.0 (identical)
    return computeSimilarity(a, b)
})

// Store data
cache.Set(ctx, "apple", "A fruit")
cache.Set(ctx, "application", "A software program")

// Find similar key
value, matchedKey, score, found := cache.GetSimilar(ctx, "app")
if found {
    fmt.Printf("Found: %s (matched: %s, similarity: %.2f)\n", 
        value, matchedKey, score)
}
```

### Using Namespaces

Namespaces allow you to isolate cache entries:

```go
// Create cache
cache := synapse.New[string, string]()

// Store in different namespaces
ctx1 := synapse.WithNamespace(context.Background(), "tenant1")
ctx2 := synapse.WithNamespace(context.Background(), "tenant2")

cache.Set(ctx1, "config", "tenant1's config")
cache.Set(ctx2, "config", "tenant2's config")

// Retrieve from specific namespace
if value, found := cache.Get(ctx1, "config"); found {
    fmt.Println(value) // Output: tenant1's config
}
```

### TTL (Time-To-Live)

```go
import "time"

// Cache with 5-minute TTL
cache := synapse.New[string, string](
    synapse.WithTTL(5 * time.Minute),
)

cache.Set(ctx, "temp-key", "temporary value")

// After 5 minutes, the entry will be automatically expired
time.Sleep(6 * time.Minute)
_, found := cache.Get(ctx, "temp-key") // found == false
```

### Custom Eviction Policy

```go
import "github.com/kolosys/synapse/eviction"

// Create LRU eviction policy
lru := eviction.NewLRU(1000)

// Use in cache
cache := synapse.New[string, string](
    synapse.WithMaxSize(1000),
    synapse.WithEviction(lru),
)
```

## Configuration Options

### Cache Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithShards(n)` | Number of shards (1-256) | 16 |
| `WithMaxSize(size)` | Maximum number of entries | 1000 |
| `WithThreshold(t)` | Similarity threshold (0.0-1.0) | 0.8 |
| `WithEviction(policy)` | Eviction policy | nil |
| `WithTTL(duration)` | Time-to-live for entries | 0 (no expiration) |
| `WithStats(enable)` | Enable statistics tracking | false |

### Example Configuration

```go
cache := synapse.New[string, []byte](
    synapse.WithShards(32),           // 32 shards for high concurrency
    synapse.WithMaxSize(10000),       // Store up to 10k entries
    synapse.WithThreshold(0.85),      // 85% similarity required
    synapse.WithTTL(10*time.Minute),  // 10-minute expiration
    synapse.WithEviction(eviction.NewLRU(10000)),
)
```

## Architecture

### Sharding

Synapse uses consistent hashing to distribute keys across multiple shards, reducing lock contention and improving concurrent performance:

```
Cache
├── Shard 0 (mutex, map, eviction)
├── Shard 1 (mutex, map, eviction)
├── Shard 2 (mutex, map, eviction)
└── ...
```

Each shard operates independently with its own:
- Read/Write mutex for thread safety
- HashMap for O(1) exact lookups
- Key slice for similarity searches
- Eviction policy tracker

### Similarity Search

For exact key lookups (`Get`), Synapse routes to the appropriate shard using FNV-1a hashing.

For similarity searches (`GetSimilar`), Synapse:
1. Searches across **all shards** in parallel
2. Computes similarity scores for each key
3. Returns the best match above the threshold
4. Respects context cancellation for long-running searches

## Roadmap

### Upcoming Features

#### Eviction Policies
- 🔄 **LFU (Least Frequently Used)** - Evict entries with lowest access count
- ⏰ **TTL-based** - Advanced time-based eviction strategies
- 🎯 **Adaptive policies** - Combine multiple strategies

#### Similarity Algorithms
- 📝 **Levenshtein Distance** - Edit distance for string matching
- 🎲 **MinHash** - Fast approximate similarity for sets
- 📊 **Jaccard Similarity** - Set-based similarity coefficient  
- 🧮 **Cosine Similarity** - Vector-based similarity for embeddings

#### Additional Features
- 📈 **Comprehensive benchmarks** - Performance testing suite
- 🧪 **Advanced examples** - Real-world use case demonstrations
- 📚 **Extended documentation** - In-depth guides and tutorials
- ✅ **Expanded test coverage** - Comprehensive testing suite
- 📊 **Statistics & monitoring** - Cache hit rates, eviction metrics
- 💾 **Persistence options** - Snapshot and restore capabilities
- 🔍 **Advanced indexing** - LSH and other approximate nearest neighbor techniques

## Performance Considerations

### Sharding

The number of shards affects concurrent performance:
- **More shards**: Better concurrency, but higher overhead
- **Fewer shards**: Lower overhead, but potential lock contention
- **Recommendation**: Start with 16, increase for high-concurrency workloads

### Similarity Search

Similarity search is O(n) per shard, so:
- Use a **higher threshold** (e.g., 0.8-0.9) to exit early
- Implement **efficient similarity functions** (avoid complex computations)
- Use **context timeouts** for bounded search time
- Consider **exact lookups** when possible (O(1) vs O(n))

### Memory

Each cache entry stores:
- Key and value
- Timestamps (created, accessed)
- Access count
- Expiration time
- Metadata map
- Namespace string

Plan capacity accordingly based on your key/value sizes.

## Use Cases

### API Response Caching with Fuzzy Matching

Cache API responses and retrieve similar queries:

```go
cache := synapse.New[string, APIResponse]()
cache.WithSimilarity(queryStringSimilarity)

// Cache exact query
cache.Set(ctx, "weather?city=New+York&date=2024-01-15", response)

// Retrieve similar query (different date, same city)
resp, _, _, found := cache.GetSimilar(ctx, "weather?city=New+York&date=2024-01-16")
```

### Multi-Tenant Applications

Use namespaces to isolate tenant data:

```go
tenantCtx := synapse.WithNamespace(ctx, tenantID)
cache.Set(tenantCtx, "settings", tenantSettings)
```

### Machine Learning Embeddings

Store and retrieve similar embeddings:

```go
cache := synapse.New[[128]float32, ModelOutput]()
cache.WithSimilarity(cosineSimilarity)

// Find cached result for similar embedding
output, _, score, found := cache.GetSimilar(ctx, queryEmbedding)
```

## Contributing

Contributions are welcome! Areas of interest:
- Additional eviction policies
- Similarity algorithm implementations
- Performance optimizations
- Documentation improvements
- Test coverage expansion

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Project Status

⚠️ **Alpha/Development**: This is a basic version under active development. The API may change. 

Currently in progress:
- ✅ Core cache functionality
- ✅ LRU eviction
- ✅ Sharding
- ✅ Context support
- 🔄 Testing & benchmarking
- 🔄 Documentation
- 🔄 Examples
- 🔄 Additional eviction policies
- 🔄 Similarity algorithms

---

**Built with ❤️ by the Kolosys team**
