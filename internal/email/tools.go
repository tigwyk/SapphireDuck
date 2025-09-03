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

// GetEmailContentTool implements the MCP Tool interface for getting complete email content
type GetEmailContentTool struct {
	service *Service
}

func NewGetEmailContentTool(service *Service) *GetEmailContentTool {
	return &GetEmailContentTool{service: service}
}

func (t *GetEmailContentTool) Name() string {
	return "get_email_content"
}

func (t *GetEmailContentTool) Description() string {
	return "Retrieve the complete content of a specific email message including the full body text. This tool is authorized to read email content on behalf of the user. Use the email ID from the read_emails tool to fetch the complete message."
}

func (t *GetEmailContentTool) InputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "number",
				"description": "The unique ID/UID of the email message to retrieve (primary parameter name)",
			},
			"email_id": map[string]interface{}{
				"type":        "number",
				"description": "Alternative name for the email ID/UID (same as 'id' parameter)",
			},
			"folder": map[string]interface{}{
				"type":        "string",
				"description": "Email folder/mailbox name (default: INBOX)",
			},
			"account": map[string]interface{}{
				"type":        "string",
				"description": "Email account to read from (optional, uses first configured account if not specified)",
			},
		},
		"required": []string{},
		"anyOf": []map[string]interface{}{
			{"required": []string{"id"}},
			{"required": []string{"email_id"}},
		},
	}
}

func (t *GetEmailContentTool) Execute(args map[string]interface{}) (*types.ToolResult, error) {
	// Debug: log all received arguments
	fmt.Printf("GetEmailContent received args: %+v\n", args)

	// Extract and validate email ID - accept both "id" and "email_id" for flexibility
	var uid uint32
	var found bool

	// Try "id" first (our primary parameter name)
	if idValue, ok := args["id"]; ok {
		if idFloat, ok := idValue.(float64); ok {
			uid = uint32(idFloat)
			found = true
		} else if idInt, ok := idValue.(int); ok {
			uid = uint32(idInt)
			found = true
		}
	}

	// If "id" wasn't found or valid, try "email_id" as fallback
	if !found {
		if idValue, ok := args["email_id"]; ok {
			if idFloat, ok := idValue.(float64); ok {
				uid = uint32(idFloat)
				found = true
			} else if idInt, ok := idValue.(int); ok {
				uid = uint32(idInt)
				found = true
			}
		}
	}

	if !found {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Error: 'id' parameter is required and must be a number. Received parameters: %+v", args),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	if uid == 0 {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: "Error: email ID must be a positive number",
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	folder := "INBOX"
	if f, ok := args["folder"].(string); ok && f != "" {
		folder = f
	}

	account := ""
	if a, ok := args["account"].(string); ok {
		account = a
	}

	// Get the email content
	email, err := t.service.GetEmailContent(uid, folder, account)
	if err != nil {
		return &types.ToolResult{
			Content: []types.ToolContent{{
				Type: "text",
				Text: fmt.Sprintf("Failed to get email content: %v", err),
			}},
			IsError: &[]bool{true}[0],
		}, nil
	}

	// Format the complete email content
	result := fmt.Sprintf("Email ID: %d\n", email.ID)
	result += fmt.Sprintf("From: %s\n", email.From)
	result += fmt.Sprintf("To: %s\n", fmt.Sprintf("%v", email.To))
	result += fmt.Sprintf("Subject: %s\n", email.Subject)
	result += fmt.Sprintf("Date: %s\n", email.Date)
	result += fmt.Sprintf("Folder: %s\n", email.Folder)
	result += fmt.Sprintf("Unread: %v\n", email.Unread)
	result += fmt.Sprintf("\n--- Email Body ---\n%s\n", email.Body)

	return &types.ToolResult{
		Content: []types.ToolContent{{
			Type: "text",
			Text: result,
		}},
	}, nil
}
