package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"ai-presence-mcp/internal/config"
	"ai-presence-mcp/internal/email"
	"ai-presence-mcp/pkg/types"
)

type Server struct {
	emailService *email.Service
	config       *config.Config
	mux          *http.ServeMux
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type ToolsListResponse struct {
	Tools []ToolInfo `json:"tools"`
}

type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// MCP Protocol types for HTTP transport
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ServerCapabilities     `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}

	// Initialize email service if configured
	if len(cfg.Email) > 0 {
		s.emailService = email.NewService(cfg.Email)
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health and info endpoints
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/info", s.handleInfo)
	
	// MCP protocol endpoints (for LM Studio compatibility)
	s.mux.HandleFunc("/mcp", s.handleMCPRequest)
	s.mux.HandleFunc("/", s.corsMiddleware(s.handleMCPRequest))
	
	// MCP-inspired REST endpoints
	s.mux.HandleFunc("/api/v1/tools", s.handleToolsList)
	
	// Email endpoints
	s.mux.HandleFunc("/api/v1/email/send", s.handleSendEmail)
	s.mux.HandleFunc("/api/v1/email/read", s.handleReadEmails)
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "0.1.0",
	}
	
	s.writeJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"name":        "AI Presence MCP Server",
		"version":     "0.1.0",
		"description": "HTTP API server providing MCP-like functionality for email and communication platforms",
		"endpoints": map[string]string{
			"health":     "/health",
			"tools":      "/api/v1/tools",
			"sendEmail":  "/api/v1/email/send",
			"readEmails": "/api/v1/email/read",
		},
		"emailAccounts": len(s.config.Email),
	}
	
	s.writeJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    info,
	})
}

func (s *Server) handleToolsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	tools := []ToolInfo{}
	
	// Add email tools if email service is available
	if s.emailService != nil {
		sendEmailTool := email.NewSendEmailTool(s.emailService)
		readEmailsTool := email.NewReadEmailsTool(s.emailService)
		
		tools = append(tools, ToolInfo{
			Name:        sendEmailTool.Name(),
			Description: sendEmailTool.Description(),
			InputSchema: sendEmailTool.InputSchema(),
		})
		
		tools = append(tools, ToolInfo{
			Name:        readEmailsTool.Name(),
			Description: readEmailsTool.Description(),
			InputSchema: readEmailsTool.InputSchema(),
		})
	}
	
	s.writeJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    ToolsListResponse{Tools: tools},
	})
}

func (s *Server) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	if s.emailService == nil {
		s.writeJSONError(w, http.StatusServiceUnavailable, "Email service not configured")
		return
	}
	
	var req types.SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	// Validate required fields
	if req.To == "" {
		s.writeJSONError(w, http.StatusBadRequest, "Missing required field: to")
		return
	}
	
	if req.Subject == "" {
		s.writeJSONError(w, http.StatusBadRequest, "Missing required field: subject")
		return
	}
	
	if req.Body == "" {
		s.writeJSONError(w, http.StatusBadRequest, "Missing required field: body")
		return
	}
	
	// Send email
	if err := s.emailService.SendEmail(req.To, req.Subject, req.Body, req.Account); err != nil {
		log.Printf("Failed to send email: %v", err)
		s.writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send email: %v", err))
		return
	}
	
	s.writeJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]string{
			"message": fmt.Sprintf("Email sent successfully to %s", req.To),
		},
	})
}

func (s *Server) handleReadEmails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	if s.emailService == nil {
		s.writeJSONError(w, http.StatusServiceUnavailable, "Email service not configured")
		return
	}
	
	// Parse query parameters
	account := r.URL.Query().Get("account")
	folder := r.URL.Query().Get("folder")
	limitStr := r.URL.Query().Get("limit")
	unreadStr := r.URL.Query().Get("unread")
	
	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	unreadOnly := false
	if unreadStr == "true" {
		unreadOnly = true
	}
	
	// Read emails
	emails, err := s.emailService.ReadEmails(account, folder, limit, unreadOnly)
	if err != nil {
		log.Printf("Failed to read emails: %v", err)
		s.writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read emails: %v", err))
		return
	}
	
	s.writeJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"emails": emails,
			"count":  len(emails),
		},
	})
}

func (s *Server) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeMCPError(w, nil, -32601, "Method not allowed - MCP requires POST", nil)
		return
	}
	
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeMCPError(w, nil, -32700, "Parse error", nil)
		return
	}
	
	switch req.Method {
	case "initialize":
		s.handleMCPInitialize(w, req.ID, req.Params)
	case "tools/list":
		s.handleMCPToolsList(w, req.ID)
	case "tools/call":
		s.handleMCPToolCall(w, req.ID, req.Params)
	default:
		s.writeMCPError(w, req.ID, -32601, "Method not found", nil)
	}
}

func (s *Server) handleMCPInitialize(w http.ResponseWriter, id interface{}, params interface{}) {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "AI Presence MCP Server",
			Version: "0.1.0",
		},
	}
	
	s.writeMCPResponse(w, id, result)
}

func (s *Server) handleMCPToolsList(w http.ResponseWriter, id interface{}) {
	var toolList []ToolInfo
	
	if s.emailService != nil {
		sendEmailTool := email.NewSendEmailTool(s.emailService)
		readEmailsTool := email.NewReadEmailsTool(s.emailService)
		
		toolList = append(toolList, ToolInfo{
			Name:        sendEmailTool.Name(),
			Description: sendEmailTool.Description(),
			InputSchema: sendEmailTool.InputSchema(),
		})
		
		toolList = append(toolList, ToolInfo{
			Name:        readEmailsTool.Name(),
			Description: readEmailsTool.Description(),
			InputSchema: readEmailsTool.InputSchema(),
		})
	}
	
	result := map[string]interface{}{
		"tools": toolList,
	}
	
	s.writeMCPResponse(w, id, result)
}

func (s *Server) handleMCPToolCall(w http.ResponseWriter, id interface{}, params interface{}) {
	if s.emailService == nil {
		s.writeMCPError(w, id, -32603, "Email service not available", nil)
		return
	}
	
	var callParams ToolCallParams
	paramBytes, err := json.Marshal(params)
	if err != nil {
		s.writeMCPError(w, id, -32602, "Invalid params", nil)
		return
	}
	
	if err := json.Unmarshal(paramBytes, &callParams); err != nil {
		s.writeMCPError(w, id, -32602, "Invalid params", nil)
		return
	}
	
	var tool interface {
		Name() string
		Description() string
		InputSchema() interface{}
		Execute(args map[string]interface{}) (*types.ToolResult, error)
	}
	
	switch callParams.Name {
	case "send_email":
		tool = email.NewSendEmailTool(s.emailService)
	case "read_emails":
		tool = email.NewReadEmailsTool(s.emailService)
	default:
		s.writeMCPError(w, id, -32601, fmt.Sprintf("Tool not found: %s", callParams.Name), nil)
		return
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
	
	s.writeMCPResponse(w, id, result)
}

func (s *Server) writeMCPResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) writeMCPError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	// Check if this might be an MCP request to root
	if r.Method == http.MethodPost {
		s.handleMCPRequest(w, r)
		return
	}
	s.writeJSONError(w, http.StatusNotFound, "Endpoint not found")
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) writeJSONError(w http.ResponseWriter, statusCode int, errorMessage string) {
	s.writeJSONResponse(w, statusCode, APIResponse{
		Success: false,
		Error:   errorMessage,
	})
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP API server on %s", addr)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /health - Health check")
	log.Printf("  GET  /api/v1/info - Server information")
	log.Printf("  GET  /api/v1/tools - List available tools")
	log.Printf("  POST /api/v1/email/send - Send email")
	log.Printf("  GET  /api/v1/email/read - Read emails")
	
	return http.ListenAndServe(addr, s.mux)
}