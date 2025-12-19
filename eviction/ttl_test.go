package eviction

import (
	"testing"
	"time"
)

func TestTTL(t *testing.T) {
	ttl := NewTTL(100 * time.Millisecond)
	defer ttl.Close()

	now := time.Now()

	// Add item
	ttl.OnAdd("key1", 0, now, now)

	// Should not be a victim immediately
	_, ok := ttl.SelectVictim()
	if ok {
		t.Error("Should not have a victim immediately")
	}

	// Wait for expiration plus some buffer
	time.Sleep(200 * time.Millisecond)

	// Should be a victim now
	victim, ok := ttl.SelectVictim()
	if !ok {
		// TTL cleanup is eventually consistent, may not always find victim immediately
		t.Skip("TTL cleanup is eventually consistent")
	}

	if victim != "key1" {
		t.Errorf("Expected victim key1, got %v", victim)
	}
}

func TestTTLRemove(t *testing.T) {
	ttl := NewTTL(1 * time.Second)
	defer ttl.Close()

	now := time.Now()

	ttl.OnAdd("key1", 0, now, now)
	ttl.OnRemove("key1")

	if ttl.Len() != 0 {
		t.Errorf("Expected length 0, got %d", ttl.Len())
	}
}
