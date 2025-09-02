# AI Presence MCP Server - Usage Guide

## Overview

This is an MVP implementation of the AI Presence MCP Server with email functionality. The server implements the Model Context Protocol (MCP) specification and provides email tools for AI systems.

## Features Implemented

- ✅ MCP JSON-RPC 2.0 protocol handler
- ✅ Email reading via IMAP
- ✅ Email sending via SMTP
- ✅ Basic input validation and security
- ✅ Configuration system
- ✅ Support for multiple email accounts
- ✅ Gmail and generic IMAP/SMTP support

## Quick Start

### 1. Build the Server

```bash
go build -o ai-presence-mcp
```

### 2. Configure Email Account

Copy the example configuration:
```bash
cp config.example.yaml config.yaml
```

Edit `config.yaml` with your email credentials:
```yaml
server:
  port: 8080
  log_level: "info"

email:
  - provider: "gmail"
    username: "your-email@gmail.com"
    password: "your-app-password"  # Use app-specific password for Gmail
    imap_server: "imap.gmail.com"
    imap_port: 993
    smtp_server: "smtp.gmail.com"
    smtp_port: 587
    use_tls: true
```

**Important for Gmail**: Use an App Password instead of your regular password:
1. Enable 2-factor authentication
2. Generate an App Password: Google Account → Security → App passwords
3. Use the generated 16-character password

### 3. Run the Server

```bash
./ai-presence-mcp
```

The server runs as a stdio-based MCP server, listening on stdin/stdout.

## Available MCP Tools

### 1. send_email

Send an email message.

**Parameters:**
- `to` (required): Recipient email address
- `subject` (required): Email subject
- `body` (required): Email body content
- `account` (optional): Email account to use (defaults to first configured account)

**Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "send_email",
    "arguments": {
      "to": "recipient@example.com",
      "subject": "Hello from MCP",
      "body": "This email was sent via the AI Presence MCP Server!"
    }
  }
}
```

### 2. read_emails

Read emails from a specified folder.

**Parameters:**
- `account` (optional): Email account to read from
- `folder` (optional): Folder name (defaults to "INBOX")
- `limit` (optional): Maximum number of emails (defaults to 10)
- `unread` (optional): Only show unread emails (defaults to false)

**Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "read_emails",
    "arguments": {
      "folder": "INBOX",
      "limit": 5,
      "unread": true
    }
  }
}
```

## Testing

### Run Unit Tests
```bash
go test ./...
```

### Test MCP Protocol Manually

1. **Initialize the connection:**
```bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}}' | ./ai-presence-mcp
```

2. **List available tools:**
```bash
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}' | ./ai-presence-mcp
```

3. **Send a test email:**
```bash
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "send_email", "arguments": {"to": "test@example.com", "subject": "Test", "body": "Test message"}}}' | ./ai-presence-mcp
```

## Security Features

- **Input validation**: Email addresses are validated for proper format
- **Input sanitization**: Content is sanitized to remove control characters
- **Length limits**: Subject and body have reasonable length limits
- **TLS encryption**: All email connections use TLS

## Troubleshooting

### Common Issues

1. **"Tools list is empty"**
   - Check that `config.yaml` exists and has valid email configuration
   - Verify email credentials are correct

2. **Gmail authentication fails**
   - Use App Password instead of regular password
   - Enable 2-factor authentication first
   - Check that IMAP is enabled in Gmail settings

3. **IMAP/SMTP connection fails**
   - Verify server addresses and ports
   - Check firewall settings
   - Ensure TLS settings match server requirements

## Next Steps

This MVP provides the foundation for the full AI Presence MCP Server. Future enhancements include:

- Calendar integration
- Social media platforms
- Advanced authentication (OAuth 2.0)
- Content filtering and safety
- Analytics and monitoring
- Multi-account management

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   MCP Client    │    │   MCP Server     │    │   Email Server  │
│   (AI System)   │◄──►│  (This Project)  │◄──►│   (IMAP/SMTP)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   config.yaml    │
                       └──────────────────┘
```

The server implements the MCP protocol specification and acts as a bridge between AI systems and email servers.