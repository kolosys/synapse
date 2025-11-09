package eviction

import (
	"container/list"
	"sync"
	"time"
)

// LRU implements a Least Recently Used eviction policy
type LRU struct {
	mu       sync.RWMutex
	list     *list.List
	items    map[any]*list.Element
	maxSize  int
}

type lruEntry struct {
	key any
}

// NewLRU creates a new LRU eviction policy
func NewLRU(maxSize int) *LRU {
	return &LRU{
		list:    list.New(),
		items:   make(map[any]*list.Element),
		maxSize: maxSize,
	}
}

// OnAccess implements EvictionPolicy
func (l *LRU) OnAccess(key any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if elem, ok := l.items[key]; ok {
		l.list.MoveToFront(elem)
	}
}

// OnAdd implements EvictionPolicy
func (l *LRU) OnAdd(key any, accessCount uint64, createdAt, accessedAt time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if elem, ok := l.items[key]; ok {
		l.list.MoveToFront(elem)
		return
	}
	
	entry := &lruEntry{key: key}
	elem := l.list.PushFront(entry)
	l.items[key] = elem
}

// OnRemove implements EvictionPolicy
func (l *LRU) OnRemove(key any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if elem, ok := l.items[key]; ok {
		l.list.Remove(elem)
		delete(l.items, key)
	}
}

// SelectVictim implements EvictionPolicy
func (l *LRU) SelectVictim() (any, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	if l.list.Len() == 0 {
		return nil, false
	}
	
	elem := l.list.Back()
	if elem == nil {
		return nil, false
	}
	
	entry := elem.Value.(*lruEntry)
	return entry.key, true
}

// Len implements EvictionPolicy
func (l *LRU) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.list.Len()
}

