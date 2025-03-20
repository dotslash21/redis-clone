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
	})
}
