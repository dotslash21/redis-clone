package tests

import (
	"testing"
)

// TestPingAndEchoCommands tests the PING and ECHO commands
func TestPingAndEchoCommands(t *testing.T) {
	// Setup test environment
	ts := NewTestSetup(t, 16381) // Different port from other tests
	defer ts.Close()

	// Test PING command with RESP format
	t.Run("PING Command (RESP Format)", func(t *testing.T) {
		response, err := ts.Client.Execute("PING")
		if err != nil {
			t.Fatalf("Failed to execute PING command: %v", err)
		}
		if response != "PONG" {
			t.Errorf("Expected 'PONG', got %q", response)
		}
	})

	// Test PING command with inline format
	t.Run("PING Command (Inline Format)", func(t *testing.T) {
		response, err := ts.Client.ExecuteInline("PING")
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
		response, err := ts.Client.Execute("ECHO", message)
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
		response, err := ts.Client.ExecuteInline("ECHO", message)
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
		response, err := ts.Client.Execute("ECHO", message)
		if err != nil {
			t.Fatalf("Failed to execute ECHO command: %v", err)
		}
		if response != message {
			t.Errorf("Expected %q, got %q", message, response)
		}
	})
}
