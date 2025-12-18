# Eviction Policies

The `eviction` package provides pluggable cache eviction strategies for managing cache size and memory usage.

**Import Path:** `github.com/kolosys/synapse/eviction`

## Overview

When a cache reaches its maximum size, it must decide which entries to remove (evict) to make room for new ones. Synapse supports pluggable eviction policies that determine which entries are selected for removal.

## The EvictionPolicy Interface

All eviction policies implement the `EvictionPolicy` interface:

```go
type EvictionPolicy interface {
    // OnAccess is called when an entry is accessed
    OnAccess(key any)

    // OnAdd is called when an entry is added
    OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time)

    // OnRemove is called when an entry is removed
    OnRemove(key any)

    // SelectVictim returns the key of the entry to evict
    SelectVictim() (any, bool)

    // Len returns the number of tracked entries
    Len() int
}
```

The cache calls these methods automatically:

- `OnAdd` when a new entry is stored
- `OnAccess` when an entry is retrieved
- `OnRemove` when an entry is deleted or evicted
- `SelectVictim` when the cache needs to make room

## Built-in Policies

### LRU (Least Recently Used)

The LRU policy evicts the entry that hasn't been accessed for the longest time. This is effective when recently accessed items are more likely to be accessed again.

```go
import "github.com/kolosys/synapse/eviction"

// Create LRU policy with capacity
lru := eviction.NewLRU(1000)

// Use with cache
cache := synapse.New[string, string](
    synapse.WithMaxSize(1000),
    synapse.WithEviction(lru),
)
```

**How it works:**

1. Maintains a doubly-linked list of keys
2. Most recently accessed keys move to the front
3. Least recently accessed keys are at the back
4. Eviction selects the key at the back

**Complexity:**

- `OnAccess`: O(1)
- `OnAdd`: O(1)
- `OnRemove`: O(1)
- `SelectVictim`: O(1)

**Best for:**

- General-purpose caching
- Workloads with temporal locality
- When recent access predicts future access

### CombinedPolicy

Combines multiple eviction policies with weighted scoring. This allows sophisticated eviction strategies that consider multiple factors.

```go
import "github.com/kolosys/synapse/eviction"

// Create individual policies
lru := eviction.NewLRU(1000)
lfu := eviction.NewLFU(1000) // If implemented

// Combine with weights (60% LRU, 40% LFU)
combined := eviction.NewCombinedPolicy(
    []eviction.EvictionPolicy{lru, lfu},
    []float64{0.6, 0.4},
)

cache := synapse.New[string, string](
    synapse.WithMaxSize(1000),
    synapse.WithEviction(combined),
)
```

**Note:** Currently, `SelectVictim` uses the first policy's selection. The combined policy primarily ensures all policies receive access notifications for accurate tracking.

## Using Eviction Policies

### Basic Configuration

```go
ctx := context.Background()

lru := eviction.NewLRU(100)
cache := synapse.New[string, string](
    synapse.WithMaxSize(100),
    synapse.WithEviction(lru),
)

// Add entries
for i := 0; i < 150; i++ {
    key := fmt.Sprintf("key-%d", i)
    cache.Set(ctx, key, "value")
}

// Only 100 entries remain (oldest 50 were evicted)
fmt.Println(cache.Len()) // 100
```

### Without Eviction Policy

If no eviction policy is set, the cache uses FIFO (First In, First Out) as a fallback:

```go
cache := synapse.New[string, string](
    synapse.WithMaxSize(100),
    // No eviction policy set
)

// When full, oldest entries are removed first
```

### Eviction with TTL

Eviction policies work alongside TTL. Expired entries are not selected for eviction—they're simply skipped during access:

```go
cache := synapse.New[string, string](
    synapse.WithMaxSize(100),
    synapse.WithTTL(5 * time.Minute),
    synapse.WithEviction(eviction.NewLRU(100)),
)
```

## Implementing Custom Policies

Create custom eviction policies by implementing the `EvictionPolicy` interface:

```go
package main

import (
    "sync"
    "time"

    "github.com/kolosys/synapse/eviction"
)

// FIFO implements First-In-First-Out eviction
type FIFO struct {
    mu    sync.Mutex
    order []any
}

func NewFIFO() *FIFO {
    return &FIFO{
        order: make([]any, 0),
    }
}

func (f *FIFO) OnAccess(key any) {
    // FIFO doesn't track access
}

func (f *FIFO) OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.order = append(f.order, key)
}

func (f *FIFO) OnRemove(key any) {
    f.mu.Lock()
    defer f.mu.Unlock()
    for i, k := range f.order {
        if k == key {
            f.order = append(f.order[:i], f.order[i+1:]...)
            break
        }
    }
}

func (f *FIFO) SelectVictim() (any, bool) {
    f.mu.Lock()
    defer f.mu.Unlock()
    if len(f.order) == 0 {
        return nil, false
    }
    return f.order[0], true
}

func (f *FIFO) Len() int {
    f.mu.Lock()
    defer f.mu.Unlock()
    return len(f.order)
}
```

Usage:

```go
fifo := NewFIFO()
cache := synapse.New[string, string](
    synapse.WithMaxSize(100),
    synapse.WithEviction(fifo),
)
```

## Architecture

### Per-Shard Eviction

Each shard maintains its own eviction policy state. When you configure eviction:

```go
lru := eviction.NewLRU(1000)
cache := synapse.New[string, string](
    synapse.WithMaxSize(1000),
    synapse.WithShards(16),
    synapse.WithEviction(lru),
)
```

The same policy instance is shared across shards, but the `maxSize` is divided:

- Total max size: 1000
- Per-shard max size: 1000 / 16 = 62

### Thread Safety

The built-in LRU policy uses `sync.RWMutex` for thread-safe operations. Custom policies must also be thread-safe since multiple shards may call policy methods concurrently.

### Entry Metadata

The `OnAdd` method receives metadata that policies can use for decisions:

- `accessCount` - Number of times the entry has been accessed
- `createdAt` - When the entry was created
- `accessedAt` - When the entry was last accessed

This enables sophisticated policies like LFU (Least Frequently Used) or time-weighted scoring.

## Policy Selection Guide

| Policy   | Best For                           | Avoid When                       |
| -------- | ---------------------------------- | -------------------------------- |
| LRU      | General caching, temporal locality | Random access patterns           |
| FIFO     | Simple, predictable eviction       | Access patterns matter           |
| LFU      | Frequency-based caching            | Access patterns change over time |
| Combined | Complex requirements               | Simplicity is preferred          |

## Performance Considerations

1. **O(1) operations** - LRU achieves constant-time operations using a hash map + doubly-linked list
2. **Memory overhead** - Each policy adds tracking overhead per entry
3. **Lock contention** - High-throughput workloads may benefit from sharding

## Further Reading

- [API Reference](../api-reference/eviction.md) - Complete API documentation
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization strategies
- [Core Concepts: Synapse](synapse.md) - Cache architecture
