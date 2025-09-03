# SapphireDuck MCP Server - Tool Documentation

This MCP server provides comprehensive email functionality through the Model Context Protocol (MCP). The server connects to configured email accounts and exposes three main tools for AI systems to manage email operations.

## Server Information
- **Name**: SapphireDuck MCP Server
- **Version**: 1.0.0
- **Protocol Version**: 2024-11-05
- **Go SDK Version**: v0.3.1
- **Email Library**: go-mail v0.6.2

## Available Tools

### send_email

**Description**: Send an email message from a configured account to a specified recipient. Supports multiple email providers with SSL/TLS security.

**Parameters**:
- `to` (string, required): Email address of the recipient
- `subject` (string, required): Subject line of the email
- `body` (string, required): Body content of the email (plain text)
- `account` (string, optional): Email account to send from. If not specified, uses the first configured account.
- `from` (string, optional): Alias for the account parameter

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

**Description**: Retrieve a list of emails with metadata from a specified folder. Returns email envelope information (from, to, subject, date, ID, read status) but does NOT include the full email body content. Use `get_email_content` to fetch complete email content.

**Parameters**:
- `account` (string, optional): Email account to read from. If not specified, uses the first configured account.
- `folder` (string, optional): Folder to read from. Defaults to "INBOX".
- `limit` (integer, optional): Maximum number of emails to retrieve. Defaults to 10.
- `unread` (boolean, optional): Only retrieve unread emails. Defaults to false.

**Returns**: List of emails with ID, from, to, subject, date, unread status, folder (NO body content)

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

### get_email_content

**Description**: Retrieve the complete content of a specific email message including the full body text. Use the email ID from the `read_emails` tool to fetch the complete message.

**Parameters**:
- `id` OR `email_id` (number, required): The unique ID/UID of the email message to retrieve (either parameter name works)
- `folder` (string, optional): Email folder/mailbox name (default: "INBOX")
- `account` (string, optional): Email account to read from. If not specified, uses the first configured account.

**Example Usage**:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_email_content",
    "arguments": {
      "id": 12345,
      "folder": "INBOX"
    }
  }
}
```

**Alternative Usage** (using `email_id` parameter):
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_email_content",
    "arguments": {
      "email_id": 12345,
      "folder": "INBOX"
    }
  }
}
```

**Success Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Email ID: 12345\nFrom: sender@example.com\nTo: ['recipient@example.com']\nSubject: Important Update\nDate: 2025-09-02T19:30:00Z\nFolder: INBOX\nUnread: true\n\n--- Email Body ---\nHello,\n\nThis is the complete email content with all the body text.\n\nBest regards,\nSender"
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

3. **Email Reading Workflow**: 
   - Use `read_emails` to get a list of emails with metadata (from, subject, date, ID)
   - Use `get_email_content` with the email ID from step 1 to retrieve the complete email content including body text
   - This two-step approach optimizes performance and allows filtering before fetching full content

4. **Email Formatting**: The `body` parameter accepts plain text. For HTML emails, ensure proper HTML formatting in the body content.

5. **Folder Names**: Common folder names include "INBOX", "Sent", "Drafts", "Trash". Use server-specific folder names as configured by the email provider.

6. **Rate Limiting**: Be mindful of email provider rate limits when sending multiple emails or reading large numbers of messages.

7. **Error Recovery**: If a tool call fails, check the error message for specific details about authentication, network, or configuration issues.