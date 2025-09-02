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

Once connected, you'll have access to these tools:

- **`send_email`**: Send emails to any recipient
- **`read_emails`**: Read emails from your inbox with filtering options

## Testing the Connection

### HTTP API Test
```bash
# Start server
./scripts/start-http-server.sh

# Test in another terminal
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/tools
```

### LM Studio Test
Try these prompts once configured:
- "Check my recent emails"
- "Send a test email to myself"
- "List my unread emails"

## Configuration Files

Pre-made configurations are in `lm-studio-configs/`:

- `http-server-config.json` - For HTTP mode (recommended)
- `executable-config.json` - For direct executable mode

## Benefits of HTTP Mode

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