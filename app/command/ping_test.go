package command

import (
	"testing"
)

func TestPingCommand_Name(t *testing.T) {
	cmd := NewPingCommand()
	if cmd.Name() != "PING" {
		t.Errorf("Expected command name to be 'PING', got %s", cmd.Name())
	}
}

func TestPingCommand_Execute(t *testing.T) {
	cmd := NewPingCommand()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "no arguments",
			args:     []string{},
			expected: "+PONG\r\n",
		},
		{
			name:     "with arguments",
			args:     []string{"hello"},
			expected: "+PONG\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cmd.Execute(tt.args)

			// Check error
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Check result
			if result != tt.expected {
				t.Errorf("Expected result %q, got %q", tt.expected, result)
			}
		})
	}
}
