package main

import (
	"log"
	"os"

	"ai-presence-mcp/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}