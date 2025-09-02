# AI Presence MCP Server - Tool Documentation

This MCP server provides email functionality through the Model Context Protocol (MCP). The server connects to configured email accounts and exposes tools for sending and reading emails.

## Server Information
- **Name**: AI Presence MCP Server
- **Version**: 0.1.0
- **Protocol Version**: 2024-11-05

## Available Tools

### send_email

**Description**: Send an email message to a specified recipient

**Parameters**:
- `to` (string, required): Email address of the recipient
- `subject` (string, required): Subject line of the email
- `body` (string, required): Body content of the email
- `account` (string, optional): Email account to send from. If not specified, uses the first configured account.

**Example Usage**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "send_email",
    "arguments": {
      "to": "recipient@example.com",
      "subject": "Meeting Reminder",
      "body": "Don't forget about our meeting tomorrow at 2 PM."
    }
  }
}
```

**Success Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Email sent successfully to recipient@example.com"
      }
    ]
  }
}
```

### read_emails

**Description**: Read emails from a specified folder with optional filters

**Parameters**:
- `account` (string, optional): Email account to read from. If not specified, uses the first configured account.
- `folder` (string, optional): Folder to read from. Defaults to "INBOX".
- `limit` (integer, optional): Maximum number of emails to retrieve. Defaults to 10.
- `unread` (boolean, optional): Only retrieve unread emails. Defaults to false.

**Example Usage**:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "read_emails",
    "arguments": {
      "folder": "INBOX",
      "limit": 5,
      "unread": true
    }
  }
}
```

**Success Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Found 2 email(s):\n\n1. From: sender@example.com\n   Subject: Important Update\n   Date: 2025-09-02T19:30:00Z\n   Unread: true\n\n2. From: team@company.com\n   Subject: Weekly Report\n   Date: 2025-09-02T18:15:00Z\n   Unread: true"
      }
    ]
  }
}
```

## Server Configuration

The server requires email configuration in `config.yaml`. Supported email providers include:
- Gmail (with app-specific passwords)
- Generic IMAP/SMTP providers
- PurelyMail (as configured in this instance)

## Error Handling

All tools return standard MCP error responses when operations fail:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "Failed to send email: authentication failed"
  }
}
```

## Usage Notes for AI Assistants

1. **Authentication**: The server handles all email authentication automatically using configured credentials.

2. **Account Selection**: When multiple email accounts are configured, specify the `account` parameter to choose which account to use. Otherwise, the first configured account is used.

3. **Email Formatting**: The `body` parameter accepts plain text. For HTML emails, ensure proper HTML formatting in the body content.

4. **Folder Names**: Common folder names include "INBOX", "Sent", "Drafts", "Trash". Use server-specific folder names as configured by the email provider.

5. **Rate Limiting**: Be mindful of email provider rate limits when sending multiple emails or reading large numbers of messages.

6. **Error Recovery**: If a tool call fails, check the error message for specific details about authentication, network, or configuration issues.