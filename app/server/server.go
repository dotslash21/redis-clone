package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	port              = 6379
	requestBufferSize = 1024
)

// Redis RESP protocol constants
const (
	simpleStringPrefix = '+'
	errorPrefix        = '-'
)

// Run the server and listen for incoming connections
func Run() {
	// Start a TCP server which listens on port 6379
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
		// Read the request from the connection
		request, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client: %s", conn.RemoteAddr())
				return
			}
			log.Printf("Error reading from connection: %v", err)
			return
		}

		request = strings.TrimSpace(request)
		log.Printf("Received data: %s", request)

		// Handle the request
		response := handleCommand(request)
		if response != "" {
			_, err := conn.Write([]byte(response))
			if err != nil {
				log.Printf("Error writing to connection: %v", err)
				return
			}

			log.Printf("Sent response: %s", response)
		}
	}
}

// Handle a command according to the RESP protocol
func handleCommand(command string) string {
	switch strings.ToUpper(command) {
	case "PING":
		return formatSimpleString("PONG")
	default:
		return formatError(fmt.Sprintf("ERR unknown command '%s'", command))
	}
}

// Format a simple string according to RESP protocol
func formatSimpleString(str string) string {
	return fmt.Sprintf("%c%s\r\n", simpleStringPrefix, str)
}

// Format an error according to RESP protocol
func formatError(err string) string {
	return fmt.Sprintf("%c%s\r\n", errorPrefix, err)
}
