package command

import (
	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/resp"
	"github.com/dotslash21/redis-clone/app/store"
)

// GetCommand implements the GET command
type GetCommand struct {
	store *store.Store
}

// NewGetCommand creates a new GET command
func NewGetCommand(s *store.Store) *GetCommand {
	return &GetCommand{store: s}
}

// Name returns the command name
func (c *GetCommand) Name() string {
	return "GET"
}

// Execute handles the GET command
func (c *GetCommand) Execute(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(errors.ErrorTypeCommand, "wrong number of arguments for 'get' command")
	}

	value, err := c.store.Get(args[0])
	if err != nil {
		// Return nil bulk string for non-existent keys
		return resp.FormatBulkString("", true), nil
	}

	return resp.FormatBulkString(value, false), nil
}
