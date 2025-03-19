package command

import (
	"testing"
	"time"

	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/store"
)

func TestSetCommand_Name(t *testing.T) {
	cmd := NewSetCommand(store.GetStore())
	if cmd.Name() != "SET" {
		t.Errorf("Expected command name to be 'SET', got %s", cmd.Name())
	}
}

func TestSetCommand_Execute(t *testing.T) {
	storeInstance := store.GetStore()
	cmd := NewSetCommand(storeInstance)

	tests := []struct {
		name     string
		args     []string
		expected string
		errMsg   string
		verify   func(t *testing.T)
	}{
		{
			name:     "basic set command",
			args:     []string{"key1", "value1"},
			expected: "+OK\r\n",
			errMsg:   "",
			verify: func(t *testing.T) {
				value, err := storeInstance.Get("key1")
				if err != nil {
					t.Errorf("Expected key to exist, got error: %v", err)
					return
				}
				if value != "value1" {
					t.Errorf("Expected value to be %q, got %q", "value1", value)
				}
			},
		},
		{
			name:     "set with expiry in seconds",
			args:     []string{"key2", "value2", "EX", "1"},
			expected: "+OK\r\n",
			errMsg:   "",
			verify: func(t *testing.T) {
				// Verify the value was set
				value, err := storeInstance.Get("key2")
				if err != nil {
					t.Errorf("Expected key to exist, got error: %v", err)
					return
				}
				if value != "value2" {
					t.Errorf("Expected value to be %q, got %q", "value2", value)
				}

				// Wait for expiry
				time.Sleep(1100 * time.Millisecond)

				// Verify the key expired
				_, err = storeInstance.Get("key2")
				if err == nil {
					t.Errorf("Expected key to be expired")
				}
			},
		},
		{
			name:     "set with expiry in milliseconds",
			args:     []string{"key3", "value3", "PX", "100"},
			expected: "+OK\r\n",
			errMsg:   "",
			verify: func(t *testing.T) {
				// Verify the value was set
				value, err := storeInstance.Get("key3")
				if err != nil {
					t.Errorf("Expected key to exist, got error: %v", err)
					return
				}
				if value != "value3" {
					t.Errorf("Expected value to be %q, got %q", "value3", value)
				}

				// Wait for expiry
				time.Sleep(150 * time.Millisecond)

				// Verify the key expired
				_, err = storeInstance.Get("key3")
				if err == nil {
					t.Errorf("Expected key to be expired")
				}
			},
		},
		{
			name:     "too few arguments",
			args:     []string{"key4"},
			expected: "",
			errMsg:   "wrong number of arguments for 'set' command",
		},
		{
			name:     "invalid expiry option",
			args:     []string{"key5", "value5", "INVALID", "10"},
			expected: "",
			errMsg:   "syntax error: invalid expire option, must be EX or PX",
		},
		{
			name:     "invalid expiry value",
			args:     []string{"key6", "value6", "EX", "notanumber"},
			expected: "",
			errMsg:   "", // We'll just check that there's an error, not the exact message
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
			} else if tt.name == "invalid expiry value" {
				// Special case for this test since the exact error message is different
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				cmdErr, ok := err.(*errors.Error)
				if !ok {
					t.Errorf("Expected errors.Error type, got %T", err)
					return
				}
				// Just check that the error is a command error
				if cmdErr.Type != errors.ErrorTypeCommand {
					t.Errorf("Expected error type %v, got %v", errors.ErrorTypeCommand, cmdErr.Type)
					return
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

			// Run verify function if provided
			if tt.verify != nil {
				tt.verify(t)
			}
		})
	}
}

func TestParseExpiry(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedTTL time.Duration
		expectError bool
	}{
		{
			name:        "EX with valid seconds",
			args:        []string{"EX", "10"},
			expectedTTL: 10 * time.Second,
			expectError: false,
		},
		{
			name:        "PX with valid milliseconds",
			args:        []string{"PX", "500"},
			expectedTTL: 500 * time.Millisecond,
			expectError: false,
		},
		{
			name:        "EX with invalid value",
			args:        []string{"EX", "abc"},
			expectError: true,
		},
		{
			name:        "PX with invalid value",
			args:        []string{"PX", "abc"},
			expectError: true,
		},
		{
			name:        "Invalid option",
			args:        []string{"XX", "10"},
			expectError: true,
		},
		{
			name:        "Empty args",
			args:        []string{},
			expectedTTL: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttl, err := parseExpiry(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}

				if ttl != tt.expectedTTL {
					t.Errorf("Expected TTL %v, got %v", tt.expectedTTL, ttl)
				}
			}
		})
	}
}
