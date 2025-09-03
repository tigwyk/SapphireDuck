# SapphireDuck MCP Server - Usage Guide

## Overview

SapphireDuck is a Model Context Protocol (MCP) server that provides comprehensive email functionality to AI systems. The server implements the official MCP specification using the Go SDK and provides three main email tools.

## Features Implemented

- ✅ MCP JSON-RPC 2.0 protocol with official Go SDK
- ✅ Email reading with metadata retrieval (IMAP)
- ✅ Complete email content fetching (IMAP with RFC822)
- ✅ Email sending with SSL/TLS support (SMTP)
- ✅ Multi-provider support (Gmail, PurelyMail, generic IMAP/SMTP)
- ✅ Robust input validation and error handling
- ✅ YAML configuration system
- ✅ Multiple email account support
- ✅ Flexible parameter handling for AI model compatibility

## Quick Start

### 1. Build the Server

```bash
go build -o ai-presence-mcp.exe .
```

### 2. Configure Email Account

Create your `config.yaml` file:
```yaml
email:
  - provider: "purelymail"
    username: "your-email@domain.com"
    password: "your-password"
    imap_server: "imap.purelymail.com"
    imap_port: 993
    smtp_server: "smtp.purelymail.com" 
    smtp_port: 465
    use_tls: true
```

**For Gmail**, use an App Password:
```yaml
email:
  - provider: "gmail"
    username: "your-email@gmail.com"
    password: "your-app-password"  # Use app-specific password
    imap_server: "imap.gmail.com"
    imap_port: 993
    smtp_server: "smtp.gmail.com"
    smtp_port: 587
    use_tls: true
```

**Gmail App Password Setup**:
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

## Available Tools

### 1. `send_email`
Send an email from a configured account.

**Parameters:**
- `to` (required): Recipient email address
- `subject` (required): Email subject line
- `body` (required): Email content (plain text)
- `account` (optional): Email account to use (defaults to first configured)
- `from` (optional): Alias for account parameter

### 2. `read_emails`
Retrieve a list of emails with metadata (does not include full body content).

**Parameters:**
- `account` (optional): Email account to read from
- `folder` (optional): Folder name (defaults to "INBOX")
- `limit` (optional): Maximum number of emails (defaults to 10)
- `unread` (optional): Only show unread emails (defaults to false)

**Returns:** List with ID, from, to, subject, date, unread status, folder

### 3. `get_email_content`
Retrieve the complete content of a specific email including full body text.

**Parameters:**
- `id` OR `email_id` (required): Email UID from read_emails
- `folder` (optional): Folder name (defaults to "INBOX")
- `account` (optional): Email account to use

**Returns:** Complete email with full body content

## Testing

### Test Mode
```bash
./ai-presence-mcp.exe -test
```
This verifies configuration and tool registration without starting the MCP server.

### MCP Client Testing

1. **Use with LM Studio:**
   - Configure `lm-studio-config.json` with server path
   - Restart LM Studio to load the MCP server
   - Ask AI to check emails or send messages

2. **Manual MCP Protocol Testing:**
```bash
# Initialize connection
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}}' | ./ai-presence-mcp.exe

# List available tools
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}' | ./ai-presence-mcp.exe

# Read emails
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "read_emails", "arguments": {"limit": 3, "unread": true}}}' | ./ai-presence-mcp.exe

# Get specific email content
echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "get_email_content", "arguments": {"id": 12345}}}' | ./ai-presence-mcp.exe
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