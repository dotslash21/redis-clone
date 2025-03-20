package tests

import (
	"testing"
	"time"
)

// TestSetAndGetCommands tests the SET and GET commands
func TestSetAndGetCommands(t *testing.T) {
	// Setup test environment
	ts := NewTestSetup(t, 16382) // Different port from other tests
	defer ts.Close()

	// Test SET and GET commands with RESP format
	t.Run("SET and GET Commands (RESP Format)", func(t *testing.T) {
		key := "mykey-resp"
		value := "myvalue-resp"

		// Test SET
		setResponse, err := ts.Client.Execute("SET", key, value)
		if err != nil {
			t.Fatalf("Failed to execute SET command: %v", err)
		}
		if setResponse != "OK" {
			t.Errorf("Expected 'OK', got %q", setResponse)
		}

		// Test GET
		getResponse, err := ts.Client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Test GET for non-existent key
		nonExistentResponse, err := ts.Client.Execute("GET", "nonexistentkey")
		if err != nil {
			t.Fatalf("Failed to execute GET command for non-existent key: %v", err)
		}
		if nonExistentResponse != "" {
			t.Errorf("Expected empty response for non-existent key, got %q", nonExistentResponse)
		}
	})

	// Test SET and GET commands with inline format
	t.Run("SET and GET Commands (Inline Format)", func(t *testing.T) {
		key := "mykey-inline"
		value := "myvalue-inline"

		// Test SET with inline format
		setResponse, err := ts.Client.ExecuteInline("SET", key, value)
		if err != nil {
			t.Fatalf("Failed to execute SET command with inline format: %v", err)
		}
		if setResponse != "OK" {
			t.Errorf("Expected 'OK', got %q", setResponse)
		}

		// Test GET with inline format
		getResponse, err := ts.Client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command with inline format: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Test GET for non-existent key with inline format
		nonExistentResponse, err := ts.Client.ExecuteInline("GET", "nonexistentkey-inline")
		if err != nil {
			t.Fatalf("Failed to execute GET command for non-existent key with inline format: %v", err)
		}
		if nonExistentResponse != "" {
			t.Errorf("Expected empty response for non-existent key, got %q", nonExistentResponse)
		}
	})

	// Test expiry functionality
	t.Run("SET with Expiry EX (RESP Format)", func(t *testing.T) {
		key := "expiringkey-resp"
		value := "expiringvalue-resp"

		// Set with 1 second expiry
		_, err := ts.Client.Execute("SET", key, value, "EX", "1")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry: %v", err)
		}

		// Verify key exists
		getResponse, err := ts.Client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := ts.Client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})

	// Test SET with expiry (inline format)
	t.Run("SET with Expiry EX (Inline Format)", func(t *testing.T) {
		key := "expiringkey-inline"
		value := "expiringvalue-inline"

		// Set with 1 second expiry using inline format
		_, err := ts.Client.ExecuteInline("SET", key, value, "EX", "1")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry using inline format: %v", err)
		}

		// Verify key exists
		getResponse, err := ts.Client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command using inline format: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := ts.Client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry using inline format: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})

	// Test SET with millisecond expiry (RESP format)
	t.Run("SET with Expiry PX (RESP Format)", func(t *testing.T) {
		key := "expiringkey-ms-resp"
		value := "expiringvalue-ms-resp"

		// Set with 1000 milliseconds expiry
		_, err := ts.Client.Execute("SET", key, value, "PX", "1000")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry: %v", err)
		}

		// Verify key exists
		getResponse, err := ts.Client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := ts.Client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})

	// Test SET with millisecond expiry (inline format)
	t.Run("SET with Expiry PX (Inline Format)", func(t *testing.T) {
		key := "expiringkey-ms-inline"
		value := "expiringvalue-ms-inline"

		// Set with 1000 milliseconds expiry using inline format
		_, err := ts.Client.ExecuteInline("SET", key, value, "PX", "1000")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry using inline format: %v", err)
		}

		// Verify key exists
		getResponse, err := ts.Client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command using inline format: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := ts.Client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry using inline format: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})
}
