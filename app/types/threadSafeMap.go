package types

import (
	"fmt"
	"hash/fnv"
	"sync"
)

const (
	// shardCount is the number of shards in the ThreadSafeMap.
	// Using multiple shards reduces lock contention in concurrent access scenarios.
	shardCount = 32
)

// shard represents a single shard of the ThreadSafeMap that contains a portion
// of the key-value pairs and its own mutex for concurrent access control.
type shard[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

// ThreadSafeMap is a concurrent map implementation that uses sharding to
// reduce lock contention. It divides the map into multiple shards, each with
// its own read-write mutex, allowing for better performance in concurrent
// environments.
type ThreadSafeMap[K comparable, V any] struct {
	shards [shardCount]*shard[K, V]
}

// NewThreadSafeMap creates and initializes a new ThreadSafeMap with the specified key and value types.
// It initializes all shards with empty maps.
func NewThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K, V] {
	mp := &ThreadSafeMap[K, V]{}
	for i := range shardCount {
		mp.shards[i] = &shard[K, V]{data: make(map[K]V)}
	}
	return mp
}

// getShard determines which shard a key belongs to by hashing the key.
// This is an internal method used for distributing keys across shards.
func (mp *ThreadSafeMap[K, V]) getShard(key K) *shard[K, V] {
	h := fnv.New32()
	_, _ = h.Write(fmt.Append(nil, key))
	return mp.shards[h.Sum32()%shardCount]
}

// Set adds or updates a key-value pair in the map.
// It is safe to call concurrently with other methods.
func (mp *ThreadSafeMap[K, V]) Set(key K, value V) {
	shard := mp.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = value
}

// Get retrieves a value by key from the map.
// Returns the value and a boolean indicating whether the key was found.
// It is safe to call concurrently with other methods.
func (mp *ThreadSafeMap[K, V]) Get(key K) (V, bool) {
	shard := mp.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	value, ok := shard.data[key]

	return value, ok
}

// Delete removes a key-value pair from the map.
// It is safe to call concurrently with other methods.
func (mp *ThreadSafeMap[K, V]) Delete(key K) {
	shard := mp.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.data, key)
}

// Len returns the total number of key-value pairs in the map across all shards.
// It is safe to call concurrently with other methods.
func (mp *ThreadSafeMap[K, V]) Len() int {
	length := 0
	for _, shard := range mp.shards {
		shard.mu.RLock()
		length += len(shard.data)
		shard.mu.RUnlock()
	}
	return length
}

// Keys returns a slice containing all keys in the map.
func (mp *ThreadSafeMap[K, V]) Keys() []K {
	keys := make([]K, 0, mp.Len())

	for _, shard := range mp.shards {
		shard.mu.RLock()
		for k := range shard.data {
			keys = append(keys, k)
		}
		shard.mu.RUnlock()
	}

	return keys
}

// Clear removes all key-value pairs from the map.
func (mp *ThreadSafeMap[K, V]) Clear() {
	for _, shard := range mp.shards {
		shard.mu.Lock()
		shard.data = make(map[K]V)
		shard.mu.Unlock()
	}
}

// Contains checks if the map contains the specified key.
func (mp *ThreadSafeMap[K, V]) Contains(key K) bool {
	shard := mp.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	_, ok := shard.data[key]
	return ok
}

// GetOrSet returns the existing value for the key if present.
// Otherwise, it sets the provided value and returns it.
func (mp *ThreadSafeMap[K, V]) GetOrSet(key K, value V) V {
	shard := mp.getShard(key)

	shard.mu.RLock()
	existingValue, ok := shard.data[key]
	shard.mu.RUnlock()

	if ok {
		return existingValue
	}

	shard.mu.Lock()
	defer shard.mu.Unlock()

	// Double-check the key doesn't exist after acquiring the write lock
	if existingValue, ok := shard.data[key]; ok {
		return existingValue
	}

	shard.data[key] = value
	return value
}

// ForEach executes the provided function once for each key-value pair in the map.
// The function receives the key and value as parameters.
// Note: The map should not be modified during iteration.
func (mp *ThreadSafeMap[K, V]) ForEach(fn func(K, V)) {
	for _, shard := range mp.shards {
		shard.mu.RLock()
		for k, v := range shard.data {
			fn(k, v)
		}
		shard.mu.RUnlock()
	}
}
