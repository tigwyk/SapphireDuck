#!/usr/bin/env python3
import json
import subprocess
import sys

def test_mcp_server():
    # Start the server process
    process = subprocess.Popen(
        ['./sapphire-duck'],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    
    # Test initialize
    init_msg = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "test-client", "version": "1.0.0"}
        }
    }
    
    try:
        # Send initialize message
        process.stdin.write(json.dumps(init_msg) + '\n')
        process.stdin.flush()
        
        # Read response (wait a bit)
        import time
        time.sleep(1)
        
        # Test tools list
        tools_msg = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list",
            "params": {}
        }
        
        process.stdin.write(json.dumps(tools_msg) + '\n')
        process.stdin.flush()
        
        time.sleep(1)
        
        # Read any available output
        process.terminate()
        stdout, stderr = process.communicate(timeout=5)
        
        print("STDOUT:", stdout)
        print("STDERR:", stderr)
        
    except Exception as e:
        print(f"Error: {e}")
        process.terminate()

if __name__ == "__main__":
    test_mcp_server()