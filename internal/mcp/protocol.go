package mcp

import (
	"encoding/json"
	"fmt"
	"log"

	"ai-presence-mcp/pkg/types"
)

const ProtocolVersion = "2024-11-05"

type Server struct {
	tools map[string]Tool
}

type Tool interface {
	Name() string
	Description() string
	InputSchema() interface{}
	Execute(args map[string]interface{}) (*types.ToolResult, error)
}

func NewServer() *Server {
	return &Server{
		tools: make(map[string]Tool),
	}
}

func (s *Server) RegisterTool(tool Tool) {
	s.tools[tool.Name()] = tool
}

func (s *Server) HandleMessage(data []byte) ([]byte, error) {
	var msg types.MCPMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return s.errorResponse(nil, -32700, "Parse error", nil)
	}

	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg.ID, msg.Params)
	case "tools/list":
		return s.handleToolsList(msg.ID)
	case "tools/call":
		return s.handleToolCall(msg.ID, msg.Params)
	default:
		return s.errorResponse(msg.ID, -32601, "Method not found", nil)
	}
}

func (s *Server) handleInitialize(id interface{}, params interface{}) ([]byte, error) {
	result := types.InitializeResult{
		ProtocolVersion: ProtocolVersion,
		Capabilities: types.ServerCapabilities{
			Tools: &types.ToolsCapability{},
		},
		ServerInfo: types.ServerInfo{
			Name:    "AI Presence MCP Server",
			Version: "0.1.0",
		},
	}

	return s.successResponse(id, result)
}

func (s *Server) handleToolsList(id interface{}) ([]byte, error) {
	var toolList []types.Tool
	
	for _, tool := range s.tools {
		toolList = append(toolList, types.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.InputSchema(),
		})
	}

	result := map[string]interface{}{
		"tools": toolList,
	}

	return s.successResponse(id, result)
}

func (s *Server) handleToolCall(id interface{}, params interface{}) ([]byte, error) {
	var callParams types.ToolCallParams
	
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return s.errorResponse(id, -32602, "Invalid params", nil)
	}
	
	if err := json.Unmarshal(paramBytes, &callParams); err != nil {
		return s.errorResponse(id, -32602, "Invalid params", nil)
	}

	tool, exists := s.tools[callParams.Name]
	if !exists {
		return s.errorResponse(id, -32601, fmt.Sprintf("Tool not found: %s", callParams.Name), nil)
	}

	result, err := tool.Execute(callParams.Arguments)
	if err != nil {
		log.Printf("Tool execution error: %v", err)
		isError := true
		result = &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Error executing tool: %v", err),
			}},
			IsError: &isError,
		}
	}

	return s.successResponse(id, result)
}

func (s *Server) successResponse(id interface{}, result interface{}) ([]byte, error) {
	response := types.MCPMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	return json.Marshal(response)
}

func (s *Server) errorResponse(id interface{}, code int, message string, data interface{}) ([]byte, error) {
	response := types.MCPMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &types.MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return json.Marshal(response)
}