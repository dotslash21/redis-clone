package helpers

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// RedisClient provides a simplified interface for Redis communications
type RedisClient struct {
	conn   net.Conn
	reader *bufio.Reader
}

// NewRedisClient creates a new Redis client connected to the specified address
func NewRedisClient(address string) (*RedisClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

// Close closes the client connection
func (c *RedisClient) Close() error {
	return c.conn.Close()
}

// Execute sends a command to the Redis server using RESP protocol format
func (c *RedisClient) Execute(command string, args ...string) (string, error) {
	// Format the command in RESP protocol format
	fullArgs := append([]string{command}, args...)
	cmd := fmt.Sprintf("*%d\r\n", len(fullArgs))

	for _, arg := range fullArgs {
		cmd += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}

	// Send the command
	_, err := c.conn.Write([]byte(cmd))
	if err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	// Read the response
	return c.readResponse()
}

// ExecuteInline sends a command to the Redis server using inline protocol format
func (c *RedisClient) ExecuteInline(command string, args ...string) (string, error) {
	// Format the command in simple space-separated format
	parts := []string{command}
	parts = append(parts, args...)

	// Join with spaces and add CRLF
	cmdLine := strings.Join(parts, " ") + "\r\n"

	// Send the command
	_, err := c.conn.Write([]byte(cmdLine))
	if err != nil {
		return "", fmt.Errorf("failed to send inline command: %w", err)
	}

	// Read the response
	return c.readResponse()
}

// readResponse reads and parses a Redis RESP protocol response
func (c *RedisClient) readResponse() (string, error) {
	// Read the first byte to determine the response type
	respType, err := c.reader.ReadByte()
	if err != nil {
		return "", fmt.Errorf("failed to read response type: %w", err)
	}

	// Unread the byte so we can read the full line
	c.reader.UnreadByte()

	// Read the first line
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read response line: %w", err)
	}

	line = strings.TrimRight(line, "\r\n")

	switch respType {
	case '+': // Simple string
		return line[1:], nil

	case '-': // Error
		return "", fmt.Errorf("redis error: %s", line[1:])

	case ':': // Integer
		return line[1:], nil

	case '$': // Bulk string
		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", fmt.Errorf("invalid bulk string length: %w", err)
		}

		// Handle nil response
		if length == -1 {
			return "", nil
		}

		// Read the string content
		data := make([]byte, length+2) // +2 for CRLF
		_, err = c.reader.Read(data)
		if err != nil {
			return "", fmt.Errorf("failed to read bulk string data: %w", err)
		}

		return string(data[:length]), nil

	case '*': // Array
		// Parse the array length
		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", fmt.Errorf("invalid array length: %w", err)
		}

		// If it's an empty array, return a specific format
		if length == 0 {
			return "*0\r\n", nil
		}

		// For arrays, we'll read all elements and return the raw RESP format
		// This will be easier to deal with in tests
		result := line + "\r\n"

		// Read each array element
		for i := 0; i < length; i++ {
			// Read element type
			elemType, err := c.reader.ReadByte()
			if err != nil {
				return "", fmt.Errorf("failed to read array element type: %w", err)
			}

			c.reader.UnreadByte()

			// Read the full element line
			elementLine, err := c.reader.ReadString('\n')
			if err != nil {
				return "", fmt.Errorf("failed to read array element line: %w", err)
			}

			result += elementLine

			// If it's a bulk string, read the data too
			if elemType == '$' {
				bulkLength, err := strconv.Atoi(strings.TrimRight(elementLine, "\r\n")[1:])
				if err != nil {
					return "", fmt.Errorf("invalid bulk string length in array: %w", err)
				}

				if bulkLength > -1 {
					// Read the string content including CRLF
					bulkData := make([]byte, bulkLength+2)
					_, err = c.reader.Read(bulkData)
					if err != nil {
						return "", fmt.Errorf("failed to read bulk string data in array: %w", err)
					}

					result += string(bulkData)
				}
			}
		}

		return result, nil

	default:
		return "", fmt.Errorf("unknown response type: %c", respType)
	}
}
