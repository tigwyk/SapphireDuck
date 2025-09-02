package main

import (
	"flag"
	"log"
	"os"

	"ai-presence-mcp/cmd/server"
)

func main() {
	testMode := flag.Bool("test", false, "Run in test mode to verify MCP server functionality")
	flag.Parse()

	if err := server.Run(*testMode); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}