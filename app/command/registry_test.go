package command

import (
	"testing"

	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/store"
)

// MockCommand implements the Command interface for testing
type MockCommand struct {
	name           string
	executeResult  string
	executeError   error
	executeCount   int
	executeArgSets [][]string
}

func NewMockCommand(name string, result string, err error) *MockCommand {
	return &MockCommand{
		name:           name,
		executeResult:  result,
		executeError:   err,
		executeArgSets: make([][]string, 0),
	}
}

func (c *MockCommand) Name() string {
	return c.name
}

func (c *MockCommand) Execute(args []string) (string, error) {
	c.executeCount++
	c.executeArgSets = append(c.executeArgSets, args)
	return c.executeResult, c.executeError
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	// Register a valid command
	mockCmd := NewMockCommand("TEST", "result", nil)
	err := registry.Register(mockCmd)
	if err != nil {
		t.Errorf("Expected no error when registering a new command, got %v", err)
	}

	// Try to register the same command again (should fail)
	mockCmd2 := NewMockCommand("TEST", "another result", nil)
	err = registry.Register(mockCmd2)
	if err == nil {
		t.Errorf("Expected error when registering a duplicate command, got nil")
	}

	// Register a different command
	mockCmd3 := NewMockCommand("ANOTHER", "result", nil)
	err = registry.Register(mockCmd3)
	if err != nil {
		t.Errorf("Expected no error when registering a different command, got %v", err)
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	mockCmd := NewMockCommand("TEST", "result", nil)
	registry.Register(mockCmd)

	tests := []struct {
		name         string
		commandName  string
		shouldFound  bool
		expectedType interface{}
	}{
		{
			name:         "existing command",
			commandName:  "TEST",
			shouldFound:  true,
			expectedType: &MockCommand{},
		},
		{
			name:        "non-existing command",
			commandName: "NONEXISTENT",
			shouldFound: false,
		},
		{
			name:        "empty command name",
			commandName: "",
			shouldFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := registry.Get(tt.commandName)

			if tt.shouldFound {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
				if cmd == nil {
					t.Errorf("Expected a command instance, got nil")
					return
				}

				// Check that we got the expected type
				switch tt.expectedType.(type) {
				case *MockCommand:
					if _, ok := cmd.(*MockCommand); !ok {
						t.Errorf("Expected *MockCommand, got %T", cmd)
					}
				default:
					t.Errorf("Unhandled expected type %T", tt.expectedType)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if cmd != nil {
					t.Errorf("Expected nil, got %v", cmd)
				}
			}
		})
	}
}

func TestRegistry_Execute(t *testing.T) {
	registry := NewRegistry()

	// Register commands
	successCmd := NewMockCommand("SUCCESS", "ok", nil)
	registry.Register(successCmd)

	errorCmd := NewMockCommand("ERROR", "", errors.New(errors.ErrorTypeCommand, "test error"))
	registry.Register(errorCmd)

	// Set up a real command for end-to-end testing
	storeInstance := store.GetStore()
	setCmd := NewSetCommand(storeInstance)
	registry.Register(setCmd)
	getCmd := NewGetCommand(storeInstance)
	registry.Register(getCmd)

	tests := []struct {
		name        string
		commandName string
		args        []string
		expected    string
		errMsg      string
	}{
		{
			name:        "execute successful command",
			commandName: "SUCCESS",
			args:        []string{"arg1", "arg2"},
			expected:    "ok",
			errMsg:      "",
		},
		{
			name:        "execute command with error",
			commandName: "ERROR",
			args:        []string{"arg1"},
			expected:    "",
			errMsg:      "test error",
		},
		{
			name:        "execute non-existent command",
			commandName: "NONEXISTENT",
			args:        []string{},
			expected:    "",
			errMsg:      "command not found",
		},
		{
			name:        "end-to-end set and get",
			commandName: "SET",
			args:        []string{"testkey", "testvalue"},
			expected:    "+OK\r\n",
			errMsg:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := registry.Execute(tt.commandName, tt.args)

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

	// Finally, verify the end-to-end test by getting the value that was set
	result, err := registry.Execute("GET", []string{"testkey"})
	if err != nil {
		t.Errorf("Expected no error from GET after SET, got %v", err)
		return
	}
	expectedResult := "$9\r\ntestvalue\r\n"
	if result != expectedResult {
		t.Errorf("Expected result %q from GET after SET, got %q", expectedResult, result)
	}
}
