package command

// Redis RESP protocol constants
const (
	SIMPLE_STRING_PREFIX = '+'
	ERROR_PREFIX         = '-'
	INTEGER_PREFIX       = ':'
	BULK_STRING_PREFIX   = '$'
	ARRAY_PREFIX         = '*'
)
