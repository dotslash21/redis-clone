package tests

import (
	"strings"
	"testing"
)

// TestConfigCommands tests the CONFIG commands
func TestConfigCommands(t *testing.T) {
	// Setup test environment
	ts := NewTestSetup(t, 16383) // Different port from other tests
	defer ts.Close()

	// Test CONFIG SET and GET commands with RESP format
	t.Run("CONFIG SET and GET Commands (RESP Format)", func(t *testing.T) {
		key := "test-config-key"
		value := "test-config-value"

		// Test CONFIG SET
		setResponse, err := ts.Client.Execute("CONFIG", "SET", key, value)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG SET command: %v", err)
		}
		if setResponse != "OK" {
			t.Errorf("Expected 'OK', got %q", setResponse)
		}

		// Test CONFIG GET for the key we just set
		getResponse, err := ts.Client.Execute("CONFIG", "GET", key)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET command: %v", err)
		}

		// Check that the response contains both the key and value
		if !strings.Contains(getResponse, key) || !strings.Contains(getResponse, value) {
			t.Errorf("Expected response to contain key %q and value %q, got %q", key, value, getResponse)
		}

		// Test CONFIG GET with wildcard pattern
		// First set another related config
		relatedKey := "test-config-related"
		relatedValue := "related-value"
		_, err = ts.Client.Execute("CONFIG", "SET", relatedKey, relatedValue)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG SET command for related key: %v", err)
		}

		// Now get all test-config* keys
		patternResponse, err := ts.Client.Execute("CONFIG", "GET", "test-config*")
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET command with pattern: %v", err)
		}

		// Check that both keys and values are in the response
		if !strings.Contains(patternResponse, key) || !strings.Contains(patternResponse, value) ||
			!strings.Contains(patternResponse, relatedKey) || !strings.Contains(patternResponse, relatedValue) {
			t.Errorf("Expected pattern response to contain all keys and values, got %q", patternResponse)
		}

		// Test CONFIG GET for non-existent key
		emptyResponse, err := ts.Client.Execute("CONFIG", "GET", "nonexistent-config-key")
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET command for non-existent key: %v", err)
		}

		// For non-existent keys, should get an empty array
		if !strings.Contains(emptyResponse, "*0") {
			t.Errorf("Expected empty array response for non-existent key, got %q", emptyResponse)
		}

		// Test setting multiple key-value pairs with a single CONFIG SET command
		multi1Key := "multi-config-key1"
		multi1Value := "multi-config-value1"
		multi2Key := "multi-config-key2"
		multi2Value := "multi-config-value2"

		multiSetResponse, err := ts.Client.Execute("CONFIG", "SET",
			multi1Key, multi1Value,
			multi2Key, multi2Value)

		if err != nil {
			t.Fatalf("Failed to execute CONFIG SET command with multiple key-value pairs: %v", err)
		}
		if multiSetResponse != "OK" {
			t.Errorf("Expected 'OK' for multi-key CONFIG SET, got %q", multiSetResponse)
		}

		// Verify both keys were set correctly
		multi1GetResponse, err := ts.Client.Execute("CONFIG", "GET", multi1Key)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET for first multi key: %v", err)
		}
		if !strings.Contains(multi1GetResponse, multi1Key) || !strings.Contains(multi1GetResponse, multi1Value) {
			t.Errorf("Expected response to contain key %q and value %q, got %q", multi1Key, multi1Value, multi1GetResponse)
		}

		multi2GetResponse, err := ts.Client.Execute("CONFIG", "GET", multi2Key)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET for second multi key: %v", err)
		}
		if !strings.Contains(multi2GetResponse, multi2Key) || !strings.Contains(multi2GetResponse, multi2Value) {
			t.Errorf("Expected response to contain key %q and value %q, got %q", multi2Key, multi2Value, multi2GetResponse)
		}
	})

	// Test CONFIG SET and CONFIG GET commands with inline format
	t.Run("CONFIG SET and GET Commands (Inline Format)", func(t *testing.T) {
		key := "inline-config-key"
		value := "inline-config-value"

		// Test CONFIG SET with inline format
		setResponse, err := ts.Client.ExecuteInline("CONFIG SET", key, value)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG SET command with inline format: %v", err)
		}
		if setResponse != "OK" {
			t.Errorf("Expected 'OK', got %q", setResponse)
		}

		// Test CONFIG GET with inline format
		getResponse, err := ts.Client.ExecuteInline("CONFIG GET", key)
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET command with inline format: %v", err)
		}

		// Check that the response contains both the key and value
		if !strings.Contains(getResponse, key) || !strings.Contains(getResponse, value) {
			t.Errorf("Expected response to contain key %q and value %q, got %q", key, value, getResponse)
		}

		// Test CONFIG GET for non-existent key with inline format
		emptyResponse, err := ts.Client.ExecuteInline("CONFIG GET", "nonexistent-inline-key")
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET command for non-existent key with inline format: %v", err)
		}

		// For non-existent keys, should get an empty array
		if !strings.Contains(emptyResponse, "*0") {
			t.Errorf("Expected empty array response for non-existent key, got %q", emptyResponse)
		}

		// Test setting multiple key-value pairs with inline format
		multiKey1 := "inline-multi-key1"
		multiValue1 := "inline-multi-value1"
		multiKey2 := "inline-multi-key2"
		multiValue2 := "inline-multi-value2"

		multiSetResponse, err := ts.Client.ExecuteInline("CONFIG SET",
			multiKey1, multiValue1,
			multiKey2, multiValue2)

		if err != nil {
			t.Fatalf("Failed to execute CONFIG SET command with multiple key-value pairs in inline format: %v", err)
		}
		if multiSetResponse != "OK" {
			t.Errorf("Expected 'OK' for multi-key CONFIG SET with inline format, got %q", multiSetResponse)
		}

		// Verify both keys were set correctly
		multiGetResponse, err := ts.Client.ExecuteInline("CONFIG GET", "inline-multi-*")
		if err != nil {
			t.Fatalf("Failed to execute CONFIG GET for multi keys: %v", err)
		}
		if !strings.Contains(multiGetResponse, multiKey1) || !strings.Contains(multiGetResponse, multiValue1) ||
			!strings.Contains(multiGetResponse, multiKey2) || !strings.Contains(multiGetResponse, multiValue2) {
			t.Errorf("Expected response to contain both keys and values, got %q", multiGetResponse)
		}
	})

	// Test invalid CONFIG commands
	t.Run("Invalid CONFIG Commands", func(t *testing.T) {
		// Test CONFIG with no subcommand
		invalidResp, err := ts.Client.Execute("CONFIG")
		if err == nil {
			t.Errorf("Expected error for CONFIG with no subcommand, got: %q", invalidResp)
		} else if !strings.Contains(err.Error(), "ERR") {
			t.Errorf("Expected error to contain 'ERR', got: %v", err)
		}

		// Test CONFIG with invalid subcommand
		invalidResp, err = ts.Client.Execute("CONFIG", "INVALID")
		if err == nil {
			t.Errorf("Expected error for CONFIG with invalid subcommand, got: %q", invalidResp)
		} else if !strings.Contains(err.Error(), "ERR") {
			t.Errorf("Expected error to contain 'ERR', got: %v", err)
		}

		// Test CONFIG GET with missing key
		invalidResp, err = ts.Client.Execute("CONFIG", "GET")
		if err == nil {
			t.Errorf("Expected error for CONFIG GET with no key, got: %q", invalidResp)
		} else if !strings.Contains(err.Error(), "ERR") {
			t.Errorf("Expected error to contain 'ERR', got: %v", err)
		}

		// Test CONFIG SET with missing value
		invalidResp, err = ts.Client.Execute("CONFIG", "SET", "some-key")
		if err == nil {
			t.Errorf("Expected error for CONFIG SET with missing value, got: %q", invalidResp)
		} else if !strings.Contains(err.Error(), "ERR") {
			t.Errorf("Expected error to contain 'ERR', got: %v", err)
		}

		// Test CONFIG SET with odd number of arguments (missing a value)
		invalidResp, err = ts.Client.Execute("CONFIG", "SET", "key1", "value1", "key2")
		if err == nil {
			t.Errorf("Expected error for CONFIG SET with odd number of arguments, got: %q", invalidResp)
		} else if !strings.Contains(err.Error(), "ERR") {
			t.Errorf("Expected error to contain 'ERR', got: %v", err)
		}
	})
}
