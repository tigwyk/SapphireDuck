package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"ai-presence-mcp/internal/config"
	"ai-presence-mcp/internal/email"
	"ai-presence-mcp/internal/mcp"
)

func Run() error {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("Warning: Could not load config file, using defaults: %v", err)
		cfg = &config.Config{
			Server: config.ServerConfig{
				Port:     8080,
				LogLevel: "info",
			},
		}
	}

	log.Printf("Starting AI Presence MCP Server...")

	// Create MCP server
	server := mcp.NewServer()

	// Register email tools if email config is available
	if len(cfg.Email) > 0 {
		emailService := email.NewService(cfg.Email)
		
		sendEmailTool := email.NewSendEmailTool(emailService)
		readEmailsTool := email.NewReadEmailsTool(emailService)
		
		server.RegisterTool(sendEmailTool)
		server.RegisterTool(readEmailsTool)
		
		log.Printf("Registered email tools for %d accounts", len(cfg.Email))
	}

	// Start stdio server (MCP protocol over stdin/stdout)
	log.Printf("MCP Server ready. Listening on stdin/stdout...")
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		response, err := server.HandleMessage([]byte(line))
		if err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}

		fmt.Println(string(response))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	return nil
}