#!/bin/bash
# SapphireDuck HTTP API Examples
# Make sure to start the server with: ./sapphire-duck -http

BASE_URL="http://localhost:8080"

echo "=== SapphireDuck HTTP API Examples ==="
echo

# Health Check
echo "1. Health Check:"
curl -s "$BASE_URL/health" | jq .
echo

# Server Info
echo "2. Server Information:"
curl -s "$BASE_URL/api/v1/info" | jq .
echo

# List Tools
echo "3. Available Tools:"
curl -s "$BASE_URL/api/v1/tools" | jq .
echo

# Read Emails
echo "4. Read Recent Emails (limit 3):"
curl -s "$BASE_URL/api/v1/email/read?limit=3" | jq .
echo

# Send Email Example (commented out to prevent accidental sends)
echo "5. Send Email Example (commented out):"
echo "# curl -X POST \"$BASE_URL/api/v1/email/send\" \\"
echo "#   -H \"Content-Type: application/json\" \\"
echo "#   -d '{\"to\": \"test@example.com\", \"subject\": \"Test from API\", \"body\": \"Hello from SapphireDuck HTTP API!\"}'"
echo

echo "=== Examples Complete ==="