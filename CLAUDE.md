# SapphireDuck Claude Desktop Integration

## Native MCP Protocol (Best Practice)

SapphireDuck now uses the **proper MCP protocol over stdio**, which is the standard way MCP servers communicate with AI clients. No HTTP, no ports, no URLs - just clean, direct communication.

## Quick Setup

### Step 1: Build the Server
```bash
go build -o sapphire-duck
```

### Step 2: Configure Claude Desktop

Add this to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "sapphire-duck": {
      "command": "<path to sapphire-duck binary>",
      "args": []
    }
  }
}
```

**That's it!** Claude Desktop will manage the MCP server process automatically.

## Configuration File Locations

### macOS
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

### Windows
```
%APPDATA%\Claude\claude_desktop_config.json
```

### Linux
```
~/.config/Claude/claude_desktop_config.json
```

## Configuration Options

### 1. Production Mode (Recommended)

```json
{
  "mcpServers": {
    "sapphire-duck": {
      "command": "G:\\GameDev\\SapphireDuck\\ai-presence-mcp.exe",
      "args": [],
      "env": {
        "CONFIG_PATH": "G:\\GameDev\\SapphireDuck\\config.yaml"
      },
      "description": "SapphireDuck Email MCP Server"
    }
  }
}
```

**Features:**
- ✅ Fast startup (pre-compiled binary)
- ✅ Stable and reliable email operations
- ✅ Three email tools: send, list, and get content
- ✅ Claude manages the process lifecycle
- ✅ Supports multiple email providers

**Requirements:**
- Build first: `go build -o ai-presence-mcp.exe .`
- Configure `config.yaml` with your email credentials
- Update paths to match your installation location

### 2. Development Mode

```json
{
  "mcpServers": {
    "sapphire-duck": {
      "command": "go",
      "args": ["run", "main.go"],
      "cwd": "<path to sapphire-duck binary>",
      "description": "SapphireDuck MCP Server - Development mode"
    }
  }
}
```

**Use Case:** When actively developing and want to test changes without rebuilding

## Why Native MCP Protocol?

### ✅ **Advantages of Native MCP:**
- **Standard Protocol**: Uses the official MCP specification
- **Process Management**: Claude Desktop handles starting/stopping the server
- **Clean Communication**: Direct stdio communication (no network overhead)
- **Security**: No network ports or HTTP endpoints to secure
- **Simplicity**: Just specify a command and args
- **Official Support**: Works with all MCP-compatible clients

### ❌ **What We Removed (HTTP approach):**
- Complex HTTP server setup
- Port management and conflicts
- CORS and web security concerns
- Manual server lifecycle management
- Network communication overhead

## Testing Your Setup

### 1. Test the Server Directly
```bash
# Test functionality
./sapphire-duck -test

## Testing

### 1. Test Server Directly
```bash
# Test configuration and tool registration
./ai-presence-mcp.exe -test

# Test MCP protocol manually
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./ai-presence-mcp.exe
```

### 2. Test with Claude Desktop
Try these prompts once configured:
- "What email tools do you have access to?"
- "Check my recent emails" - Uses `read_emails`
- "Show me the full content of email ID 12345" - Uses `get_email_content`
- "Send a test email to myself" - Uses `send_email`
- "List only my unread emails" - Uses `read_emails` with filters

## Available Tools

Once connected, Claude will have access to three email tools:

- **`send_email`**: Send emails from configured accounts
  - Parameters: to, subject, body, account (optional)
  - Supports multiple accounts and providers

- **`read_emails`**: List emails with metadata  
  - Parameters: account, folder, limit, unread (all optional)
  - Returns: ID, from, to, subject, date, read status
  - Does NOT include full email body (use get_email_content for that)

- **`get_email_content`**: Get complete email content including body text
  - Parameters: id/email_id, folder, account (folder and account optional)
  - Returns: Complete email with full body content
  - Use email ID from read_emails result

## Email Workflow

1. **List emails**: Claude calls `read_emails` to see available messages
2. **Read content**: Claude calls `get_email_content` with specific email ID
3. **Respond**: Claude can compose and send replies using `send_email`

## Troubleshooting

### "Failed to start MCP server"
- **Check Path**: Ensure the `command` path is absolute and points to `ai-presence-mcp.exe`
- **Check Config**: Verify `CONFIG_PATH` environment variable points to valid `config.yaml`
- **Check Build**: Ensure executable was built successfully with `go build -o ai-presence-mcp.exe .`

### "No tools available"
- **Email Config**: Ensure `config.yaml` has valid email configuration section
- **Credentials**: Verify email credentials are correct and account has IMAP/SMTP access
- **Test Mode**: Run `./ai-presence-mcp.exe -test` to verify configuration and tool registration

### "Email operation failed"
- **Authentication**: Check email username/password are correct
- **Ports**: Verify IMAP/SMTP ports and TLS settings match your provider
- **Firewall**: Ensure outbound connections to email servers are allowed
- Use absolute paths in the configuration
- Verify the executable exists: `ls -la /path/to/sapphire-duck`

## Best Practices

1. **Use Absolute Paths**: Always specify full paths in configurations
2. **Build First**: Compile to executable for production use
3. **Test Locally**: Verify with `./sapphire-duck -test` before configuring
4. **Check Logs**: Claude Desktop provides MCP server logs for debugging
5. **Keep It Simple**: Native MCP protocol is cleaner than HTTP approaches

## Example Full Configuration

Create or edit your Claude Desktop config file:

```json
{
  "mcpServers": {
    "sapphire-duck": {
      "command": "<path to sapphire-duck binary>",
      "args": [],
      "env": {
        "LOG_LEVEL": "info"
      }
    }
  },
  "globalShortcut": "Cmd+Shift+."
}
```

## Quick Commands Reference

```bash
# Build the server
go build -o sapphire-duck

# Test functionality
./sapphire-duck -test

# Test MCP protocol
python3 test-native-mcp.py

# Check executable
./sapphire-duck --help  # (if help is implemented)
```

---

**This is the proper way to do MCP!** Native protocol over stdio is the standard, official approach used by all MCP-compatible clients. Much cleaner than HTTP workarounds. 🎯