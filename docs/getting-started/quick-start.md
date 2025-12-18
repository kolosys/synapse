# Quick Start

This guide will help you get started with synapse quickly with a basic example.

## Basic Usage

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/synapse/algorithms"
    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/eviction"
)

func main() {
    // Basic usage example
    fmt.Println("Welcome to synapse!")
    
    // TODO: Add your code here
}
```

## Common Use Cases

### Using algorithms

**Import Path:** `github.com/kolosys/synapse/algorithms`



```go
package main

import (
    "fmt"
    "github.com/kolosys/synapse/algorithms"
)

func main() {
    // Example usage of algorithms
    fmt.Println("Using algorithms package")
}
```

#### Available Functions
- **DamerauLevenshtein** - DamerauLevenshtein computes the Damerau-Levenshtein distance Similar to Levenshtein but also allows transpositions
- **Euclidean** - Euclidean computes the Euclidean distance between two points Returns a similarity score (inverse of distance)
- **Hamming** - Hamming computes the Hamming distance between two strings Strings must be of equal length. Returns a normalized score between 0.0 and 1.0
- **HammingBytes** - HammingBytes computes the Hamming distance between two byte slices
- **Levenshtein** - Levenshtein computes the Levenshtein distance between two strings Returns a normalized score between 0.0 (completely different) and 1.0 (identical)
- **Manhattan** - Manhattan computes the Manhattan distance between two points

For detailed API documentation, see the [algorithms API Reference](../api-reference/algorithms.md).

### Using synapse

**Import Path:** `github.com/kolosys/synapse`



```go
package main

import (
    "fmt"
    "github.com/kolosys/synapse"
)

func main() {
    // Example usage of synapse
    fmt.Println("Using synapse package")
}
```

#### Available Types
- **Cache** - Cache is a generic similarity-based cache with sharding
- **Entry** - Entry represents a cache entry with metadata
- **EvictionPolicy** - EvictionPolicy is re-exported from the eviction package
- **Option** - Option is a function that modifies Options
- **Options** - Options contains configuration options for the cache
- **Shard** - Shard represents a single shard of the cache
- **Similarity** - Similarity is an interface for similarity computation
- **SimilarityFunc** - SimilarityFunc is a function type that computes similarity between two keys It should return a score between 0.0 (completely different) and 1.0 (identical)

#### Available Functions
- **GetMetadata** - GetMetadata retrieves a metadata value from the context
- **GetNamespace** - GetNamespace retrieves the namespace from the context
- **WithMetadata** - WithMetadata adds metadata to the context
- **WithNamespace** - WithNamespace adds a namespace to the context

For detailed API documentation, see the [synapse API Reference](../api-reference/synapse.md).

### Using eviction

**Import Path:** `github.com/kolosys/synapse/eviction`



```go
package main

import (
    "fmt"
    "github.com/kolosys/synapse/eviction"
)

func main() {
    // Example usage of eviction
    fmt.Println("Using eviction package")
}
```

#### Available Types
- **CombinedPolicy** - CombinedPolicy combines multiple eviction policies with weighted scoring
- **EvictionPolicy** - EvictionPolicy defines the interface for cache eviction strategies
- **LRU** - LRU implements a Least Recently Used eviction policy

For detailed API documentation, see the [eviction API Reference](../api-reference/eviction.md).

## Step-by-Step Tutorial

### Step 1: Import the Package

First, import the necessary packages in your Go file:

```go
import (
    "fmt"
    "github.com/kolosys/synapse/algorithms"
    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/eviction"
)
```

### Step 2: Initialize

Set up the basic configuration:

```go
func main() {
    // Initialize your application
    fmt.Println("Initializing synapse...")
}
```

### Step 3: Use the Library

Implement your specific use case:

```go
func main() {
    // Your implementation here
}
```

## Running Your Code

To run your Go program:

```bash
go run main.go
```

To build an executable:

```bash
go build -o myapp
./myapp
```

## Configuration Options

synapse can be configured to suit your needs. Check the [Core Concepts](../core-concepts/) section for detailed information about configuration options.

## Error Handling

Always handle errors appropriately:

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

## Best Practices

- Always handle errors returned by library functions
- Check the API documentation for detailed parameter information
- Use meaningful variable and function names
- Add comments to document your code

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/synapse/algorithms"
    "github.com/kolosys/synapse"
    "github.com/kolosys/synapse/eviction"
)

func main() {
    fmt.Println("Starting synapse application...")
    
    // Add your implementation here
    
    fmt.Println("Application completed successfully!")
}
```

## Next Steps

Now that you've seen the basics, explore:

- **[Core Concepts](../core-concepts/)** - Understanding the library architecture
- **[API Reference](../api-reference/)** - Complete API documentation
- **[Examples](../examples/README.md)** - More detailed examples
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced patterns

## Getting Help

If you run into issues:

1. Check the [API Reference](../api-reference/)
2. Browse the [Examples](../examples/README.md)
3. Visit the [GitHub Issues](https://github.com/kolosys/synapse/issues) page

