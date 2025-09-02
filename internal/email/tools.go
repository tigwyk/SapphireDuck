package email

import (
	"fmt"

	"ai-presence-mcp/pkg/types"
	"ai-presence-mcp/pkg/utils"
)

// SendEmailTool implements the MCP Tool interface for sending emails
type SendEmailTool struct {
	service *Service
}

func NewSendEmailTool(service *Service) *SendEmailTool {
	return &SendEmailTool{service: service}
}

func (t *SendEmailTool) Name() string {
	return "send_email"
}

func (t *SendEmailTool) Description() string {
	return "Send an email message from the configured email account to a specified recipient. This tool is authorized to send emails on behalf of the user."
}

func (t *SendEmailTool) InputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"to": map[string]interface{}{
				"type":        "string",
				"description": "Email address of the recipient",
			},
			"subject": map[string]interface{}{
				"type":        "string",
				"description": "Subject line of the email",
			},
			"body": map[string]interface{}{
				"type":        "string",
				"description": "Body content of the email",
			},
			"account": map[string]interface{}{
				"type":        "string",
				"description": "Email account to send from (optional, uses first configured account if not specified)",
			},
			"from": map[string]interface{}{
				"type":        "string",
				"description": "Email address to send from (alias for account, optional)",
			},
		},
		"required": []string{"to", "subject", "body"},
	}
}

func (t *SendEmailTool) Execute(args map[string]interface{}) (*types.ToolResult, error) {
	// Debug: log all received arguments
	fmt.Printf("SendEmail received args: %+v\n", args)
	
	to, ok := args["to"].(string)
	if !ok || to == "" {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: "Error: 'to' parameter is required and must be a string",
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	// Validate email format
	if err := utils.ValidateEmail(to); err != nil {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Error: %v", err),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	subject, ok := args["subject"].(string)
	if !ok {
		subject = ""
	}

	// Validate and sanitize subject
	subject = utils.SanitizeInput(subject)
	if err := utils.IsValidSubject(subject); err != nil {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Error: %v", err),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	body, ok := args["body"].(string)
	if !ok {
		body = ""
	}

	// Validate and sanitize body
	body = utils.SanitizeInput(body)
	if err := utils.IsValidBody(body); err != nil {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Error: %v", err),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	account, _ := args["account"].(string)
	if account == "" {
		// Also check "from" parameter as an alias
		account, _ = args["from"].(string)
	}

	if err := t.service.SendEmail(to, subject, body, account); err != nil {
		// Log the specific error for debugging
		fmt.Printf("SendEmail error - to: %s, subject: %s, account: %s, error: %v\n", to, subject, account, err)
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Failed to send email: %v", err),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	return &types.ToolResult{
		Content: []types.ToolContent{{
			Type: "text",
			Text: fmt.Sprintf("Email sent successfully to %s", to),
		}},
	}, nil
}

// ReadEmailsTool implements the MCP Tool interface for reading emails
type ReadEmailsTool struct {
	service *Service
}

func NewReadEmailsTool(service *Service) *ReadEmailsTool {
	return &ReadEmailsTool{service: service}
}

func (t *ReadEmailsTool) Name() string {
	return "read_emails"
}

func (t *ReadEmailsTool) Description() string {
	return "Read and retrieve emails from the configured email account. This tool is authorized to access the configured email accounts and retrieve messages for the user."
}

func (t *ReadEmailsTool) InputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"account": map[string]interface{}{
				"type":        "string",
				"description": "Email account to read from (optional, uses first configured account if not specified)",
			},
			"folder": map[string]interface{}{
				"type":        "string",
				"description": "Folder to read from (optional, defaults to INBOX)",
			},
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of emails to retrieve (optional, defaults to 10)",
			},
			"unread": map[string]interface{}{
				"type":        "boolean",
				"description": "Only retrieve unread emails (optional, defaults to false)",
			},
		},
	}
}

func (t *ReadEmailsTool) Execute(args map[string]interface{}) (*types.ToolResult, error) {
	account, _ := args["account"].(string)
	folder, _ := args["folder"].(string)
	
	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	
	unread := false
	if u, ok := args["unread"].(bool); ok {
		unread = u
	}

	emails, err := t.service.ReadEmails(account, folder, limit, unread)
	if err != nil {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Failed to read emails: %v", err),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	if len(emails) == 0 {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: "No emails found matching the criteria",
			}},
		}, nil
	}

	result := fmt.Sprintf("Found %d email(s):\n\n", len(emails))
	for i, email := range emails {
		result += fmt.Sprintf("%d. From: %s\n   Subject: %s\n   Date: %s\n   Unread: %v\n\n",
			i+1, email.From, email.Subject, email.Date, email.Unread)
	}

	return &types.ToolResult{
		Content: []types.ToolContent{{
			Type: "text",
			Text: result,
		}},
	}, nil
}