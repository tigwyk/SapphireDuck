package mcp

import (
	"encoding/json"
	"testing"

	"ai-presence-mcp/pkg/types"
)

func TestInitialize(t *testing.T) {
	server := NewServer()
	
	initMsg := types.MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
		},
	}
	
	msgBytes, err := json.Marshal(initMsg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	response, err := server.HandleMessage(msgBytes)
	if err != nil {
		t.Fatalf("Failed to handle message: %v", err)
	}
	
	var responseMsg types.MCPMessage
	if err := json.Unmarshal(response, &responseMsg); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if responseMsg.Error != nil {
		t.Errorf("Expected no error, got: %v", responseMsg.Error)
	}
	
	if responseMsg.Result == nil {
		t.Error("Expected result, got nil")
	}
}

func TestToolsList(t *testing.T) {
	server := NewServer()
	
	listMsg := types.MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}
	
	msgBytes, err := json.Marshal(listMsg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	response, err := server.HandleMessage(msgBytes)
	if err != nil {
		t.Fatalf("Failed to handle message: %v", err)
	}
	
	var responseMsg types.MCPMessage
	if err := json.Unmarshal(response, &responseMsg); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if responseMsg.Error != nil {
		t.Errorf("Expected no error, got: %v", responseMsg.Error)
	}
	
	if responseMsg.Result == nil {
		t.Error("Expected result, got nil")
	}
}

func TestInvalidMethod(t *testing.T) {
	server := NewServer()
	
	invalidMsg := types.MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "invalid/method",
	}
	
	msgBytes, err := json.Marshal(invalidMsg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	response, err := server.HandleMessage(msgBytes)
	if err != nil {
		t.Fatalf("Failed to handle message: %v", err)
	}
	
	var responseMsg types.MCPMessage
	if err := json.Unmarshal(response, &responseMsg); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if responseMsg.Error == nil {
		t.Error("Expected error for invalid method, got nil")
	}
	
	if responseMsg.Error.Code != -32601 {
		t.Errorf("Expected error code -32601, got: %d", responseMsg.Error.Code)
	}
}