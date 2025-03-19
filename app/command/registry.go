package command

import (
	"fmt"
	"sync"

	"github.com/dotslash21/redis-clone/app/errors"
)

// Command represents a Redis command
type Command interface {
	// Name returns the command name
	Name() string
	// Execute executes the command with given arguments
	Execute(args []string) (string, error)
}

// Registry is a thread-safe registry of commands
type Registry struct {
	commands sync.Map
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers a command in the registry
func (r *Registry) Register(cmd Command) error {
	name := cmd.Name()
	// Add logging
	fmt.Printf("Registering command: %q\n", name)

	if _, exists := r.commands.Load(name); exists {
		return errors.New(errors.ErrorTypeCommand, "command already registered")
	}
	r.commands.Store(name, cmd)
	return nil
}

// Get returns a command by name
func (r *Registry) Get(name string) (Command, error) {
	// Add logging
	fmt.Printf("Looking up command: %q\n", name)

	cmd, exists := r.commands.Load(name)
	if !exists {
		return nil, errors.New(errors.ErrorTypeCommand, "command not found")
	}
	return cmd.(Command), nil
}

// Execute executes a command by name with the given arguments
func (r *Registry) Execute(name string, args []string) (string, error) {
	cmd, err := r.Get(name)
	if err != nil {
		return "", err
	}
	return cmd.Execute(args)
}
