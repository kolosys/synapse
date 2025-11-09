package eviction

import (
	"time"
)

// EvictionPolicy defines the interface for cache eviction strategies
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

// CombinedPolicy combines multiple eviction policies with weighted scoring
type CombinedPolicy struct {
	policies []EvictionPolicy
	weights  []float64
}

// NewCombinedPolicy creates a new combined eviction policy
func NewCombinedPolicy(policies []EvictionPolicy, weights []float64) *CombinedPolicy {
	if len(policies) != len(weights) {
		panic("policies and weights must have the same length")
	}
	
	// Normalize weights
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	normalized := make([]float64, len(weights))
	for i, w := range weights {
		normalized[i] = w / sum
	}
	
	return &CombinedPolicy{
		policies: policies,
		weights:  normalized,
	}
}

// OnAccess implements EvictionPolicy
func (c *CombinedPolicy) OnAccess(key any) {
	for _, policy := range c.policies {
		policy.OnAccess(key)
	}
}

// OnAdd implements EvictionPolicy
func (c *CombinedPolicy) OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time) {
	for _, policy := range c.policies {
		policy.OnAdd(key, accessCount, createdAt, accessedAt)
	}
}

// OnRemove implements EvictionPolicy
func (c *CombinedPolicy) OnRemove(key any) {
	for _, policy := range c.policies {
		policy.OnRemove(key)
	}
}

// SelectVictim implements EvictionPolicy
// It uses the first policy's victim selection
func (c *CombinedPolicy) SelectVictim() (any, bool) {
	if len(c.policies) == 0 {
		return nil, false
	}
	return c.policies[0].SelectVictim()
}

// Len implements EvictionPolicy
func (c *CombinedPolicy) Len() int {
	if len(c.policies) == 0 {
		return 0
	}
	return c.policies[0].Len()
}

