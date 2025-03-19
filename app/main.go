package main

import (
	"log"

	"github.com/dotslash21/redis-clone/app/server"
)

func main() {
	srv, err := server.NewServer(6379)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
