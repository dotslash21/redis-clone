package command

import (
	"github.com/dotslash21/redis-clone/app/config"
	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/resp"
)

// ConfigCommand implements the CONFIG command
type ConfigCommand struct {
}

// NewConfigCommand creates a new CONFIG command
func NewConfigCommand() *ConfigCommand {
	return &ConfigCommand{}
}

// Name returns the command name
func (c *ConfigCommand) Name() string {
	return "CONFIG"
}

// Execute handles the CONFIG command
func (c *ConfigCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New(errors.ErrorTypeCommand, "wrong number of arguments for 'config' command")
	}

	subcommand := args[0]
	if subcommand == "GET" {
		if len(args) != 2 {
			return "", errors.New(errors.ErrorTypeCommand, "wrong number of arguments for 'config get' command")
		}

		searchpattern := args[1]
		results, err := config.GetConfig(searchpattern)
		if err != nil {
			return "", errors.New(errors.ErrorTypeCommand, "error getting config: "+err.Error())
		}

		if len(results) == 0 {
			return resp.FormatArray([]string{}), nil
		}

		// Convert the map to a slice of strings for RESP format
		respArray := make([]string, 0, len(results)*2)
		for key, value := range results {
			respArray = append(respArray, resp.FormatBulkString(key, false), resp.FormatBulkString(value, false))
		}

		return resp.FormatArray(respArray), nil
	} else if subcommand == "SET" {
		if len(args) != 3 {
			return "", errors.New(errors.ErrorTypeCommand, "wrong number of arguments for 'config set' command")
		}

		key := args[1]
		value := args[2]

		config.SetConfig(key, value)

		return resp.FormatSimpleString("OK"), nil
	} else {
		return "", errors.New(errors.ErrorTypeCommand, "unknown subcommand 'config "+subcommand+"'")
	}
}
