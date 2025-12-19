# synapse API

Complete API documentation for the synapse package.

**Import Path:** `github.com/kolosys/synapse`

## Package Documentation



## Types

### Cache
Cache is a generic similarity-based cache with sharding

#### Example Usage

```go
// Create a new Cache
cache := Cache{

}
```

#### Type Definition

```go
type Cache struct {
}
```

### Constructor Functions

### New

New creates a new cache with the given options

```go
func New(opts ...Option) **ast.IndexListExpr
```

**Parameters:**
- `opts` (...Option)

**Returns:**
- **ast.IndexListExpr

## Methods

### Delete

Delete removes a key from the cache

```go
func (**ast.IndexListExpr) Delete(ctx context.Context, key K) bool
```

**Parameters:**
- `ctx` (context.Context)
- `key` (K)

**Returns:**
- bool

### Get

Get retrieves a value by exact key match

```go
func (**ast.IndexListExpr) Get(ctx context.Context, key K) (V, bool)
```

**Parameters:**
- `ctx` (context.Context)
- `key` (K)

**Returns:**
- V
- bool

### GetSimilar

GetSimilar finds the most similar key above the threshold

```go
func (**ast.IndexListExpr) GetSimilar(ctx context.Context, key K) (V, K, float64, bool)
```

**Parameters:**
- `ctx` (context.Context)
- `key` (K)

**Returns:**
- V
- K
- float64
- bool

### Len

Len returns the total number of entries in the cache

```go
func (**ast.IndexListExpr) Len() int
```

**Parameters:**
  None

**Returns:**
- int

### Set

Set stores a value

```go
func (**ast.IndexListExpr) Set(ctx context.Context, key K, value V) error
```

**Parameters:**
- `ctx` (context.Context)
- `key` (K)
- `value` (V)

**Returns:**
- error

### Stats

Stats returns aggregated statistics from all shards Returns zero values if stats are not enabled

```go
func (**ast.IndexListExpr) Stats() Stats
```

**Parameters:**
  None

**Returns:**
- Stats

### WithSimilarity

WithSimilarity sets the similarity function for the cache

```go
func (**ast.IndexListExpr) WithSimilarity(fn *ast.IndexExpr) **ast.IndexListExpr
```

**Parameters:**
- `fn` (*ast.IndexExpr)

**Returns:**
- **ast.IndexListExpr

### Entry
Entry represents a cache entry with metadata

#### Example Usage

```go
// Create a new Entry
entry := Entry{
    Key: K{},
    Value: V{},
    CreatedAt: /* value */,
    AccessedAt: /* value */,
    AccessCount: 42,
    ExpiresAt: /* value */,
    Metadata: map[],
    Namespace: "example",
}
```

#### Type Definition

```go
type Entry struct {
    Key K
    Value V
    CreatedAt time.Time
    AccessedAt time.Time
    AccessCount uint64
    ExpiresAt time.Time
    Metadata map[string]any
    Namespace string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Key | `K` |  |
| Value | `V` |  |
| CreatedAt | `time.Time` |  |
| AccessedAt | `time.Time` |  |
| AccessCount | `uint64` |  |
| ExpiresAt | `time.Time` |  |
| Metadata | `map[string]any` |  |
| Namespace | `string` |  |

## Methods

### IsExpired

IsExpired checks if the entry has expired

```go
func (**ast.IndexListExpr) IsExpired() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Touch

Touch updates the access time and increments access count

```go
func (**ast.IndexListExpr) Touch()
```

**Parameters:**
  None

**Returns:**
  None

### EvictionPolicy
EvictionPolicy is re-exported from the eviction package

#### Example Usage

```go
// Example usage of EvictionPolicy
var value EvictionPolicy
// Initialize with appropriate value
```

#### Type Definition

```go
type EvictionPolicy eviction.EvictionPolicy
```

### Option
Option is a function that modifies Options

#### Example Usage

```go
// Example usage of Option
var value Option
// Initialize with appropriate value
```

#### Type Definition

```go
type Option func(*Options)
```

### Constructor Functions

### WithEviction

WithEviction sets the eviction policy

```go
func WithEviction(policy EvictionPolicy) Option
```

**Parameters:**
- `policy` (EvictionPolicy)

**Returns:**
- Option

### WithMaxSize

WithMaxSize sets the maximum cache size

```go
func WithMaxSize(size int) Option
```

**Parameters:**
- `size` (int)

**Returns:**
- Option

### WithShards

WithShards sets the number of shards

```go
func WithShards(n int) Option
```

**Parameters:**
- `n` (int)

**Returns:**
- Option

