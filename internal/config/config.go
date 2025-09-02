package config

import (
	"fmt"
	"os"

	"ai-presence-mcp/pkg/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Email  []types.EmailConfig `yaml:"email"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}

func Load(configPath string) (*Config, error) {
	// Set defaults
	config := &Config{
		Server: ServerConfig{
			Port:     8080,
			LogLevel: "info",
		},
	}

	// Try to read config file
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if port := os.Getenv("MCP_PORT"); port != "" {
		// Convert port to int if needed
		config.Server.Port = 8080 // Simple default for now
	}

	return config, nil
}

func (c *Config) GetEmailAccount(account string) (*types.EmailConfig, error) {
	if account == "" && len(c.Email) > 0 {
		return &c.Email[0], nil
	}

	for i := range c.Email {
		if c.Email[i].Username == account {
			return &c.Email[i], nil
		}
	}

	return nil, fmt.Errorf("email account not found: %s", account)
}