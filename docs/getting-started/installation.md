# Installation

This guide will help you install and set up Synapse in your Go project.

## Prerequisites

Before installing Synapse, ensure you have:

- **Go 1.24+** installed
- A Go module initialized in your project (run `go mod init` if needed)

## Installation

Use `go get` to install Synapse:

```bash
go get github.com/kolosys/synapse
```

This will download the package and add it to your `go.mod` file.

## Available Packages

Synapse includes several packages that you can import based on your needs:

### Core Package

The main cache functionality:

```go
import "github.com/kolosys/synapse"
```

### Algorithms Package

Built-in similarity algorithms for string and vector comparisons:

```go
import "github.com/kolosys/synapse/algorithms"
```

Provides:

- `Levenshtein` - Edit distance for strings
- `DamerauLevenshtein` - Edit distance with transpositions
- `Hamming` / `HammingBytes` - Hamming distance for equal-length strings/bytes
- `Euclidean` - Euclidean distance for vectors
- `Manhattan` - Manhattan distance for vectors

### Eviction Package

Cache eviction policies:

```go
import "github.com/kolosys/synapse/eviction"
```

Provides:

- `LRU` - Least Recently Used eviction
- `CombinedPolicy` - Combine multiple policies with weighted scoring

## Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "context"
    "fmt"

    "github.com/kolosys/synapse"
)

func main() {
    cache := synapse.New[string, string]()
    ctx := context.Background()

    cache.Set(ctx, "hello", "world")
    if value, found := cache.Get(ctx, "hello"); found {
        fmt.Println("Synapse installed successfully:", value)
    }
}
```

Run the test:

```bash
go run main.go
```

## Updating

To update to the latest version:

```bash
go get -u github.com/kolosys/synapse
```

To update to a specific version:

```bash
go get github.com/kolosys/synapse@v1.2.3
```

## Development Setup

If you want to contribute or modify the library:

```bash
git clone https://github.com/kolosys/synapse.git
cd synapse
go mod download
go test -race ./...
```

## Troubleshooting

### Module Not Found

If you encounter a "module not found" error:

1. Ensure your `GOPATH` is set correctly
2. Check that you have network access to GitHub
3. Try running `go clean -modcache` and reinstall

### Private Repository Access

For private repositories, configure Git to use SSH or a personal access token:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Or set up GOPRIVATE:

```bash
export GOPRIVATE=github.com/kolosys/*
```

## Next Steps

- [Quick Start Guide](quick-start.md) - Learn how to use Synapse
- [Core Concepts](../core-concepts/synapse.md) - Understand the architecture
- [API Reference](../api-reference/synapse.md) - Complete API documentation
