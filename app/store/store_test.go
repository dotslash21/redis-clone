package store

import (
	"testing"
	"time"
)

func TestGetStore(t *testing.T) {
	// Test that GetStore returns a non-nil store
	s := GetStore()
	if s == nil {
		t.Fatal("GetStore returned nil")
	}

	// Test that GetStore returns the same instance each time (singleton)
	s2 := GetStore()
	if s != s2 {
		t.Error("GetStore did not return the same instance")
	}
}

func TestSetAndGet(t *testing.T) {
	s := GetStore()
	// Reset the store for testing
	s.data.Clear()

	// Test setting and getting values
	s.Set("key1", "value1", 0)
	s.Set("key2", "value2", 0)

	val1, err1 := s.Get("key1")
	if err1 != nil || val1 != "value1" {
		t.Errorf("Expected (value1, nil), got (%s, %v)", val1, err1)
	}

	val2, err2 := s.Get("key2")
	if err2 != nil || val2 != "value2" {
		t.Errorf("Expected (value2, nil), got (%s, %v)", val2, err2)
	}

	// Test getting a non-existent key
	_, err3 := s.Get("nonexistent")
	if err3 == nil {
		t.Error("Expected error for nonexistent key, got nil")
	}

	// Test overwriting an existing key
	s.Set("key1", "newvalue1", 0)
	val4, _ := s.Get("key1")
	if val4 != "newvalue1" {
		t.Errorf("Expected newvalue1 after update, got %s", val4)
	}
}

func TestExpiry(t *testing.T) {
	s := GetStore()
	s.data.Clear()

	// Set a key with a short TTL
	s.Set("expiring", "value", 50*time.Millisecond)

	// Key should exist before expiry
	val, err := s.Get("expiring")
	if err != nil || val != "value" {
		t.Errorf("Expected value to exist before expiry: got (%s, %v)", val, err)
	}

	// Wait for the key to expire
	time.Sleep(100 * time.Millisecond)

	// Key should be expired now
	_, err = s.Get("expiring")
	if err == nil {
		t.Error("Expected error for expired key, got nil")
	}

	// Test FlushExpired functionality
	s.Set("expiring1", "value1", 50*time.Millisecond)
	s.Set("expiring2", "value2", 50*time.Millisecond)
	s.Set("nonexpiring", "value3", 0)

	// Wait for keys to expire
	time.Sleep(100 * time.Millisecond)

	// Flush expired keys
	s.FlushExpired()

	// Check that expired keys are gone and non-expiring keys remain
	_, err1 := s.Get("expiring1")
	if err1 == nil {
		t.Error("Expected expiring1 to be removed")
	}

	_, err2 := s.Get("expiring2")
	if err2 == nil {
		t.Error("Expected expiring2 to be removed")
	}

	val3, err3 := s.Get("nonexpiring")
	if err3 != nil || val3 != "value3" {
		t.Errorf("Expected nonexpiring key to remain, got (%s, %v)", val3, err3)
	}
}

func TestIsExpired(t *testing.T) {
	s := GetStore()
	s.data.Clear()

	// Set up test data
	s.Set("nonexpiring", "value1", 0)
	s.Set("expiring", "value2", 50*time.Millisecond)

	// Test non-expiring key
	if s.isExpired("nonexpiring") {
		t.Error("Non-expiring key incorrectly reported as expired")
	}

	// Test non-existent key
	if !s.isExpired("nonexistent") {
		t.Error("Non-existent key should be reported as expired")
	}

	// Test expiring key before expiry
	if s.isExpired("expiring") {
		t.Error("Key should not be expired yet")
	}

	// Wait for expiry
	time.Sleep(100 * time.Millisecond)

	// Test expiring key after expiry
	if !s.isExpired("expiring") {
		t.Error("Key should be expired now")
	}

	// Verify the expired key was deleted
	_, exists := s.data.Get("expiring")
	if exists {
		t.Error("Expired key should be deleted from the store")
	}
}

func TestExpiryHeap(t *testing.T) {
	s := GetStore()
	s.data.Clear()

	// Clear expiry heap for testing
	s.exp.mu.Lock()
	s.exp.items = []expiryItem{}
	s.exp.mu.Unlock()

	// Set keys with different expiry times
	s.Set("key1", "value1", 300*time.Millisecond)
	s.Set("key2", "value2", 100*time.Millisecond)
	s.Set("key3", "value3", 200*time.Millisecond)

	// The heap should be ordered by expiry time (earliest first)
	s.exp.mu.RLock()
	if len(s.exp.items) != 3 {
		t.Errorf("Expected 3 items in expiry heap, got %d", len(s.exp.items))
	}

	// Check that the earliest expiring key is first in the heap
	if len(s.exp.items) > 0 && s.exp.items[0].key != "key2" {
		t.Errorf("Expected key2 to be at the top of the heap, got %s", s.exp.items[0].key)
	}
	s.exp.mu.RUnlock()

	// Wait for the first key to expire
	time.Sleep(150 * time.Millisecond)

	// Trigger expiry check
	_, _ = s.Get("key2")

	// Wait for second key to expire
	time.Sleep(100 * time.Millisecond)

	// Flush expired keys
	s.FlushExpired()

	// Only key1 should remain
	_, err2 := s.Get("key2")
	if err2 == nil {
		t.Error("key2 should be expired")
	}

	_, err3 := s.Get("key3")
	if err3 == nil {
		t.Error("key3 should be expired")
	}

	val1, err1 := s.Get("key1")
	if err1 != nil {
		t.Error("key1 should still exist")
	}
	if val1 != "value1" {
		t.Errorf("Expected value1, got %s", val1)
	}

	// Wait for the last key to expire
	time.Sleep(100 * time.Millisecond)

	// All keys should be expired now
	s.FlushExpired()
	_, err1 = s.Get("key1")
	if err1 == nil {
		t.Error("key1 should be expired")
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := GetStore()
	s.data.Clear()

	const numGoroutines = 10
	const numOperations = 100

	// Channel to coordinate goroutines
	done := make(chan bool, numGoroutines*2)

	// Test concurrent writing
	for i := 0; i < numGoroutines; i++ {
		go func(base int) {
			for j := 0; j < numOperations; j++ {
				key := "key" + string(rune('A'+base)) + string(rune('0'+j%10))
				value := "value" + string(rune('A'+base)) + string(rune('0'+j%10))
				ttl := time.Duration(j%3) * 500 * time.Millisecond // Mix of TTLs
				s.Set(key, value, ttl)
			}
			done <- true
		}(i)
	}

	// Test concurrent reading
	for i := 0; i < numGoroutines; i++ {
		go func(base int) {
			for j := 0; j < numOperations; j++ {
				key := "key" + string(rune('A'+base)) + string(rune('0'+j%10))
				s.Get(key) // Just attempt to get, don't check result as it might expire
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < numGoroutines*2; i++ {
		<-done
	}

	// Basic check that the store is still functioning
	s.Set("final", "test", 0)
	val, err := s.Get("final")
	if err != nil || val != "test" {
		t.Errorf("Store not functioning after concurrent access: (%s, %v)", val, err)
	}

	// Wait for keys with TTL to expire
	time.Sleep(1500 * time.Millisecond)
	s.FlushExpired()

	// Some keys should remain (those set with TTL 0)
	if s.data.Len() == 0 {
		t.Error("Expected some keys to remain after expiry")
	}
}
