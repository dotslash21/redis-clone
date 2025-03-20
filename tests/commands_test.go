package tests

import (
	"testing"
)

// TestCommands is a simple wrapper to run all Redis command tests
// This file acts as an entry point to the more focused command test files
func TestCommands(t *testing.T) {
	// Start a server on the original test port
	ts := NewTestSetup(t, 16380) // Same port as the original test
	defer ts.Close()

	t.Log("Running Redis command tests. See specific test files for detailed tests:")
	t.Log("- ping_echo_test.go: Tests for PING and ECHO commands")
	t.Log("- set_get_test.go: Tests for SET and GET commands")
	t.Log("- config_test.go: Tests for CONFIG commands")

	// Verify basic connectivity by running a simple PING test
	response, err := ts.Client.Execute("PING")
	if err != nil {
		t.Fatalf("Failed to execute PING command: %v", err)
	}
	if response != "PONG" {
		t.Errorf("Expected 'PONG', got %q", response)
	}
}
