package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"ai-presence-mcp/internal/config"
	"ai-presence-mcp/internal/email"
	"ai-presence-mcp/internal/mcp"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func Run(testMode bool) error {
	// Configure logger to use stderr (stdout must be reserved for JSON-RPC in MCP)
	log.SetOutput(os.Stderr)

	// Get config path from environment variable or use default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Load configuration
	cfg, err := config.Load(configPath)
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

	if testMode {
		return runTestMode(server, cfg)
	}

	// Start MCP server with stdio transport (standard MCP protocol)
	log.Printf("MCP Server ready. Listening on stdin/stdout...")

	transport := &sdkmcp.StdioTransport{}
	ctx := context.Background()

	if err := server.Run(ctx, transport); err != nil {
		return fmt.Errorf("failed to run MCP server: %w", err)
	}

	return nil
}

func runTestMode(server *mcp.Server, cfg *config.Config) error {
	log.Printf("Running in test mode...")

	if len(cfg.Email) > 0 {
		log.Printf("✅ MCP Server initialized with %d email accounts", len(cfg.Email))
	} else {
		log.Printf("⚠️  No email configuration found")
	}

	log.Printf("✅ Test mode completed! Server would be ready to accept MCP connections.")
	return nil
}
