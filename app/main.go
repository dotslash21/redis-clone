package main

import (
	"github.com/dotslash21/redis-clone/app/server"
)

func main() {
	server.Run(6379)
}
