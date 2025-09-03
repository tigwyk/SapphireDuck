# SapphireDuck Integration Guide

## Quick Start for LM Studio

### 🎯 Recommended Method: HTTP Server Mode

**Simplest setup - just use a URL instead of file paths!**

#### Step 1: Start the Server
```bash
# Option A: Use the startup script (recommended)
./scripts/start-http-server.sh

# Option B: Start manually
./sapphire-duck -http
```

#### Step 2: Configure LM Studio
Add this to your LM Studio MCP configuration:

```json
{
  "mcpServers": {
    "sapphire-duck": {
      "url": "http://localhost:8080",
      "description": "AI Presence MCP Server - HTTP mode"
    }
  }
}
```

**That's it!** The HTTP server now implements the full MCP protocol over HTTP, so LM Studio can communicate with it just like a native MCP server, but with the convenience of a simple URL.

### Alternative: Direct Executable Mode
```json
{
  "mcpServers": {
    "sapphire-duck": {
      "command": "/full/path/to/sapphire-duck",
      "args": []
    }
  }
}
```

## Available Tools

Once connected, you'll have access to these email tools:

- **`send_email`**: Send emails to any recipient from configured accounts
- **`read_emails`**: List emails with metadata (from, subject, date, ID, read status)
- **`get_email_content`**: Retrieve complete email content including full body text

## Testing the Connection

### Test Mode
```bash
# Test server configuration and tool registration
./ai-presence-mcp.exe -test
```

### LM Studio Test
Try these prompts once configured:
- "Check my recent emails" - Uses `read_emails`
- "Show me the full content of email ID 12345" - Uses `get_email_content`
- "Send a test email to myself" - Uses `send_email`
- "List my unread emails from the last week" - Uses `read_emails` with filters

### Manual MCP Protocol Test
```bash
# List available tools
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./ai-presence-mcp.exe

# Read email metadata
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "read_emails", "arguments": {"limit": 5}}}' | ./ai-presence-mcp.exe
```

## Email Workflow Example

1. **List emails**: AI calls `read_emails` to see available messages
2. **Get content**: AI calls `get_email_content` with email ID to read full text
3. **Respond**: AI can draft responses using `send_email`

## Configuration Files

The main configuration is in `config.yaml`:
```yaml
email:
  - provider: "your-provider"
    username: "email@domain.com"
    password: "password"
    imap_server: "imap.server.com"
    imap_port: 993
    smtp_server: "smtp.server.com"
    smtp_port: 465
    use_tls: true
```

LM Studio configuration is in `lm-studio-config.json`:
```json
{
  "mcpServers": {
    "ai-presence-mcp": {
      "command": "path/to/ai-presence-mcp.exe",
      "args": [],
      "env": {"CONFIG_PATH": "path/to/config.yaml"},
      "alwaysAllow": ["send_email", "read_emails", "get_email_content"]
    }
  }
}
```

## Benefits of Current Implementation

✅ **Easier setup** - Just a URL, no file paths  
✅ **Better debugging** - Can test endpoints directly  
✅ **More flexible** - Works with any HTTP client  
✅ **Cleaner configs** - No need to specify full executable paths  
✅ **Multi-client** - Multiple AI systems can connect simultaneously  

## Troubleshooting

**Q: Server won't start on port 8080?**  
A: Use a different port: `./sapphire-duck -http -port 8081`

**Q: LM Studio can't connect?**  
A: Make sure the server is running and update the URL in your config.

**Q: No email functionality?**  
A: Create a `config.yaml` file with your email credentials.

---

**Pro Tip:** Keep the HTTP server running in a terminal while using LM Studio for the smoothest experience!