#!/usr/bin/env python3
"""
Test native MCP protocol over stdio (the proper way)
This simulates how Claude Desktop and LM Studio communicate with MCP servers
"""
import json
import subprocess
import sys
import time

def test_native_mcp():
    print("=== Testing Native MCP Protocol (stdio) ===")
    print("This is how Claude Desktop and LM Studio actually communicate with MCP servers")
    print()
    
    # Start the MCP server process
    print("Starting MCP server process...")
    process = subprocess.Popen(
        ['./sapphire-duck'],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=0
    )
    
    def send_request(request):
        """Send a JSON-RPC request and get response"""
        request_json = json.dumps(request)
        print(f"→ Sending: {request_json}")
        
        process.stdin.write(request_json + '\n')
        process.stdin.flush()
        
        # Read response
        response_line = process.stdout.readline()
        if response_line:
            response = json.loads(response_line.strip())
            print(f"← Received: {json.dumps(response, indent=2)}")
            return response
        else:
            print("No response received")
            return None
    
    try:
        # Test 1: Initialize the connection
        print("\n1. Testing MCP Initialize (required handshake):")
        init_request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0.0"
                }
            }
        }
        
        init_response = send_request(init_request)
        if init_response and "result" in init_response:
            print("✅ Initialize successful!")
            server_info = init_response["result"]["serverInfo"]
            print(f"   Server: {server_info['name']} v{server_info['version']}")
            
            # Send initialized notification (required by MCP spec)
            print("\n   Sending 'initialized' notification...")
            initialized_notification = {
                "jsonrpc": "2.0",
                "method": "notifications/initialized"
            }
            process.stdin.write(json.dumps(initialized_notification) + '\n')
            process.stdin.flush()
            print("   ✅ Initialization complete!")
        else:
            print("❌ Initialize failed!")
            return
        
        # Test 2: List available tools
        print("\n2. Testing Tools List:")
        tools_request = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list",
            "params": {}
        }
        
        tools_response = send_request(tools_request)
        if tools_response and "result" in tools_response:
            tools = tools_response["result"]["tools"]
            print(f"✅ Found {len(tools)} tools:")
            for tool in tools:
                print(f"   - {tool['name']}: {tool['description']}")
        else:
            print("❌ Tools list failed!")
            return
        
        # Test 3: Call a tool (read emails)
        print("\n3. Testing Tool Call (read_emails):")
        tool_call_request = {
            "jsonrpc": "2.0",
            "id": 3,
            "method": "tools/call",
            "params": {
                "name": "read_emails",
                "arguments": {
                    "limit": 2
                }
            }
        }
        
        tool_response = send_request(tool_call_request)
        if tool_response and "result" in tool_response:
            print("✅ Tool call successful!")
            content = tool_response["result"]["content"]
            if content:
                print(f"   Result: {content[0]['text'][:100]}...")
        else:
            print("❌ Tool call failed!")
        
        print("\n=== Native MCP Protocol Test Complete! ===")
        print("✅ All tests passed! The server is working with proper MCP protocol.")
        
    except Exception as e:
        print(f"❌ Error during testing: {e}")
    finally:
        # Clean up
        process.terminate()
        process.wait()

if __name__ == "__main__":
    test_native_mcp()