package command

import (
	"github.com/dotslash21/redis-clone/app/resp"
)

// PingCommand implements the PING command
type PingCommand struct{}

// NewPingCommand creates a new PING command
func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

// Name returns the command name
func (c *PingCommand) Name() string {
	return "PING"
}

// Execute handles the PING command
func (c *PingCommand) Execute(args []string) (string, error) {
	return resp.FormatSimpleString("PONG"), nil
}
