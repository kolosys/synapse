package algorithms

import (
	"math"
)

// Levenshtein computes the Levenshtein distance between two strings
// Returns a normalized score between 0.0 (completely different) and 1.0 (identical)
func Levenshtein(a, b string) float64 {
	if a == b {
		return 1.0
	}

	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}

	for i := 0; i <= len(a); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	distance := matrix[len(a)][len(b)]
	maxLen := max(len(a), len(b))

	return 1.0 - float64(distance)/float64(maxLen)
}

// Hamming computes the Hamming distance between two strings
// Strings must be of equal length. Returns a normalized score between 0.0 and 1.0
func Hamming(a, b string) float64 {
	if a == b {
		return 1.0
	}

	if len(a) != len(b) {
		return 0.0
	}

	if len(a) == 0 {
		return 1.0
	}

	differences := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			differences++
		}
	}

	return 1.0 - float64(differences)/float64(len(a))
}

// HammingBytes computes the Hamming distance between two byte slices
func HammingBytes(a, b []byte) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	if len(a) == 0 {
		return 1.0
	}

	differences := 0
	for i := range a {
		xor := a[i] ^ b[i]
		differences += popcount(xor)
	}

	totalBits := len(a) * 8
	return 1.0 - float64(differences)/float64(totalBits)
}

// popcount counts the number of set bits
func popcount(x byte) int {
	count := 0
	for x != 0 {
		count++
		x &= x - 1
	}
	return count
}

// DamerauLevenshtein computes the Damerau-Levenshtein distance
// Similar to Levenshtein but also allows transpositions
func DamerauLevenshtein(a, b string) float64 {
	if a == b {
		return 1.0
	}

	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	lenA, lenB := len(a), len(b)
	maxDist := lenA + lenB

	h := make(map[rune]int)
	matrix := make([][]int, lenA+2)
	for i := range matrix {
		matrix[i] = make([]int, lenB+2)
	}

	matrix[0][0] = maxDist
	for i := 0; i <= lenA; i++ {
		matrix[i+1][0] = maxDist
		matrix[i+1][1] = i
	}
	for j := 0; j <= lenB; j++ {
		matrix[0][j+1] = maxDist
		matrix[1][j+1] = j
	}

	for i := 1; i <= lenA; i++ {
		db := 0
		for j := 1; j <= lenB; j++ {
			k := h[rune(b[j-1])]
			l := db
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
				db = j
			}

			matrix[i+1][j+1] = min(
				matrix[i][j]+cost,              // substitution
				matrix[i+1][j]+1,               // insertion
				matrix[i][j+1]+1,               // deletion
				matrix[k][l]+(i-k-1)+1+(j-l-1), // transposition
			)
		}
		h[rune(a[i-1])] = i
	}

	distance := matrix[lenA+1][lenB+1]
	maxLen := max(lenA, lenB)

	return 1.0 - float64(distance)/float64(maxLen)
}

func min(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// Euclidean computes the Euclidean distance between two points
// Returns a similarity score (inverse of distance)
func Euclidean(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	if len(a) == 0 {
		return 1.0
	}

	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	distance := math.Sqrt(sum)

	return 1.0 / (1.0 + distance)
}

// Manhattan computes the Manhattan distance between two points
func Manhattan(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	if len(a) == 0 {
		return 1.0
	}

	sum := 0.0
	for i := range a {
		sum += math.Abs(a[i] - b[i])
	}

	return 1.0 / (1.0 + sum)
}
