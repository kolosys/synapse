# Quick Start

This guide will help you get started with Synapse quickly with practical examples.

## Basic Cache Operations

Create a cache and perform basic CRUD operations:

```go
package main

import (
    "context"
    "fmt"

    "github.com/kolosys/synapse"
)

func main() {
    ctx := context.Background()

    // Create a new cache with default settings
    cache := synapse.New[string, int]()

    // Set values
    cache.Set(ctx, "apples", 10)
    cache.Set(ctx, "oranges", 5)

    // Get value by exact key
    if value, found := cache.Get(ctx, "apples"); found {
        fmt.Println("Apples:", value) // Output: Apples: 10
    }

    // Delete a key
    cache.Delete(ctx, "oranges")

    // Get cache size
    fmt.Println("Size:", cache.Len()) // Output: Size: 1
}
```

## Similarity-Based Lookups

The core feature of Synapse is finding similar keys when exact matches don't exist:

```go
package main

import (
    "context"
    "fmt"

    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/algorithms"
)

func main() {
    ctx := context.Background()

    // Create cache with similarity threshold
    cache := synapse.New[string, string](
        synapse.WithThreshold(0.7), // 70% similarity required
    )

    // Set similarity function
    cache.WithSimilarity(algorithms.Levenshtein)

    // Populate cache
    cache.Set(ctx, "user:alice", "Alice's profile")
    cache.Set(ctx, "user:bob", "Bob's profile")
    cache.Set(ctx, "user:charlie", "Charlie's profile")

    // Exact match works as normal
    if value, found := cache.Get(ctx, "user:alice"); found {
        fmt.Println("Exact:", value)
    }

    // Find similar key (typo: "user:alic" instead of "user:alice")
    value, matchedKey, score, found := cache.GetSimilar(ctx, "user:alic")
    if found {
        fmt.Printf("Found: %s\n", value)
        fmt.Printf("Matched key: %s\n", matchedKey)
        fmt.Printf("Similarity: %.2f\n", score)
    }
    // Output:
    // Found: Alice's profile
    // Matched key: user:alice
    // Similarity: 0.91
}
```

## Configuration Options

Synapse uses the functional options pattern for configuration:

```go
cache := synapse.New[string, string](
    synapse.WithShards(16),           // Number of shards (1-256)
    synapse.WithMaxSize(10000),       // Maximum entries
    synapse.WithThreshold(0.8),       // Similarity threshold (0.0-1.0)
    synapse.WithTTL(5 * time.Minute), // Entry expiration
    synapse.WithStats(true),          // Enable statistics
)
```

| Option             | Default | Description                          |
| ------------------ | ------- | ------------------------------------ |
| `WithShards(n)`    | 16      | Number of cache shards (1-256)       |
| `WithMaxSize(n)`   | 1000    | Maximum entries across all shards    |
| `WithThreshold(t)` | 0.8     | Minimum similarity score for matches |
| `WithEviction(p)`  | nil     | Eviction policy (e.g., LRU)          |
| `WithTTL(d)`       | 0       | Time-to-live for entries             |
| `WithStats(b)`     | false   | Enable statistics tracking           |

## Using Eviction Policies

Configure LRU eviction to manage cache size:

```go
package main

import (
    "context"
    "fmt"

    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/eviction"
)

func main() {
    ctx := context.Background()

    // Create LRU policy
    lru := eviction.NewLRU(100)

    // Create cache with LRU eviction
    cache := synapse.New[string, string](
        synapse.WithMaxSize(100),
        synapse.WithEviction(lru),
    )

    // Add entries - LRU will evict oldest when full
    for i := 0; i < 150; i++ {
        key := fmt.Sprintf("key-%d", i)
        cache.Set(ctx, key, "value")
    }

    fmt.Println("Size:", cache.Len()) // Output: Size: 100
}
```

## Namespace Isolation

Partition cache entries by namespace for multi-tenant applications:

```go
package main

import (
    "context"
    "fmt"

    "github.com/kolosys/synapse"
)

func main() {
    cache := synapse.New[string, string]()

    // Create contexts with different namespaces
    tenant1 := synapse.WithNamespace(context.Background(), "tenant-1")
    tenant2 := synapse.WithNamespace(context.Background(), "tenant-2")

    // Store same key in different namespaces
    cache.Set(tenant1, "config", "tenant 1 config")
    cache.Set(tenant2, "config", "tenant 2 config")

    // Retrieve from specific namespace
    if value, found := cache.Get(tenant1, "config"); found {
        fmt.Println("Tenant 1:", value) // Output: Tenant 1: tenant 1 config
    }

    if value, found := cache.Get(tenant2, "config"); found {
        fmt.Println("Tenant 2:", value) // Output: Tenant 2: tenant 2 config
    }
}
```

