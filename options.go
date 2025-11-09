package synapse

import (
	"time"
)

// Options contains configuration options for the cache
type Options struct {
	NumShards           int
	MaxSize             int
	SimilarityThreshold float64
	EvictionPolicy      EvictionPolicy
	TTL                 time.Duration
	EnableStats         bool
}

// Option is a function that modifies Options
type Option func(*Options)

// defaultOptions returns default configuration
func defaultOptions() *Options {
	return &Options{
		NumShards:           16,
		MaxSize:             1000,
		SimilarityThreshold: 0.8,
		TTL:                 0, // No expiration by default
		EnableStats:         false,
	}
}

// WithShards sets the number of shards
func WithShards(n int) Option {
	return func(o *Options) {
		if n > 0 && n <= 256 {
			o.NumShards = n
		}
	}
}

// WithMaxSize sets the maximum cache size
func WithMaxSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.MaxSize = size
		}
	}
}

// WithThreshold sets the similarity threshold
func WithThreshold(t float64) Option {
	return func(o *Options) {
		if t >= 0.0 && t <= 1.0 {
			o.SimilarityThreshold = t
		}
	}
}

// WithEviction sets the eviction policy
func WithEviction(policy EvictionPolicy) Option {
	return func(o *Options) {
		o.EvictionPolicy = policy
	}
}

// WithTTL sets the time-to-live for cache entries
func WithTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.TTL = ttl
	}
}

// WithStats enables statistics tracking
func WithStats(enable bool) Option {
	return func(o *Options) {
		o.EnableStats = enable
	}
}

