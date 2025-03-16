package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/dotslash21/redis-clone/app/command"
)

const (
	CRLF = "\r\n"
)

// Run the server and listen for incoming connections on the specified port
func Run(port int) {
	// Start a TCP server which listens on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to bind to port %d: %v", port, err)
	}
	log.Printf("Listening on port %d for incoming connections", port)

	for {
		// Accept a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue // Try to accept another connection instead of crashing
		}

		// Handle the connection
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

// Handle a connection
func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure connection is closed when function exits

	reader := bufio.NewReader(conn)

	for {
		cmd, args, err := parseRESP(reader)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client")
				return
			}
			log.Printf("Error parsing command: %v", err)
			continue // Continue to read the next command
		}

		response := command.HandleCommand(cmd, args)
		if _, err = conn.Write([]byte(response)); err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}

// parseRESP reads a RESP array from the client and returns the command and args.
func parseRESP(reader *bufio.Reader) (string, []string, error) {
	line, err := readLine(reader)
	if err != nil {
		return "", nil, err
	}

	// Check if the command is an inline command
	if len(line) != 0 && rune(line[0]) != command.ARRAY_PREFIX {
		return parseInlineCommand(line)
	}

	return parseRESPArray(reader, line)
}

// readLine reads a line and trims CRLF
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, CRLF), nil
}

// parseRESPArray parses a RESP array and returns command and arguments
func parseRESPArray(reader *bufio.Reader, firstLine string) (string, []string, error) {
	// Parse array length
	arrayLen, err := strconv.Atoi(firstLine[1:])
	if err != nil {
		return "", nil, fmt.Errorf("invalid array length: %w", err)
	}

	// Read each bulk string
	parts := make([]string, 0, arrayLen)
	for i := 0; i < arrayLen; i++ {
		// Read the bulk string
		str, err := readBulkString(reader)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, str)
	}

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}

	return parts[0], parts[1:], nil
}

// readBulkString reads a RESP bulk string
func readBulkString(reader *bufio.Reader) (string, error) {
	// Read the length line
	lengthLine, err := readLine(reader)
	if err != nil {
		return "", err
	}

	if len(lengthLine) == 0 || rune(lengthLine[0]) != command.BULK_STRING_PREFIX {
		return "", fmt.Errorf("expected bulk string prefix '%c'", command.BULK_STRING_PREFIX)
	}

	// Parse the length
	strLen, err := strconv.Atoi(lengthLine[1:])
	if err != nil {
		return "", fmt.Errorf("invalid bulk string length: %w", err)
	}

	// Read the string data
	argBytes := make([]byte, strLen)
	if _, err := io.ReadFull(reader, argBytes); err != nil {
		return "", err
	}

	// Discard the trailing CRLF
	if _, err := reader.Discard(2); err != nil {
		return "", err
	}

	return string(argBytes), nil
}

// parseInlineCommand parses a simple inline command from the client.
func parseInlineCommand(line string) (string, []string, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}

	return parts[0], parts[1:], nil
}
