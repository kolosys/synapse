package synapse

// SimilarityFunc is a function type that computes similarity between two keys
// It should return a score between 0.0 (completely different) and 1.0 (identical)
type SimilarityFunc[K comparable] func(a, b K) float64

// Similarity is an interface for similarity computation
type Similarity[K comparable] interface {
	// Score computes the similarity score between two keys
	// Returns a value between 0.0 and 1.0
	Score(a, b K) float64

	// Threshold returns the minimum similarity score for a match
	Threshold() float64
}

// similarityAdapter adapts a SimilarityFunc to the Similarity interface
type similarityAdapter[K comparable] struct {
	fn        SimilarityFunc[K]
	threshold float64
}

// Score implements the Similarity interface
func (s *similarityAdapter[K]) Score(a, b K) float64 {
	return s.fn(a, b)
}

// Threshold implements the Similarity interface
func (s *similarityAdapter[K]) Threshold() float64 {
	return s.threshold
}

// NewSimilarity creates a Similarity from a SimilarityFunc
func NewSimilarity[K comparable](fn SimilarityFunc[K], threshold float64) Similarity[K] {
	return &similarityAdapter[K]{
		fn:        fn,
		threshold: threshold,
	}
}
