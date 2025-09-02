#!/bin/bash

# Test script for MCP protocol functionality
echo "Testing AI Presence MCP Server"
echo "================================"

# Start the server in the background
./ai-presence-mcp &
SERVER_PID=$!

# Give the server time to start
sleep 2

# Test 1: Initialize
echo "Test 1: Initialize"
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}}' | ./ai-presence-mcp

echo ""

# Test 2: List tools
echo "Test 2: List tools"
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}' | ./ai-presence-mcp

echo ""

# Test 3: Try to call a tool (this will fail without email config, but tests the protocol)
echo "Test 3: Call send_email tool (will fail without config)"
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "send_email", "arguments": {"to": "test@example.com", "subject": "Test", "body": "Test message"}}}' | ./ai-presence-mcp

# Clean up
kill $SERVER_PID 2>/dev/null

echo ""
echo "Tests completed. Check output above for results."