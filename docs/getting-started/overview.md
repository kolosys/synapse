# Overview

Synapse is a high-performance, generic similarity-based cache for Go with intelligent sharding and pluggable eviction policies.

## What is Synapse?

Synapse is a thread-safe, context-aware cache that goes beyond traditional key-value storage by supporting **similarity-based lookups**. When an exact key match isn't found, Synapse can find the "closest" matching key based on a configurable similarity function and threshold.

This makes Synapse ideal for:

- **Fuzzy matching** - Find cached results for similar queries
- **Typo tolerance** - Match user input despite spelling errors
- **Semantic caching** - Cache responses for semantically similar requests
- **Deduplication** - Identify near-duplicate entries

## Key Features

| Feature                      | Description                                                           |
| ---------------------------- | --------------------------------------------------------------------- |
| **Similarity-Based Lookups** | Find approximate matches when exact keys don't exist                  |
| **Generic Types**            | Fully type-safe with Go generics (1.24+)                              |
| **High Performance**         | Automatic sharding distributes load across concurrent-safe partitions |
| **Pluggable Similarity**     | Define custom similarity algorithms for your use case                 |
| **Eviction Policies**        | LRU with support for combined policies                                |
| **TTL Support**              | Automatic expiration of cache entries                                 |
| **Namespace Isolation**      | Partition cache entries by namespace via context                      |
| **Thread-Safe**              | RWMutex-protected reads and writes per shard                          |
| **Metadata Support**         | Attach custom metadata to cache entries                               |
| **Context-Aware**            | Full `context.Context` integration for cancellation                   |
| **Zero Dependencies**        | Only uses Go standard library                                         |

## Project Information

- **Repository**: [github.com/kolosys/synapse](https://github.com/kolosys/synapse)
- **Import Path**: `github.com/kolosys/synapse`
- **License**: MIT
- **Go Version**: 1.24+

## Package Structure

Synapse is organized into three packages:

### `synapse` (main package)

The core cache implementation with:

- `Cache[K, V]` - Generic similarity-based cache
- `Entry[K, V]` - Cache entry with metadata
- `Shard[K, V]` - Individual cache shard
- Configuration via functional options
- Context utilities for namespace and metadata

### `synapse/algorithms`

Built-in similarity algorithms:

- **Levenshtein** - Edit distance for strings
- **Damerau-Levenshtein** - Edit distance with transpositions
- **Hamming** - Hamming distance for equal-length data
- **Euclidean** - Euclidean distance for vectors
- **Manhattan** - Manhattan distance for vectors

### `synapse/eviction`

Pluggable eviction policies:

- **LRU** - Least Recently Used eviction
- **CombinedPolicy** - Combine multiple policies with weighted scoring
- **EvictionPolicy** - Interface for custom implementations

## Documentation Structure

| Section                                | Description                           |
| -------------------------------------- | ------------------------------------- |
| [Getting Started](../getting-started/) | Installation and quick start guides   |
| [Core Concepts](../core-concepts/)     | Architecture and fundamental concepts |
| [Advanced Topics](../advanced/)        | Performance tuning and best practices |
| [API Reference](../api-reference/)     | Complete API documentation            |
| [Examples](../examples/)               | Working code examples                 |

## Quick Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/algorithms"
    "github.com/kolosys/synapse/eviction"
)

func main() {
    ctx := context.Background()

    // Create cache with LRU eviction
    cache := synapse.New[string, string](
        synapse.WithMaxSize(1000),
        synapse.WithThreshold(0.7),
        synapse.WithEviction(eviction.NewLRU(1000)),
    )

    // Set similarity function using Levenshtein distance
    cache.WithSimilarity(algorithms.Levenshtein)

    // Store values
    cache.Set(ctx, "hello world", "greeting")
    cache.Set(ctx, "goodbye world", "farewell")

    // Exact match
    if value, found := cache.Get(ctx, "hello world"); found {
        fmt.Println("Exact:", value) // Output: Exact: greeting
    }

    // Similarity search
    value, key, score, found := cache.GetSimilar(ctx, "hello wrold")
    if found {
        fmt.Printf("Similar: %s (key: %s, score: %.2f)\n", value, key, score)
        // Output: Similar: greeting (key: hello world, score: 0.91)
    }
}
```

## Community & Support

- **GitHub Issues**: [github.com/kolosys/synapse/issues](https://github.com/kolosys/synapse/issues)
- **Discussions**: [github.com/kolosys/synapse/discussions](https://github.com/kolosys/synapse/discussions)

## Next Steps

- [Installation Guide](installation.md) - Install Synapse
- [Quick Start](quick-start.md) - Get started with examples
- [Core Concepts](../core-concepts/synapse.md) - Understand the architecture