## TTL Expiration

Automatically expire cache entries after a duration:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/kolosys/synapse"
)

func main() {
    ctx := context.Background()

    cache := synapse.New[string, string](
        synapse.WithTTL(1 * time.Second),
    )

    cache.Set(ctx, "temp", "temporary value")

    // Value exists immediately
    if value, found := cache.Get(ctx, "temp"); found {
        fmt.Println("Before:", value)
    }

    // Wait for expiration
    time.Sleep(2 * time.Second)

    // Value is expired
    if _, found := cache.Get(ctx, "temp"); !found {
        fmt.Println("After: expired")
    }
}
```

## Custom Similarity Functions

Define your own similarity function for domain-specific matching:

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/kolosys/synapse"
)

// prefixSimilarity compares string prefixes
func prefixSimilarity(a, b string) float64 {
    a, b = strings.ToLower(a), strings.ToLower(b)
    if a == b {
        return 1.0
    }

    minLen := len(a)
    if len(b) < minLen {
        minLen = len(b)
    }

    matches := 0
    for i := 0; i < minLen; i++ {
        if a[i] == b[i] {
            matches++
        } else {
            break
        }
    }

    maxLen := len(a)
    if len(b) > maxLen {
        maxLen = len(b)
    }

    return float64(matches) / float64(maxLen)
}

func main() {
    ctx := context.Background()

    cache := synapse.New[string, string](
        synapse.WithThreshold(0.5),
    )
    cache.WithSimilarity(prefixSimilarity)

    cache.Set(ctx, "application", "An app")
    cache.Set(ctx, "banana", "A fruit")

    // "app" matches "application" by prefix
    value, key, score, found := cache.GetSimilar(ctx, "app")
    if found {
        fmt.Printf("Found: %s (key: %s, score: %.2f)\n", value, key, score)
        // Output: Found: An app (key: application, score: 0.27)
    }
}
```

## Using Built-in Algorithms

Synapse provides several built-in similarity algorithms:

```go
package main

import (
    "fmt"

    "github.com/kolosys/synapse/algorithms"
)

func main() {
    // Levenshtein distance (edit distance)
    score := algorithms.Levenshtein("hello", "helo")
    fmt.Printf("Levenshtein: %.2f\n", score) // 0.80

    // Damerau-Levenshtein (with transpositions)
    score = algorithms.DamerauLevenshtein("hello", "hlelo")
    fmt.Printf("Damerau-Levenshtein: %.2f\n", score) // 0.80

    // Hamming distance (same length strings)
    score = algorithms.Hamming("hello", "hallo")
    fmt.Printf("Hamming: %.2f\n", score) // 0.80

    // Euclidean distance (vectors)
    score = algorithms.Euclidean([]float64{1, 2, 3}, []float64{1, 2, 4})
    fmt.Printf("Euclidean: %.2f\n", score) // 0.50

    // Manhattan distance (vectors)
    score = algorithms.Manhattan([]float64{1, 2, 3}, []float64{1, 2, 4})
    fmt.Printf("Manhattan: %.2f\n", score) // 0.50
}
```

## Context Cancellation

Synapse respects context cancellation for all operations:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/algorithms"
)

func main() {
    cache := synapse.New[string, string]()
    cache.WithSimilarity(algorithms.Levenshtein)

    // Populate with many entries
    ctx := context.Background()
    for i := 0; i < 10000; i++ {
        cache.Set(ctx, fmt.Sprintf("key-%d", i), "value")
    }

    // Create context with timeout
    searchCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
    defer cancel()

    // Similarity search respects cancellation
    _, _, _, found := cache.GetSimilar(searchCtx, "search-term")
    if !found {
        fmt.Println("Search cancelled or no match found")
    }
}
```

## Next Steps

- [Core Concepts](../core-concepts/synapse.md) - Understand the architecture
- [Algorithms](../core-concepts/algorithms.md) - Learn about similarity algorithms
- [Eviction Policies](../core-concepts/eviction.md) - Configure cache eviction
- [Best Practices](../advanced/best-practices.md) - Production recommendations
- [Performance Tuning](../advanced/performance-tuning.md) - Optimize for your workload
