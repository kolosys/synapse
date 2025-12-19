package synapse

import (
	"sync/atomic"
)

// Stats contains cache performance statistics
type Stats struct {
	Hits            uint64
	Misses          uint64
	Sets            uint64
	Deletes         uint64
	SimilarSearches uint64
	SimilarHits     uint64
	Evictions       uint64
	Expired         uint64
}

// shardStats contains per-shard statistics using atomic counters
type shardStats struct {
	hits            atomic.Uint64
	misses          atomic.Uint64
	sets            atomic.Uint64
	deletes         atomic.Uint64
	similarSearches atomic.Uint64
	similarHits     atomic.Uint64
	evictions       atomic.Uint64
	expired         atomic.Uint64
}

// newShardStats creates a new shard stats tracker
func newShardStats() *shardStats {
	return &shardStats{}
}

// recordHit increments the hit counter
func (s *shardStats) recordHit() {
	s.hits.Add(1)
}

// recordMiss increments the miss counter
func (s *shardStats) recordMiss() {
	s.misses.Add(1)
}

// recordSet increments the set counter
func (s *shardStats) recordSet() {
	s.sets.Add(1)
}

// recordDelete increments the delete counter
func (s *shardStats) recordDelete() {
	s.deletes.Add(1)
}

// recordSimilarSearch increments the similarity search counter
func (s *shardStats) recordSimilarSearch() {
	s.similarSearches.Add(1)
}

// recordSimilarHit increments the similarity hit counter
func (s *shardStats) recordSimilarHit() {
	s.similarHits.Add(1)
}

// recordEviction increments the eviction counter
func (s *shardStats) recordEviction() {
	s.evictions.Add(1)
}

// recordExpired increments the expired counter
func (s *shardStats) recordExpired() {
	s.expired.Add(1)
}

// snapshot returns a snapshot of current statistics
func (s *shardStats) snapshot() Stats {
	return Stats{
		Hits:            s.hits.Load(),
		Misses:          s.misses.Load(),
		Sets:            s.sets.Load(),
		Deletes:         s.deletes.Load(),
		SimilarSearches: s.similarSearches.Load(),
		SimilarHits:     s.similarHits.Load(),
		Evictions:       s.evictions.Load(),
		Expired:         s.expired.Load(),
	}
}
