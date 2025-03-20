package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/dotslash21/redis-clone/app/server"
	"github.com/dotslash21/redis-clone/tests/helpers"
)

// TestSetup encapsulates the server and client setup for tests
type TestSetup struct {
	Server *server.Server
	Client *helpers.RedisClient
	Port   int
}

// NewTestSetup creates a new test setup with a server and client
func NewTestSetup(t *testing.T, port int) *TestSetup {
	// Create and start the Redis server
	srv, err := server.NewServer(port)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		srv.Run()
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a Redis client
	addr := fmt.Sprintf("localhost:%d", port)
	client, err := helpers.NewRedisClient(addr)
	if err != nil {
		srv.Shutdown()
		t.Fatalf("Failed to connect to Redis server: %v", err)
	}

	return &TestSetup{
		Server: srv,
		Client: client,
		Port:   port,
	}
}

// Close shuts down the test setup
func (ts *TestSetup) Close() {
	if ts.Client != nil {
		ts.Client.Close()
	}
	if ts.Server != nil {
		ts.Server.Shutdown()
	}
}
