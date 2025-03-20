package command

import (
	"strings"
	"testing"

	"github.com/dotslash21/redis-clone/app/config"
	"github.com/dotslash21/redis-clone/app/errors"
)

func TestConfigCommand_Name(t *testing.T) {
	cmd := NewConfigCommand()
	if cmd.Name() != "CONFIG" {
		t.Errorf("Expected command name to be 'CONFIG', got %s", cmd.Name())
	}
}

func TestConfigCommand_Execute(t *testing.T) {
	cmd := NewConfigCommand()

	// Setup some initial config for testing GET
	config.SetConfig("test-key", "test-value")
	config.SetConfig("another-key", "another-value")

	tests := []struct {
		name     string
		args     []string
		expected string
		errMsg   string
	}{
		{
			name:     "too few arguments",
			args:     []string{},
			expected: "",
			errMsg:   "wrong number of arguments for 'config' command",
		},
		{
			name:     "invalid subcommand",
			args:     []string{"INVALID"},
			expected: "",
			errMsg:   "unknown subcommand 'config INVALID'",
		},
		{
			name:     "get with too few arguments",
			args:     []string{"GET"},
			expected: "",
			errMsg:   "wrong number of arguments for 'config get' command",
		},
		{
			name:     "get with too many arguments",
			args:     []string{"GET", "test-key", "extra"},
			expected: "",
			errMsg:   "wrong number of arguments for 'config get' command",
		},
		{
			name:     "set with too few arguments",
			args:     []string{"SET", "test-key"},
			expected: "",
			errMsg:   "wrong number of arguments for 'config set' command",
		},
		{
			name:     "set with odd number of arguments",
			args:     []string{"SET", "key1", "value1", "key2"},
			expected: "",
			errMsg:   "wrong number of arguments for 'config set' command",
		},
		{
			name:     "get existing key",
			args:     []string{"GET", "test-key"},
			expected: "*2\r\n$8\r\ntest-key\r\n$10\r\ntest-value\r\n",
			errMsg:   "",
		},
		{
			name:     "get non-existing key",
			args:     []string{"GET", "non-existing-key"},
			expected: "*0\r\n",
			errMsg:   "",
		},
		{
			name:     "get with wildcard",
			args:     []string{"GET", "*key"},
			expected: "*4\r\n$8\r\ntest-key\r\n$10\r\ntest-value\r\n$11\r\nanother-key\r\n$13\r\nanother-value\r\n",
			errMsg:   "",
		},
		{
			name:     "set new key",
			args:     []string{"SET", "new-key", "new-value"},
			expected: "+OK\r\n",
			errMsg:   "",
		},
		{
			name:     "set existing key",
			args:     []string{"SET", "test-key", "updated-value"},
			expected: "+OK\r\n",
			errMsg:   "",
		},
		{
			name:     "set multiple keys",
			args:     []string{"SET", "multi-key1", "multi-value1", "multi-key2", "multi-value2"},
			expected: "+OK\r\n",
			errMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cmd.Execute(tt.args)

			// Check error
			if tt.errMsg != "" {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				cmdErr, ok := err.(*errors.Error)
				if !ok {
					t.Errorf("Expected errors.Error type, got %T", err)
					return
				}
				if cmdErr.Error() != tt.errMsg {
					t.Errorf("Expected error message %q, got %q", tt.errMsg, cmdErr.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
			}

			// Check result
			if result != tt.expected {
				t.Errorf("Expected result %q, got %q", tt.expected, result)
			}

			// Verify SET operations for some cases
			if tt.name == "set new key" {
				result, _ := config.GetConfig("new-key")
				if _, exists := result["new-key"]; !exists || result["new-key"] != "new-value" {
					t.Errorf("Failed to set new key: %v", result)
				}
			} else if tt.name == "set existing key" {
				result, _ := config.GetConfig("test-key")
				if _, exists := result["test-key"]; !exists || result["test-key"] != "updated-value" {
					t.Errorf("Failed to update existing key: %v", result)
				}
			} else if tt.name == "set multiple keys" {
				// Verify multiple keys were set correctly
				result, _ := config.GetConfig("multi-key*")
				if _, exists := result["multi-key1"]; !exists || result["multi-key1"] != "multi-value1" {
					t.Errorf("Failed to set multi-key1: %v", result)
				}
				if _, exists := result["multi-key2"]; !exists || result["multi-key2"] != "multi-value2" {
					t.Errorf("Failed to set multi-key2: %v", result)
				}
			}
		})
	}
}

func TestConfigCommand_Execute_GetWithPattern(t *testing.T) {
	cmd := NewConfigCommand()

	// Setup test data
	config.SetConfig("prefix:key1", "value1")
	config.SetConfig("prefix:key2", "value2")
	config.SetConfig("other:key", "value3")

	// Test with pattern matching
	result, err := cmd.Execute([]string{"GET", "prefix:*"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Since map iteration order is not guaranteed, we need to check if the result contains expected key-value pairs
	if !strings.Contains(result, "prefix:key1") || !strings.Contains(result, "value1") ||
		!strings.Contains(result, "prefix:key2") || !strings.Contains(result, "value2") {
		t.Errorf("Expected result to contain prefix keys and values, got: %s", result)
	}

	if strings.Contains(result, "other:key") || strings.Contains(result, "value3") {
		t.Errorf("Result should not contain non-matching keys, got: %s", result)
	}
}
