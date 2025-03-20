package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dotslash21/redis-clone/app/command"
	"github.com/dotslash21/redis-clone/app/errors"
	"github.com/dotslash21/redis-clone/app/store"
)

// Server represents a Redis server
type Server struct {
	listener  net.Listener
	registry  *command.Registry
	store     *store.Store
	conns     sync.Map
	shutdown  chan struct{}
	waitGroup sync.WaitGroup
}

// NewServer creates a new Redis server
func NewServer(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrorTypeServer, fmt.Sprintf("failed to bind to port %d", port))
	}

	s := &Server{
		listener: listener,
		registry: command.NewRegistry(),
		store:    store.GetStore(),
		shutdown: make(chan struct{}),
	}

	// Register commands
	s.registerCommands()

	return s, nil
}

// registerCommands registers all supported Redis commands
func (s *Server) registerCommands() {
	s.registry.Register(command.NewPingCommand())
	s.registry.Register(command.NewEchoCommand())
	s.registry.Register(command.NewSetCommand(s.store))
	s.registry.Register(command.NewGetCommand(s.store))
	s.registry.Register(command.NewConfigCommand())
}

// Run starts the server and listens for connections
func (s *Server) Run() error {
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start accepting connections
	go s.acceptConnections()

	// Wait for shutdown signal
	<-sigChan
	return s.Shutdown()
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Printf("Error accepting connection: %v", err)
				}
				return
			}

			s.waitGroup.Add(1)
			s.conns.Store(conn.RemoteAddr(), conn)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles a client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.conns.Delete(conn.RemoteAddr())
		s.waitGroup.Done()
	}()

	reader := bufio.NewReader(conn)

	for {
		select {
		case <-s.shutdown:
			return
		default:
			// Set read deadline to handle hanging connections
			// conn.SetReadDeadline(time.Now().Add(time.Second))

			cmd, args, err := parseRESP(reader)
			if err != nil {
				if err == io.EOF {
					return
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Error parsing command: %v", err)
				continue
			}

			response, err := s.registry.Execute(cmd, args)
			if err != nil {
				if errors.IsCommandError(err) {
					log.Printf("Command error executing %s: %v", cmd, err)
					response = fmt.Sprintf("-ERR %v\r\n", err)
				} else {
					log.Printf("Internal error executing %s: %v", cmd, err)
					response = "-ERR internal server error\r\n"
				}
			}

			if _, err = conn.Write([]byte(response)); err != nil {
				log.Printf("Error writing to connection: %v", err)
				return
			}
		}
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	// Signal shutdown
	close(s.shutdown)

	// Close listener to stop accepting new connections
	if err := s.listener.Close(); err != nil {
		return errors.Wrap(err, errors.ErrorTypeServer, "failed to close listener")
	}

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for all connections to close or timeout
	done := make(chan struct{})
	go func() {
		s.waitGroup.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Printf("Server shutdown complete")
	case <-ctx.Done():
		log.Printf("Server shutdown timed out")
	}

	return nil
}

// parseRESP reads a RESP array from the client and returns the command and args.
func parseRESP(reader *bufio.Reader) (string, []string, error) {
	line, err := readLine(reader)
	if err != nil {
		return "", nil, err
	}

	log.Printf("RESP parsing first line: %q", line)

	// Check if the command is an inline command
	if len(line) != 0 && !strings.HasPrefix(line, "*") {
		log.Printf("Parsing as inline command")
		return parseInlineCommand(line)
	}

	log.Printf("Parsing as RESP array")
	cmd, args, err := parseRESPArray(reader, line)
	if err != nil {
		return "", nil, err
	}

	log.Printf("Parsed command: %s, args: %v", cmd, args)
	return cmd, args, nil
}

// readLine reads a line and trims CRLF
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// parseRESPArray parses a RESP array and returns command and arguments
func parseRESPArray(reader *bufio.Reader, firstLine string) (string, []string, error) {
	// Parse array length
	arrayLen, err := parseArrayLength(firstLine)
	if err != nil {
		return "", nil, err
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
		return "", nil, errors.New(errors.ErrorTypeCommand, "empty command")
	}

	return strings.ToUpper(parts[0]), parts[1:], nil
}

// parseArrayLength parses the array length from a RESP array header
func parseArrayLength(line string) (int, error) {
	if !strings.HasPrefix(line, "*") {
		return 0, errors.New(errors.ErrorTypeCommand, "expected array")
	}

	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, errors.New(errors.ErrorTypeCommand, "invalid array length")
	}

	return length, nil
}

// readBulkString reads a RESP bulk string
func readBulkString(reader *bufio.Reader) (string, error) {
	// Read the length line
	lengthLine, err := readLine(reader)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrorTypeCommand, "failed to read bulk string length line")
	}

	if !strings.HasPrefix(lengthLine, "$") {
		return "", errors.New(errors.ErrorTypeCommand, "expected bulk string")
	}

	// Parse the length
	strLen, err := strconv.Atoi(lengthLine[1:])
	if err != nil {
		return "", errors.New(errors.ErrorTypeCommand, "invalid bulk string length")
	}

	// Handle null bulk string
	if strLen == -1 {
		return "", nil
	}

	// Read the string data
	argBytes := make([]byte, strLen)
	if _, err := io.ReadFull(reader, argBytes); err != nil {
		return "", errors.Wrap(err, errors.ErrorTypeCommand, "failed to read bulk string data")
	}

	// Discard the trailing CRLF
	if _, err := reader.Discard(2); err != nil {
		return "", errors.Wrap(err, errors.ErrorTypeCommand, "failed to read bulk string CRLF")
	}

	return string(argBytes), nil
}

// parseInlineCommand parses a simple inline command from the client
func parseInlineCommand(line string) (string, []string, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil, errors.New(errors.ErrorTypeCommand, "empty command")
	}

	return strings.ToUpper(parts[0]), parts[1:], nil
}
