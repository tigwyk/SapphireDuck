package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"ai-presence-mcp/pkg/types"
)

type Server struct {
	mcpServer *sdkmcp.Server
}

type Tool interface {
	Name() string
	Description() string
	InputSchema() interface{}
	Execute(args map[string]interface{}) (*types.ToolResult, error)
}

func NewServer() *Server {
	mcpServer := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "ai-presence-mcp",
		Version: "0.1.0",
	}, nil)
	
	return &Server{
		mcpServer: mcpServer,
	}
}

func (s *Server) RegisterTool(tool Tool) {
	toolDef := &sdkmcp.Tool{
		Name:        tool.Name(),
		Description: tool.Description(),
		InputSchema: nil,
	}
	
	handler := func(ctx context.Context, req *sdkmcp.CallToolRequest, rawArgs any) (*sdkmcp.CallToolResult, any, error) {
		args := make(map[string]interface{})
		if arguments, ok := rawArgs.(map[string]interface{}); ok {
			args = arguments
		}
		
		result, err := tool.Execute(args)
		if err != nil {
			return &sdkmcp.CallToolResult{
				Content: []sdkmcp.Content{&sdkmcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}, nil, nil
		}
		
		var content []sdkmcp.Content
		for _, c := range result.Content {
			content = append(content, &sdkmcp.TextContent{Text: c.Text})
		}
		
		isError := result.IsError != nil && *result.IsError
		return &sdkmcp.CallToolResult{
			Content: content,
			IsError: isError,
		}, nil, nil
	}
	
	sdkmcp.AddTool(s.mcpServer, toolDef, handler)
}

func (s *Server) Run(ctx context.Context, transport sdkmcp.Transport) error {
	return s.mcpServer.Run(ctx, transport)
}