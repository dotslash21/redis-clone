package store

import (
	"container/heap"
	"errors"
	"sync"
	"time"

	"github.com/dotslash21/redis-clone/app/types"
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
type expiryHeap struct {
	items []expiryItem
	mu    sync.RWMutex
}

func (eh *expiryHeap) Len() int           { return len(eh.items) }
func (eh *expiryHeap) Less(i, j int) bool { return eh.items[i].expireAt.Before(eh.items[j].expireAt) }
func (eh *expiryHeap) Swap(i, j int)      { eh.items[i], eh.items[j] = eh.items[j], eh.items[i] }
func (eh *expiryHeap) Push(x interface{}) {
	eh.items = append(eh.items, x.(expiryItem))
}

func (eh *expiryHeap) Pop() interface{} {
	old := eh.items
	n := len(old)
	item := old[n-1]
	eh.items = old[0 : n-1]
	return item
}

// Store is a Redis-like key-value store with a min-heap for expiry.
type Store struct {
	data *types.ThreadSafeMap[string, *RedisValue]
	exp  expiryHeap
}

// store is a singleton instance of Store
var store *Store

// GetStore returns the store instance, creating it if it doesn't exist.
func GetStore() *Store {
	if store == nil {
		store = &Store{
			data: types.NewThreadSafeMap[string, *RedisValue](),
			exp:  expiryHeap{items: []expiryItem{}},
		}
		heap.Init(&store.exp)
	}
	return store
}

// flushExpired deletes all expired keys at once.
func (s *Store) FlushExpired() {
	s.exp.mu.Lock()
	defer s.exp.mu.Unlock()

	now := time.Now()
	for len(s.exp.items) > 0 && !s.exp.items[0].expireAt.IsZero() && s.exp.items[0].expireAt.Before(now) {
		top := heap.Pop(&s.exp).(expiryItem)
		if v, ok := s.data.Get(top.key); ok && !v.ExpireAt.IsZero() && v.ExpireAt.Before(now) {
			s.data.Delete(top.key)
		}
	}
}

// isExpired returns true if the key has expired and removes it from the store.
func (s *Store) isExpired(key string) bool {
	val, exists := s.data.Get(key)
	if !exists {
		return true
	}
	if !val.ExpireAt.IsZero() && time.Now().After(val.ExpireAt) {
		s.data.Delete(key)
		return true
	}
	return false
}

// Set stores a string value with optional expiry in the store.
func (s *Store) Set(key, value string, ttl time.Duration) {
	expiry := time.Time{}
	if ttl > 0 {
		expiry = time.Now().Add(ttl)

		s.exp.mu.Lock()
		heap.Push(&s.exp, expiryItem{key: key, expireAt: expiry})
		s.exp.mu.Unlock()
	}

	s.data.Set(key, &RedisValue{
		Value:    value,
		ExpireAt: expiry,
	})
}

// Get retrieves a string value from the store, returning ok = false if missing or expired.
func (s *Store) Get(key string) (value string, err error) {
	if s.isExpired(key) {
		s.data.Delete(key)
		return "", errors.New("key expired")
	}

	val, exists := s.data.Get(key)
	if !exists || s.isExpired(key) {
		return "", errors.New("key not found")
	}

	return val.Value, nil
}
