package types

import (
	"sync"
	"testing"
)

func TestNewThreadSafeMap(t *testing.T) {
	m := NewThreadSafeMap[string, int]()
	if m == nil {
		t.Fatal("NewThreadSafeMap returned nil")
	}
	if len(m.shards) != shardCount {
		t.Errorf("Expected %d shards, got %d", shardCount, len(m.shards))
	}
}

func TestThreadSafeMap_SetAndGet(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	// Test setting and getting values
	m.Set("key1", 100)
	m.Set("key2", 200)

	val1, exists1 := m.Get("key1")
	if !exists1 || val1 != 100 {
		t.Errorf("Expected (100, true), got (%d, %t)", val1, exists1)
	}

	val2, exists2 := m.Get("key2")
	if !exists2 || val2 != 200 {
		t.Errorf("Expected (200, true), got (%d, %t)", val2, exists2)
	}

	// Test getting a non-existent key
	_, exists3 := m.Get("nonexistent")
	if exists3 {
		t.Error("Expected nonexistent key to return false")
	}

	// Test overwriting an existing key
	m.Set("key1", 150)
	val4, _ := m.Get("key1")
	if val4 != 150 {
		t.Errorf("Expected 150 after update, got %d", val4)
	}
}

func TestThreadSafeMap_Delete(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	m.Set("key1", 100)
	m.Set("key2", 200)

	// Delete an existing key
	m.Delete("key1")
	_, exists := m.Get("key1")
	if exists {
		t.Error("Key should not exist after deletion")
	}

	// Check that other keys are unaffected
	val, exists := m.Get("key2")
	if !exists || val != 200 {
		t.Errorf("Expected (200, true), got (%d, %t)", val, exists)
	}

	// Delete a non-existent key (should not panic)
	m.Delete("nonexistent")
}

func TestThreadSafeMap_Len(t *testing.T) {
	m := NewThreadSafeMap[string, int]()
	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}

	m.Set("key1", 100)
	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}

	m.Set("key2", 200)
	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}

	m.Delete("key1")
	if m.Len() != 1 {
		t.Errorf("Expected length 1 after deletion, got %d", m.Len())
	}
}

func TestThreadSafeMap_Keys(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	// Test with empty map
	keys := m.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(keys))
	}

	// Add keys and check
	m.Set("key1", 100)
	m.Set("key2", 200)

	keys = m.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	// Check if all keys are present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	if !keyMap["key1"] || !keyMap["key2"] {
		t.Error("Keys() did not return all expected keys")
	}
}

func TestThreadSafeMap_Clear(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	m.Set("key1", 100)
	m.Set("key2", 200)

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("Expected length 0 after Clear(), got %d", m.Len())
	}

	_, exists := m.Get("key1")
	if exists {
		t.Error("Expected key1 to be removed after Clear()")
	}
}

func TestThreadSafeMap_Contains(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	if m.Contains("key1") {
		t.Error("Expected Contains() to return false for nonexistent key")
	}

	m.Set("key1", 100)

	if !m.Contains("key1") {
		t.Error("Expected Contains() to return true for existing key")
	}
}

func TestThreadSafeMap_GetOrSet(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	// Test with a new key
	val := m.GetOrSet("key1", 100)
	if val != 100 {
		t.Errorf("Expected GetOrSet to return 100, got %d", val)
	}

	// The key should now exist in the map
	val, exists := m.Get("key1")
	if !exists || val != 100 {
		t.Errorf("Expected (100, true), got (%d, %t)", val, exists)
	}

	// Test with an existing key - should return existing value
	val = m.GetOrSet("key1", 200)
	if val != 100 {
		t.Errorf("Expected GetOrSet to return existing value 100, got %d", val)
	}

	// The value should not be updated
	val, _ = m.Get("key1")
	if val != 100 {
		t.Errorf("Expected value to remain 100, got %d", val)
	}
}

func TestThreadSafeMap_ForEach(t *testing.T) {
	m := NewThreadSafeMap[string, int]()

	m.Set("key1", 100)
	m.Set("key2", 200)

	// Count the number of key-value pairs visited
	count := 0
	sum := 0

	m.ForEach(func(k string, v int) {
		count++
		sum += v
	})

	if count != 2 {
		t.Errorf("Expected ForEach to visit 2 items, visited %d", count)
	}

	if sum != 300 {
		t.Errorf("Expected sum of values to be 300, got %d", sum)
	}
}

func TestThreadSafeMap_ConcurrentAccess(t *testing.T) {
	m := NewThreadSafeMap[int, int]()
	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // For both readers and writers

	// Test concurrent writing
	for i := 0; i < numGoroutines; i++ {
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := base*numOperations + j
				m.Set(key, key*10)
			}
		}(i)
	}

	// Test concurrent reading and other operations
	for i := 0; i < numGoroutines; i++ {
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := base*numOperations + j

				// Mix different operations
				switch j % 4 {
				case 0:
					m.Get(key)
				case 1:
					m.Contains(key)
				case 2:
					m.GetOrSet(key, key*10)
				case 3:
					if j%100 == 0 { // Delete occasionally
						m.Delete(key)
					}
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify the map has the expected number of elements
	// We can't know the exact count due to concurrent deletions
	count := m.Len()
	t.Logf("Final map size: %d", count)

	// The count should be close to numGoroutines * numOperations, minus the deletions
	if count == 0 {
		t.Error("Expected the map to contain elements after concurrent operations")
	}
}
