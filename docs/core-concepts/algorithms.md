# Algorithms

The `algorithms` package provides built-in similarity functions for computing distance and similarity scores between values.

**Import Path:** `github.com/kolosys/synapse/algorithms`

## Overview

All algorithms in this package return a normalized similarity score between `0.0` (completely different) and `1.0` (identical). This makes them directly compatible with Synapse's `SimilarityFunc[K]` type for string keys.

## String Distance Algorithms

### Levenshtein

Computes the [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance) (edit distance) between two strings. This measures the minimum number of single-character edits (insertions, deletions, or substitutions) required to transform one string into another.

```go
score := algorithms.Levenshtein("hello", "helo")
// score = 0.80 (1 deletion needed)

score = algorithms.Levenshtein("kitten", "sitting")
// score = 0.57 (3 edits needed)

score = algorithms.Levenshtein("hello", "hello")
// score = 1.0 (identical)
```

**Use cases:**

- Typo correction
- Fuzzy string matching
- Spell checking

**Complexity:** O(n × m) time and space, where n and m are string lengths.

### Damerau-Levenshtein

Similar to Levenshtein but also counts transpositions (swapping adjacent characters) as a single edit. This better handles common typing errors where characters are accidentally swapped.

```go
score := algorithms.DamerauLevenshtein("hello", "hlelo")
// score = 0.80 (1 transposition)

// Compare with standard Levenshtein
score = algorithms.Levenshtein("hello", "hlelo")
// score = 0.60 (2 edits: delete + insert)
```

**Use cases:**

- Typo detection where transpositions are common
- Keyboard input correction
- DNA sequence comparison

**Complexity:** O(n × m) time and space.

### Hamming

Computes the [Hamming distance](https://en.wikipedia.org/wiki/Hamming_distance) between two strings of **equal length**. This counts the number of positions where corresponding characters differ.

```go
score := algorithms.Hamming("hello", "hallo")
// score = 0.80 (1 position differs)

score = algorithms.Hamming("hello", "world")
// score = 0.20 (4 positions differ)

// Returns 0.0 for different lengths
score = algorithms.Hamming("hello", "hi")
// score = 0.0
```

**Use cases:**

- Fixed-format identifiers
- Error detection in data transmission
- Comparing binary representations

**Complexity:** O(n) time.

### HammingBytes

Computes Hamming distance between byte slices at the bit level. Useful for comparing binary data or hash values.

```go
a := []byte{0xFF, 0x00}  // 11111111 00000000
b := []byte{0xFE, 0x01}  // 11111110 00000001
score := algorithms.HammingBytes(a, b)
// score = 0.875 (2 bits differ out of 16)
```

**Use cases:**

- Perceptual hashing (image similarity)
- SimHash comparison
- Binary fingerprint matching

**Complexity:** O(n) time.

## Vector Distance Algorithms

### Euclidean

Computes the [Euclidean distance](https://en.wikipedia.org/wiki/Euclidean_distance) between two points in n-dimensional space. Returns a similarity score inversely proportional to the distance.

```go
a := []float64{0, 0}
b := []float64{3, 4}
score := algorithms.Euclidean(a, b)
// score = 0.167 (distance = 5)

// Identical points
score = algorithms.Euclidean([]float64{1, 2, 3}, []float64{1, 2, 3})
// score = 1.0
```

The similarity is calculated as `1.0 / (1.0 + distance)`.

**Use cases:**

- Feature vector comparison
- Embedding similarity
- K-nearest neighbors

**Complexity:** O(n) time.

### Manhattan

Computes the [Manhattan distance](https://en.wikipedia.org/wiki/Taxicab_geometry) (L1 norm) between two points. This is the sum of absolute differences along each dimension.

```go
a := []float64{0, 0}
b := []float64{3, 4}
score := algorithms.Manhattan(a, b)
// score = 0.125 (distance = 7)

score = algorithms.Manhattan([]float64{1, 2}, []float64{4, 6})
// score = 0.125 (distance = 7)
```

**Use cases:**

- Grid-based pathfinding
- Feature comparison where dimensions are independent
- Sparse vector comparison

**Complexity:** O(n) time.

## Using with Synapse

All string algorithms can be used directly as similarity functions:

```go
cache := synapse.New[string, string](
    synapse.WithThreshold(0.8),
)

// Use Levenshtein directly
cache.WithSimilarity(algorithms.Levenshtein)

// Or use Damerau-Levenshtein
cache.WithSimilarity(algorithms.DamerauLevenshtein)
```

For vector algorithms, wrap them in a function:

```go
type Vector []float64

cache := synapse.New[Vector, string](
    synapse.WithThreshold(0.5),
)

cache.WithSimilarity(func(a, b Vector) float64 {
    return algorithms.Euclidean(a, b)
})
```

## Custom Algorithms

You can implement custom similarity functions. They must:

1. Accept two keys of the same type
2. Return a `float64` between `0.0` and `1.0`
3. Return `1.0` for identical keys

```go
// Custom case-insensitive prefix similarity
func prefixSimilarity(a, b string) float64 {
    a, b = strings.ToLower(a), strings.ToLower(b)
    if a == b {
        return 1.0
    }

    shorter, longer := a, b
    if len(b) < len(a) {
        shorter, longer = b, a
    }

    if strings.HasPrefix(longer, shorter) {
        return float64(len(shorter)) / float64(len(longer))
    }
    return 0.0
}

cache.WithSimilarity(prefixSimilarity)
```

## Algorithm Selection Guide

| Algorithm           | Best For                       | Requirements         |
| ------------------- | ------------------------------ | -------------------- |
| Levenshtein         | General string matching, typos | None                 |
| Damerau-Levenshtein | Keyboard typos, transpositions | None                 |
| Hamming             | Fixed-length identifiers       | Equal length strings |
| HammingBytes        | Binary data, hashes            | Equal length slices  |
| Euclidean           | Embeddings, dense vectors      | Same dimensions      |
| Manhattan           | Sparse vectors, grid distances | Same dimensions      |

## Performance Considerations

- **Levenshtein/Damerau-Levenshtein**: O(n×m) complexity can be slow for very long strings. Consider truncating or pre-filtering for large datasets.
- **Hamming**: Very fast O(n), but requires equal lengths.
- **Vector algorithms**: Fast O(n), scale linearly with dimensions.

For high-throughput similarity searches, consider:

- Using shorter keys when possible
- Pre-computing hashes for approximate matching
- Limiting the number of entries per shard

## Further Reading

- [API Reference](../api-reference/algorithms.md) - Complete function signatures
- [Performance Tuning](../advanced/performance-tuning.md) - Optimization tips
- [Quick Start](../getting-started/quick-start.md) - Usage examples
