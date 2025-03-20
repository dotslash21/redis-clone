package tests

import (
	"testing"
	"time"

	"github.com/dotslash21/redis-clone/app/server"
	"github.com/dotslash21/redis-clone/tests/helpers"
)

func TestCommands(t *testing.T) {
	// Start the Redis server on a test port
	testPort := 16380 // Different from other test to avoid conflicts
	srv, err := server.NewServer(testPort)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		srv.Run()
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make sure we shut down the server at the end
	defer srv.Shutdown()

	// Create a Redis client
	client, err := helpers.NewRedisClient("localhost:16380")
	if err != nil {
		t.Fatalf("Failed to connect to Redis server: %v", err)
	}
	defer client.Close()

	// Test PING command with RESP format
	t.Run("PING Command (RESP Format)", func(t *testing.T) {
		response, err := client.Execute("PING")
		if err != nil {
			t.Fatalf("Failed to execute PING command: %v", err)
		}
		if response != "PONG" {
			t.Errorf("Expected 'PONG', got %q", response)
		}
	})

	// Test PING command with inline format
	t.Run("PING Command (Inline Format)", func(t *testing.T) {
		response, err := client.ExecuteInline("PING")
		if err != nil {
			t.Fatalf("Failed to execute PING command with inline format: %v", err)
		}
		if response != "PONG" {
			t.Errorf("Expected 'PONG', got %q", response)
		}
	})

	// Test ECHO command with RESP format
	t.Run("ECHO Command (RESP Format)", func(t *testing.T) {
		message := "Hello,Redis!"
		response, err := client.Execute("ECHO", message)
		if err != nil {
			t.Fatalf("Failed to execute ECHO command: %v", err)
		}
		if response != message {
			t.Errorf("Expected %q, got %q", message, response)
		}
	})

	// Test ECHO command with inline format
	t.Run("ECHO Command (Inline Format)", func(t *testing.T) {
		message := "Hello,Redis!"
		response, err := client.ExecuteInline("ECHO", message)
		if err != nil {
			t.Fatalf("Failed to execute ECHO command with inline format: %v", err)
		}
		if response != message {
			t.Errorf("Expected %q, got %q", message, response)
		}
	})

	// Test ECHO command with spaces (RESP format)
	t.Run("ECHO Command with Spaces (RESP Format)", func(t *testing.T) {
		message := "Hello Redis World with spaces"
		response, err := client.Execute("ECHO", message)
		if err != nil {
			t.Fatalf("Failed to execute ECHO command: %v", err)
		}
		if response != message {
			t.Errorf("Expected %q, got %q", message, response)
		}
	})

	// Test ECHO command with spaces (inline format)
	// t.Run("ECHO Command with Spaces (Inline Format)", func(t *testing.T) {
	// 	message := "Hello Redis World with spaces"
	// 	response, err := client.ExecuteInline("ECHO", message)
	// 	if err != nil {
	// 		t.Fatalf("Failed to execute ECHO command with inline format: %v", err)
	// 	}
	// 	if response != message {
	// 		t.Errorf("Expected %q, got %q", message, response)
	// 	}
	// })

	// Test SET and GET commands with RESP format
	t.Run("SET and GET Commands (RESP Format)", func(t *testing.T) {
		key := "mykey-resp"
		value := "myvalue-resp"

		// Test SET
		setResponse, err := client.Execute("SET", key, value)
		if err != nil {
			t.Fatalf("Failed to execute SET command: %v", err)
		}
		if setResponse != "OK" {
			t.Errorf("Expected 'OK', got %q", setResponse)
		}

		// Test GET
		getResponse, err := client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Test GET for non-existent key
		nonExistentResponse, err := client.Execute("GET", "nonexistentkey")
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
		setResponse, err := client.ExecuteInline("SET", key, value)
		if err != nil {
			t.Fatalf("Failed to execute SET command with inline format: %v", err)
		}
		if setResponse != "OK" {
			t.Errorf("Expected 'OK', got %q", setResponse)
		}

		// Test GET with inline format
		getResponse, err := client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command with inline format: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Test GET for non-existent key with inline format
		nonExistentResponse, err := client.ExecuteInline("GET", "nonexistentkey-inline")
		if err != nil {
			t.Fatalf("Failed to execute GET command for non-existent key with inline format: %v", err)
		}
		if nonExistentResponse != "" {
			t.Errorf("Expected empty response for non-existent key, got %q", nonExistentResponse)
		}
	})

	// Test SET with expiry (RESP format)
	t.Run("SET with Expiry EX (RESP Format)", func(t *testing.T) {
		key := "expiringkey-resp"
		value := "expiringvalue-resp"

		// Set with 1 second expiry
		_, err := client.Execute("SET", key, value, "EX", "1")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry: %v", err)
		}

		// Verify key exists
		getResponse, err := client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := client.Execute("GET", key)
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
		_, err := client.ExecuteInline("SET", key, value, "EX", "1")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry using inline format: %v", err)
		}

		// Verify key exists
		getResponse, err := client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command using inline format: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry using inline format: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})

	// Test SET with expiry (RESP format)
	t.Run("SET with Expiry PX (RESP Format)", func(t *testing.T) {
		key := "expiringkey-resp"
		value := "expiringvalue-resp"

		// Set with 1 second expiry
		_, err := client.Execute("SET", key, value, "PX", "1000")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry: %v", err)
		}

		// Verify key exists
		getResponse, err := client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := client.Execute("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})

	// Test SET with expiry (inline format)
	t.Run("SET with Expiry PX (Inline Format)", func(t *testing.T) {
		key := "expiringkey-inline"
		value := "expiringvalue-inline"

		// Set with 1 second expiry using inline format
		_, err := client.ExecuteInline("SET", key, value, "PX", "1000")
		if err != nil {
			t.Fatalf("Failed to execute SET command with expiry using inline format: %v", err)
		}

		// Verify key exists
		getResponse, err := client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command using inline format: %v", err)
		}
		if getResponse != value {
			t.Errorf("Expected %q, got %q", value, getResponse)
		}

		// Wait for key to expire
		time.Sleep(1500 * time.Millisecond)

		// Verify key is gone
		expiredResponse, err := client.ExecuteInline("GET", key)
		if err != nil {
			t.Fatalf("Failed to execute GET command after expiry using inline format: %v", err)
		}
		if expiredResponse != "" {
			t.Errorf("Expected empty response for expired key, got %q", expiredResponse)
		}
	})
}
