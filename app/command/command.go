package command

import (
	"fmt"
	"strings"
)

// Handle a command according to the RESP protocol
func HandleCommand(command string) string {
	command, args, err := parseCommand(command)
	if err != nil {
		return formatError(fmt.Sprintf("ERR %v", err))
	}

	switch strings.ToUpper(command) {
	case "PING":
		return formatSimpleString("PONG")
	case "ECHO":
		if len(args) == 0 {
			return formatError("ERR wrong number of arguments for 'ECHO' command")
		}
		return formatSimpleString(args[0])
	default:
		return formatError(fmt.Sprintf("ERR unknown command '%s'", command))
	}
}

// Parse a command according to the RESP protocol
func parseCommand(command string) (string, []string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("invalid command")
	}

	return parts[0], parts[1:], nil
}

// Format a simple string according to RESP protocol
func formatSimpleString(str string) string {
	return fmt.Sprintf("%c%s\r\n", SIMPLE_STRING_PREFIX, str)
}

// Format an error according to RESP protocol
func formatError(err string) string {
	return fmt.Sprintf("%c%s\r\n", ERROR_PREFIX, err)
}
