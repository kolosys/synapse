# algorithms API

Complete API documentation for the algorithms package.

**Import Path:** `github.com/kolosys/synapse/algorithms`

## Package Documentation



## Functions

### DamerauLevenshtein
DamerauLevenshtein computes the Damerau-Levenshtein distance Similar to Levenshtein but also allows transpositions

```go
func DamerauLevenshtein(a, b string) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `a` | `string` | |
| `b` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of DamerauLevenshtein
result := DamerauLevenshtein(/* parameters */)
```

### Euclidean
Euclidean computes the Euclidean distance between two points Returns a similarity score (inverse of distance)

```go
func Euclidean(a, b []float64) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `a` | `[]float64` | |
| `b` | `[]float64` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of Euclidean
result := Euclidean(/* parameters */)
```

### Hamming
Hamming computes the Hamming distance between two strings Strings must be of equal length. Returns a normalized score between 0.0 and 1.0

```go
func Hamming(a, b string) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `a` | `string` | |
| `b` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of Hamming
result := Hamming(/* parameters */)
```

### HammingBytes
HammingBytes computes the Hamming distance between two byte slices

```go
func HammingBytes(a, b []byte) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `a` | `[]byte` | |
| `b` | `[]byte` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of HammingBytes
result := HammingBytes(/* parameters */)
```

### Levenshtein
Levenshtein computes the Levenshtein distance between two strings Returns a normalized score between 0.0 (completely different) and 1.0 (identical)

```go
func Levenshtein(a, b string) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `a` | `string` | |
| `b` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of Levenshtein
result := Levenshtein(/* parameters */)
```

### Manhattan
Manhattan computes the Manhattan distance between two points

```go
func Manhattan(a, b []float64) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `a` | `[]float64` | |
| `b` | `[]float64` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of Manhattan
result := Manhattan(/* parameters */)
```

## External Links

- [Package Overview](../packages/algorithms.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/synapse/algorithms)
- [Source Code](https://github.com/kolosys/synapse/tree/main/algorithms)
