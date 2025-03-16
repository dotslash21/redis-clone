package command

import (
	"errors"
	"strings"
	"time"
)

func handlePing() string {
	return formatSimpleString("PONG")
}

func handleEcho(args []string) string {
	if len(args) == 0 {
		return formatError("ERR wrong number of arguments for 'ECHO' command")
	}
	return formatSimpleString(args[0])
}

func handleSet(args []string) string {
	if len(args) < 2 {
		return formatError("ERR wrong number of arguments for 'SET' command")
	}
	key, value := args[0], args[1]
	expiry, err := getExpiryDuration(args)
	if err != nil {
		return formatSimpleString(err.Error())
	}

	redis_store.Set(key, value, expiry)
	return formatSimpleString("OK")
}

func getExpiryDuration(args []string) (time.Duration, error) {
	if len(args) < 4 {
		return -1, nil
	}

	expiryType, expiryValue := strings.ToUpper(args[2]), args[3]
	var unit string

	switch expiryType {
	case "EX":
		unit = "s"
	case "PX":
		unit = "ms"
	default:
		return -1, errors.New("ERR invalid expire type in 'SET' command")
	}

	expiry, err := time.ParseDuration(expiryValue + unit)
	if err != nil {
		return -1, errors.New("ERR invalid expire time in 'SET' command")
	}

	return expiry, nil
}

func handleGet(args []string) string {
	if len(args) < 1 {
		return formatError("ERR wrong number of arguments for 'GET' command")
	}
	value, err := redis_store.Get(args[0])
	if err != nil {
		return formatBulkString("", true)
	}
	return formatBulkString(value, false)
}
