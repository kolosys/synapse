# main

This example demonstrates basic usage of the library.

## Source Code

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
	for i := range minLen {
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

## Running the Example

To run this example:

```bash
cd string_similarity
go run main.go
```

## Expected Output

```
Hello from Proton examples!
```
