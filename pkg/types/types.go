package types

// MCP Protocol Types

type MCPMessage struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      *ClientInfo            `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type ToolResult struct {
	Content []ToolContent `json:"content"`
	IsError *bool         `json:"isError,omitempty"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Email Types

type EmailConfig struct {
	Provider     string `yaml:"provider"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	IMAPServer   string `yaml:"imap_server"`
	IMAPPort     int    `yaml:"imap_port"`
	SMTPServer   string `yaml:"smtp_server"`
	SMTPPort     int    `yaml:"smtp_port"`
	UseTLS       bool   `yaml:"use_tls"`
}

type EmailMessage struct {
	ID          uint32   `json:"id"`
	From        string   `json:"from"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Date        string   `json:"date"`
	Unread      bool     `json:"unread"`
	Folder      string   `json:"folder"`
}

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Account string `json:"account,omitempty"`
}

type ReadEmailsRequest struct {
	Account  string `json:"account,omitempty"`
	Folder   string `json:"folder,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Unread   *bool  `json:"unread,omitempty"`
}