### WithStats

WithStats enables statistics tracking

```go
func WithStats(enable bool) Option
```

**Parameters:**
- `enable` (bool)

**Returns:**
- Option

### WithTTL

WithTTL sets the time-to-live for cache entries

```go
func WithTTL(ttl time.Duration) Option
```

**Parameters:**
- `ttl` (time.Duration)

**Returns:**
- Option

### WithThreshold

WithThreshold sets the similarity threshold

```go
func WithThreshold(t float64) Option
```

**Parameters:**
- `t` (float64)

**Returns:**
- Option

### Options
Options contains configuration options for the cache

#### Example Usage

```go
// Create a new Options
options := Options{
    NumShards: 42,
    MaxSize: 42,
    SimilarityThreshold: 3.14,
    EvictionPolicy: EvictionPolicy{},
    TTL: /* value */,
    EnableStats: true,
}
```

#### Type Definition

```go
type Options struct {
    NumShards int
    MaxSize int
    SimilarityThreshold float64
    EvictionPolicy EvictionPolicy
    TTL time.Duration
    EnableStats bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| NumShards | `int` |  |
| MaxSize | `int` |  |
| SimilarityThreshold | `float64` |  |
| EvictionPolicy | `EvictionPolicy` |  |
| TTL | `time.Duration` |  |
| EnableStats | `bool` |  |

### Shard
Shard represents a single shard of the cache

#### Example Usage

```go
// Create a new Shard
shard := Shard{

}
```

#### Type Definition

```go
type Shard struct {
}
```

### Similarity
Similarity is an interface for similarity computation

#### Example Usage

```go
// Example implementation of Similarity
type MySimilarity struct {
    // Add your fields here
}

func (m MySimilarity) Score(param1 K) float64 {
    // Implement your logic here
    return
}

func (m MySimilarity) Threshold() float64 {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Similarity interface {
    Score(a, b K) float64
    Threshold() float64
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Constructor Functions

### NewSimilarity

NewSimilarity creates a Similarity from a SimilarityFunc

```go
func NewSimilarity(fn *ast.IndexExpr, threshold float64) *ast.IndexExpr
```

**Parameters:**
- `fn` (*ast.IndexExpr)
- `threshold` (float64)

**Returns:**
- *ast.IndexExpr

### SimilarityFunc
SimilarityFunc is a function type that computes similarity between two keys It should return a score between 0.0 (completely different) and 1.0 (identical)

#### Example Usage

```go
// Example usage of SimilarityFunc
var value SimilarityFunc
// Initialize with appropriate value
```

#### Type Definition

```go
type SimilarityFunc func(a, b K) float64
```

### Stats
Stats contains cache performance statistics

#### Example Usage

```go
// Create a new Stats
stats := Stats{
    Hits: 42,
    Misses: 42,
    Sets: 42,
    Deletes: 42,
    SimilarSearches: 42,
    SimilarHits: 42,
    Evictions: 42,
    Expired: 42,
}
```

#### Type Definition

```go
type Stats struct {
    Hits uint64
    Misses uint64
    Sets uint64
    Deletes uint64
    SimilarSearches uint64
    SimilarHits uint64
    Evictions uint64
    Expired uint64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Hits | `uint64` |  |
| Misses | `uint64` |  |
| Sets | `uint64` |  |
| Deletes | `uint64` |  |
| SimilarSearches | `uint64` |  |
| SimilarHits | `uint64` |  |
| Evictions | `uint64` |  |
| Expired | `uint64` |  |

## Functions

### GetMetadata
GetMetadata retrieves a metadata value from the context

```go
func GetMetadata(ctx context.Context, key string) (any, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `key` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `any` | |
| `bool` | |

**Example:**

```go
// Example usage of GetMetadata
result := GetMetadata(/* parameters */)
```

### GetNamespace
GetNamespace retrieves the namespace from the context

```go
func GetNamespace(ctx context.Context) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of GetNamespace
result := GetNamespace(/* parameters */)
```

### WithMetadata
WithMetadata adds metadata to the context

```go
func WithMetadata(ctx context.Context, key string, value any) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `key` | `string` | |
| `value` | `any` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithMetadata
result := WithMetadata(/* parameters */)
```

### WithNamespace
WithNamespace adds a namespace to the context

```go
func WithNamespace(ctx context.Context, namespace string) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `namespace` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithNamespace
result := WithNamespace(/* parameters */)
```

## External Links

- [Package Overview](../packages/synapse.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/synapse)
- [Source Code](https://github.com/kolosys/synapse/tree/main/github.com/kolosys/synapse)
