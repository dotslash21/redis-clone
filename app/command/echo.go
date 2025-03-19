package command

import (
	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/resp"
)

// EchoCommand implements the ECHO command
type EchoCommand struct{}

// NewEchoCommand creates a new ECHO command
func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

// Name returns the command name
func (c *EchoCommand) Name() string {
	return "ECHO"
}

// Execute handles the ECHO command
func (c *EchoCommand) Execute(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New(errors.ErrorTypeCommand, "wrong number of arguments for 'echo' command")
	}
	return resp.FormatBulkString(args[0], false), nil
}
