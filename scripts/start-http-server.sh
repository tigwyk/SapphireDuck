#!/bin/bash

# SapphireDuck HTTP Server Startup Script

echo "🦆 Starting SapphireDuck HTTP Server..."
echo "📍 Server will be available at: http://localhost:8080"
echo

# Check if binary exists
if [ ! -f "./sapphire-duck" ]; then
    echo "❌ Error: sapphire-duck binary not found"
    echo "   Please run: go build -o sapphire-duck"
    exit 1
fi

# Check if config exists
if [ ! -f "./config.yaml" ]; then
    echo "⚠️  Warning: config.yaml not found - server will use defaults"
    echo "   Create config.yaml for email functionality"
    echo
fi

echo "🚀 Starting server in HTTP mode..."
echo "   Use Ctrl+C to stop"
echo
echo "📋 Available endpoints:"
echo "   GET  /health                    - Health check"
echo "   GET  /api/v1/info               - Server info"
echo "   GET  /api/v1/tools              - List tools"
echo "   POST /api/v1/email/send         - Send email"
echo "   GET  /api/v1/email/read         - Read emails"
echo
echo "🔗 AI Client Config (Claude Desktop / LM Studio):"
echo "   URL: http://localhost:8080"
echo "   See CLAUDE.md and INTEGRATION.md for setup details"
echo "═══════════════════════════════════════════════"

./sapphire-duck -http