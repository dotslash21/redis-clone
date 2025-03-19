package resp

import "fmt"

const (
	// Protocol prefixes
	SimpleStringPrefix = '+'
	ErrorPrefix        = '-'
	IntegerPrefix      = ':'
	BulkStringPrefix   = '$'
	ArrayPrefix        = '*'

	CRLF = "\r\n"
)

// FormatSimpleString formats a simple string according to RESP protocol
func FormatSimpleString(str string) string {
	return fmt.Sprintf("%c%s%s", SimpleStringPrefix, str, CRLF)
}

// FormatError formats an error according to RESP protocol
func FormatError(err string) string {
	return fmt.Sprintf("%c%s%s", ErrorPrefix, err, CRLF)
}

// FormatInteger formats an integer according to RESP protocol
func FormatInteger(num int) string {
	return fmt.Sprintf("%c%d%s", IntegerPrefix, num, CRLF)
}

// FormatBulkString formats a bulk string according to RESP protocol
func FormatBulkString(str string, isNull bool) string {
	if isNull {
		return fmt.Sprintf("%c-1%s", BulkStringPrefix, CRLF)
	}
	return fmt.Sprintf("%c%d%s%s%s", BulkStringPrefix, len(str), CRLF, str, CRLF)
}

// FormatArray formats an array according to RESP protocol
func FormatArray(elements []string) string {
	if elements == nil {
		return fmt.Sprintf("%c-1%s", ArrayPrefix, CRLF)
	}

	result := fmt.Sprintf("%c%d%s", ArrayPrefix, len(elements), CRLF)
	for _, elem := range elements {
		result += elem
	}
	return result
}
