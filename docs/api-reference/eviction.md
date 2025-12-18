# eviction API

Complete API documentation for the eviction package.

**Import Path:** `github.com/kolosys/synapse/eviction`

## Package Documentation



## Types

### CombinedPolicy
CombinedPolicy combines multiple eviction policies with weighted scoring

#### Example Usage

```go
// Create a new CombinedPolicy
combinedpolicy := CombinedPolicy{

}
```

#### Type Definition

```go
type CombinedPolicy struct {
}
```

### Constructor Functions

### NewCombinedPolicy

NewCombinedPolicy creates a new combined eviction policy

```go
func NewCombinedPolicy(policies []EvictionPolicy, weights []float64) *CombinedPolicy
```

**Parameters:**
- `policies` ([]EvictionPolicy)
- `weights` ([]float64)

**Returns:**
- *CombinedPolicy

## Methods

### Len

Len implements EvictionPolicy

```go
func (*LRU) Len() int
```

**Parameters:**
  None

**Returns:**
- int

### OnAccess

OnAccess implements EvictionPolicy

```go
func (*LRU) OnAccess(key any)
```

**Parameters:**
- `key` (any)

**Returns:**
  None

### OnAdd

OnAdd implements EvictionPolicy

```go
func (*LRU) OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time)
```

**Parameters:**
- `key` (any)
- `accessCount` (uint64)
- `createdAt` (time.Time)
- `accessedAt` (time.Time)

**Returns:**
  None

### OnRemove

OnRemove implements EvictionPolicy

```go
func (*LRU) OnRemove(key any)
```

**Parameters:**
- `key` (any)

**Returns:**
  None

### SelectVictim

SelectVictim implements EvictionPolicy It uses the first policy's victim selection

```go
func (*LRU) SelectVictim() (any, bool)
```

**Parameters:**
  None

**Returns:**
- any
- bool

### EvictionPolicy
EvictionPolicy defines the interface for cache eviction strategies

#### Example Usage

```go
// Example implementation of EvictionPolicy
type MyEvictionPolicy struct {
    // Add your fields here
}

func (m MyEvictionPolicy) OnAccess(param1 any)  {
    // Implement your logic here
    return
}

func (m MyEvictionPolicy) OnAdd(param1 any, param2 uint64, param3 time.Time)  {
    // Implement your logic here
    return
}

func (m MyEvictionPolicy) OnRemove(param1 any)  {
    // Implement your logic here
    return
}

func (m MyEvictionPolicy) SelectVictim() any {
    // Implement your logic here
    return
}

func (m MyEvictionPolicy) Len() int {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type EvictionPolicy interface {
    OnAccess(key any)
    OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time)
    OnRemove(key any)
    SelectVictim() (any, bool)
    Len() int
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### LRU
LRU implements a Least Recently Used eviction policy

#### Example Usage

```go
// Create a new LRU
lru := LRU{

}
```

#### Type Definition

```go
type LRU struct {
}
```

### Constructor Functions

### NewLRU

NewLRU creates a new LRU eviction policy

```go
func NewLRU(maxSize int) *LRU
```

**Parameters:**
- `maxSize` (int)

**Returns:**
- *LRU

## Methods

### Len

Len implements EvictionPolicy

```go
func (*LRU) Len() int
```

**Parameters:**
  None

**Returns:**
- int

### OnAccess

OnAccess implements EvictionPolicy

```go
func (*LRU) OnAccess(key any)
```

**Parameters:**
- `key` (any)

**Returns:**
  None

### OnAdd

OnAdd implements EvictionPolicy

```go
func (*LRU) OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time)
```

**Parameters:**
- `key` (any)
- `accessCount` (uint64)
- `createdAt` (time.Time)
- `accessedAt` (time.Time)

**Returns:**
  None

### OnRemove

OnRemove implements EvictionPolicy

```go
func (*LRU) OnRemove(key any)
```

**Parameters:**
- `key` (any)

**Returns:**
  None

### SelectVictim

SelectVictim implements EvictionPolicy

```go
func (*LRU) SelectVictim() (any, bool)
```

**Parameters:**
  None

**Returns:**
- any
- bool

## External Links

- [Package Overview](../packages/eviction.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/synapse/eviction)
- [Source Code](https://github.com/kolosys/synapse/tree/main/eviction)
