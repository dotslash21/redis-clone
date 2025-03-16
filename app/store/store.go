package store

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

// RedisValue holds both the value and metadata (type info, expiry).
type RedisValue struct {
	Value    string
	ExpireAt time.Time
}

// expiryItem holds the key and its expiry time for the min-heap.
type expiryItem struct {
	key      string
	expireAt time.Time
}

// expiryHeap is a min-heap based on expiry time.
type expiryHeap []expiryItem

func (eh expiryHeap) Len() int           { return len(eh) }
func (eh expiryHeap) Less(i, j int) bool { return eh[i].expireAt.Before(eh[j].expireAt) }
func (eh expiryHeap) Swap(i, j int)      { eh[i], eh[j] = eh[j], eh[i] }

func (eh *expiryHeap) Push(x interface{}) {
	*eh = append(*eh, x.(expiryItem))
}

func (eh *expiryHeap) Pop() interface{} {
	old := *eh
	n := len(old)
	item := old[n-1]
	*eh = old[0 : n-1]
	return item
}

// Store is a Redis-like key-value store with a min-heap for expiry.
type Store struct {
	data map[string]*RedisValue
	mu   sync.RWMutex
	exp  expiryHeap
}

// store is a singleton instance of Store
var store *Store

// GetStore returns the store instance, creating it if it doesn't exist.
func GetStore() *Store {
	if store == nil {
		store = &Store{
			data: make(map[string]*RedisValue),
			mu:   sync.RWMutex{},
			exp:  expiryHeap{},
		}
		heap.Init(&store.exp)
	}
	return store
}

// flushExpired deletes all expired keys at once.
func (s *Store) FlushExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for len(s.exp) > 0 && !s.exp[0].expireAt.IsZero() && s.exp[0].expireAt.Before(now) {
		top := heap.Pop(&s.exp).(expiryItem)
		if v, ok := s.data[top.key]; ok && !v.ExpireAt.IsZero() && v.ExpireAt.Before(now) {
			delete(s.data, top.key)
		}
	}
}

// isExpired returns true if the key has expired and removes it from the store.
func (s *Store) isExpired(key string) bool {
	val, exists := s.data[key]
	if !exists {
		return true
	}
	if !val.ExpireAt.IsZero() && time.Now().After(val.ExpireAt) {
		delete(s.data, key)
		return true
	}
	return false
}

// Set stores a string value with optional expiry in the store.
func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiry := time.Time{}
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
		heap.Push(&s.exp, expiryItem{key: key, expireAt: expiry})
	}

	s.data[key] = &RedisValue{
		Value:    value,
		ExpireAt: expiry,
	}
}

// Get retrieves a string value from the store, returning ok = false if missing or expired.
func (s *Store) Get(key string) (value string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.isExpired(key) {
		delete(s.data, key)
		return "", errors.New("key expired")
	}

	val, exists := s.data[key]
	if !exists || s.isExpired(key) {
		return "", errors.New("key not found")
	}

	return val.Value, nil
}
