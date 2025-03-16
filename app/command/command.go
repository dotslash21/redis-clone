package command

import (
	"fmt"
	"strings"

	"github.com/dotslash21/redis-clone/app/store"
)

var redis_store = store.GetStore()

// Handle a command according to the RESP protocol
func HandleCommand(cmdStr string) string {
	cmd, args, err := parseCommand(cmdStr)
	if err != nil {
		return formatError(fmt.Sprintf("ERR %v", err))
	}

	switch strings.ToUpper(cmd) {
	case "PING":
		return handlePing()
	case "ECHO":
		return handleEcho(args)
	case "SET":
		return handleSet(args)
	case "GET":
		return handleGet(args)
	default:
		return formatError(fmt.Sprintf("ERR unknown command '%s'", cmd))
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

// Format a bulk string according to RESP protocol
func formatBulkString(str string, isNull bool) string {
	if isNull {
		return fmt.Sprintf("%c-1\r\n", BULK_STRING_PREFIX)
	}
	return fmt.Sprintf("%c%d\r\n%s\r\n", BULK_STRING_PREFIX, len(str), str)
}
