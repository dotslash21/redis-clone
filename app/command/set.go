package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/resp"
	"github.com/dotslash21/redis-clone/app/store"
)

// SetCommand implements the SET command
type SetCommand struct {
	store *store.Store
}

// NewSetCommand creates a new SET command
func NewSetCommand(s *store.Store) *SetCommand {
	return &SetCommand{store: s}
}

// Name returns the command name
func (c *SetCommand) Name() string {
	return "SET"
}

// Execute handles the SET command
func (c *SetCommand) Execute(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New(errors.ErrorTypeCommand, "wrong number of arguments for 'set' command")
	}

	key := args[0]
	value := args[1]
	var ttl time.Duration

	// Process optional arguments (EX, PX)
	if len(args) >= 3 {
		var err error
		ttl, err = parseExpiry(args[2:])
		if err != nil {
			return "", err
		}
	}

	c.store.Set(key, value, ttl)
	return resp.FormatSimpleString("OK"), nil
}

// parseExpiry parses expiry options (EX seconds or PX milliseconds)
func parseExpiry(args []string) (time.Duration, error) {
	if len(args) < 2 {
		return 0, nil
	}

	option := strings.ToUpper(args[0])
	value, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrorTypeCommand, "invalid expire time in 'set' command")
	}

	switch option {
	case "EX":
		return time.Duration(value) * time.Second, nil
	case "PX":
		return time.Duration(value) * time.Millisecond, nil
	default:
		return 0, errors.New(errors.ErrorTypeCommand, "syntax error: invalid expire option, must be EX or PX")
	}
}
