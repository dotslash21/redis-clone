package command

import (
	"testing"

	"github.com/dotslash21/redis-clone/app/errors"
)

func TestEchoCommand_Name(t *testing.T) {
	cmd := NewEchoCommand()
	if cmd.Name() != "ECHO" {
		t.Errorf("Expected command name to be 'ECHO', got %s", cmd.Name())
	}
}

func TestEchoCommand_Execute(t *testing.T) {
	cmd := NewEchoCommand()

	tests := []struct {
		name     string
		args     []string
		expected string
		errMsg   string
	}{
		{
			name:     "valid echo command",
			args:     []string{"hello"},
			expected: "$5\r\nhello\r\n",
			errMsg:   "",
		},
		{
			name:     "no arguments",
			args:     []string{},
			expected: "",
			errMsg:   "wrong number of arguments for 'echo' command",
		},
		{
			name:     "too many arguments",
			args:     []string{"hello", "world"},
			expected: "",
			errMsg:   "wrong number of arguments for 'echo' command",
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
		})
	}
}
