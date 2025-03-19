package command

import (
	"testing"
	"time"

	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/store"
)

func TestGetCommand_Name(t *testing.T) {
	cmd := NewGetCommand(store.GetStore())
	if cmd.Name() != "GET" {
		t.Errorf("Expected command name to be 'GET', got %s", cmd.Name())
	}
}

func TestGetCommand_Execute(t *testing.T) {
	storeInstance := store.GetStore()
	cmd := NewGetCommand(storeInstance)

	// Set up some test data
	storeInstance.Set("testkey", "testvalue", 0)
	storeInstance.Set("expiring", "expirevalue", 100*time.Millisecond)

	tests := []struct {
		name     string
		args     []string
		expected string
		errMsg   string
		setup    func()
	}{
		{
			name:     "valid get command",
			args:     []string{"testkey"},
			expected: "$9\r\ntestvalue\r\n",
			errMsg:   "",
		},
		{
			name:     "key does not exist",
			args:     []string{"nonexistent"},
			expected: "$-1\r\n",
			errMsg:   "",
		},
		{
			name:     "no arguments",
			args:     []string{},
			expected: "",
			errMsg:   "wrong number of arguments for 'get' command",
		},
		{
			name:     "too many arguments",
			args:     []string{"key1", "key2"},
			expected: "",
			errMsg:   "wrong number of arguments for 'get' command",
		},
		{
			name:     "expired key returns nil",
			args:     []string{"expiring"},
			expected: "$-1\r\n",
			errMsg:   "",
			setup: func() {
				// Sleep to ensure the key expires
				time.Sleep(150 * time.Millisecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute setup if provided
			if tt.setup != nil {
				tt.setup()
			}

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
		})
	}
}
