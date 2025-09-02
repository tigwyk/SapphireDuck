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
**File:** `claude-desktop-configs/executable-config.json`

```json
{
  "mcpServers": {
    "sapphire-duck": {
      "command": "<path to sapphire-duck binary>",
      "args": [],
      "description": "SapphireDuck MCP Server - Production mode"
    }
  }
}
```

**Pros:**
- ✅ Fast startup (pre-compiled binary)
- ✅ Stable and reliable
- ✅ No compilation step needed
- ✅ Claude Desktop manages the process lifecycle

**Requirements:**
- Must build first: `go build -o sapphire-duck`
- Update path to match your installation location

### 2. Development Mode
**File:** `claude-desktop-configs/go-run-config.json`

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

# Test MCP protocol
python3 test-native-mcp.py
```

### 2. Test with Claude Desktop
Try these prompts once configured:
- "What tools do you have access to?"
- "Check my recent emails"
- "Send a test email to myself"
- "List my unread emails"

## Available Tools

Once connected, Claude will have access to:

- **`send_email`**: Send emails to any recipient
- **`read_emails`**: Read emails from your inbox with filtering options

## Troubleshooting

### "Failed to start MCP server"
- **Check Path**: Ensure the `command` path is absolute and correct
- **Check Permissions**: Make sure the executable has proper permissions
- **Check Dependencies**: Verify `config.yaml` exists with email credentials

### "No tools available"
- **Email Config**: Ensure `config.yaml` has valid email configuration
- **Server Logs**: Check Claude Desktop logs for startup errors
- **Test Mode**: Run `./sapphire-duck -test` to verify functionality

### "Command not found"
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