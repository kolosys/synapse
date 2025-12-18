package algorithms

import (
	"math"
	"testing"
)

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		a, b     string
		expected float64
	}{
		{"hello", "hello", 1.0},
		{"hello", "hallo", 0.8},
		{"", "", 1.0},
		{"abc", "xyz", 0.0},
		{"kitten", "sitting", 0.5714285714285714},
	}
	
	for _, tt := range tests {
		result := Levenshtein(tt.a, tt.b)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Levenshtein(%q, %q) = %f; want %f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestHamming(t *testing.T) {
	tests := []struct {
		a, b     string
		expected float64
	}{
		{"hello", "hello", 1.0},
		{"hello", "hallo", 0.8},
		{"abc", "abc", 1.0},
		{"abc", "xyz", 0.0},
		{"", "", 1.0},
	}
	
	for _, tt := range tests {
		result := Hamming(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("Hamming(%q, %q) = %f; want %f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestHammingBytes(t *testing.T) {
	tests := []struct {
		a, b     []byte
		expected float64
	}{
		{[]byte{0xFF}, []byte{0xFF}, 1.0},
		{[]byte{0xFF}, []byte{0x00}, 0.0},
		{[]byte{0xF0}, []byte{0x0F}, 0.0},
	}
	
	for _, tt := range tests {
		result := HammingBytes(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("HammingBytes(%v, %v) = %f; want %f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestEuclidean(t *testing.T) {
	tests := []struct {
		a, b     []float64
		minScore float64
	}{
		{[]float64{1, 2, 3}, []float64{1, 2, 3}, 1.0},
		{[]float64{0, 0}, []float64{3, 4}, 0.1},
		{[]float64{}, []float64{}, 1.0},
	}
	
	for _, tt := range tests {
		result := Euclidean(tt.a, tt.b)
		if result < tt.minScore {
			t.Errorf("Euclidean(%v, %v) = %f; want >= %f", tt.a, tt.b, result, tt.minScore)
		}
	}
}

func TestManhattan(t *testing.T) {
	tests := []struct {
		a, b     []float64
		minScore float64
	}{
		{[]float64{1, 2, 3}, []float64{1, 2, 3}, 1.0},
		{[]float64{0, 0}, []float64{5, 5}, 0.05},
	}
	
	for _, tt := range tests {
		result := Manhattan(tt.a, tt.b)
		if result < tt.minScore {
			t.Errorf("Manhattan(%v, %v) = %f; want >= %f", tt.a, tt.b, result, tt.minScore)
		}
	}
}

func TestDamerauLevenshtein(t *testing.T) {
	tests := []struct {
		a, b     string
		minScore float64
	}{
		{"hello", "hello", 1.0},
		{"abc", "abc", 1.0},
		{"ab", "ba", 0.5},
	}
	
	for _, tt := range tests {
		result := DamerauLevenshtein(tt.a, tt.b)
		if result < tt.minScore {
			t.Errorf("DamerauLevenshtein(%q, %q) = %f; want >= %f", tt.a, tt.b, result, tt.minScore)
		}
	}
}

