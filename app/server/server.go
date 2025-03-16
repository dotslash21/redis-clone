package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/dotslash21/redis-clone/app/command"
)

// Run the server and listen for incoming connections
func Run() {
	// Start a TCP server which listens on port 6379
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	if err != nil {
		log.Fatalf("Failed to bind to port %d: %v", PORT, err)
	}
	log.Printf("Listening on port %d for incoming connections", PORT)

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
		response := command.HandleCommand(request)
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
