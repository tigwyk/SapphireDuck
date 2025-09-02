# AI Presence MCP Server

![Build Status](https://img.shields.io/badge/build-in%20progress-yellow)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)
![License](https://img.shields.io/badge/license-MIT-green)

A Model Context Protocol (MCP) server built in Go that provides AI systems with authenticated access to email, social media, and communication platforms, enabling autonomous online presence and interaction.

## 🚀 Features

### Core Capabilities
- **Email Integration**: Full IMAP/SMTP support for Gmail, Outlook, and generic providers
- **Social Media Management**: Twitter/X, LinkedIn, and extensible platform support
- **Calendar Management**: Google Calendar, Outlook, and CalDAV integration
- **Messaging Systems**: SMS, WhatsApp Business, Telegram, and more
- **Secure Authentication**: OAuth 2.0, API keys, and encrypted credential storage

### Intelligence & Automation
- **Smart Response System**: Context-aware automated responses
- **Content Safety**: Built-in spam detection and content moderation
- **Analytics Dashboard**: Engagement metrics and productivity insights
- **Multi-Account Support**: Manage multiple accounts across platforms
- **Workflow Automation**: Intelligent scheduling and task management

## 📋 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git for version control
- Valid API credentials for platforms you want to integrate

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/ai-presence-mcp.git
   cd ai-presence-mcp
   ```

2. **Initialize the Go module**
   ```bash
   go mod init ai-presence-mcp
   go mod tidy
   ```

3. **Configure your environment**
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your credentials
   ```

4. **Build and run**
   ```bash
   go build -o ai-presence-mcp
   ./ai-presence-mcp
   ```

## 🔧 Configuration

Create a `config.yaml` file in the project root:

```yaml
server:
  port: 8080
  log_level: "info"

platforms:
  email:
    - provider: "gmail"
      username: "your-email@gmail.com"
      auth_type: "oauth2"
      client_id: "your-client-id"
      client_secret: "your-client-secret"
  
  social:
    twitter:
      api_key: "your-api-key"
      api_secret: "your-api-secret"
      access_token: "your-access-token"
      access_secret: "your-access-secret"
    
    linkedin:
      client_id: "your-linkedin-client-id"
      client_secret: "your-linkedin-client-secret"

security:
  encryption_key: "your-32-byte-encryption-key"
  rate_limit: 100  # requests per minute
```

## 📚 API Documentation

### Core MCP Methods

#### Initialize Connection
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {}
  }
}
```

#### List Available Tools
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

#### Send Email
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "send_email",
    "arguments": {
      "to": "recipient@example.com",
      "subject": "Hello World",
      "body": "This is a test email",
      "account": "your-email@gmail.com"
    }
  }
}
```

### Available Tools

| Tool | Description | Platform |
|------|-------------|----------|
| `send_email` | Send an email message | Email |
| `read_emails` | Retrieve emails with filters | Email |
| `schedule_event` | Create calendar event | Calendar |
| `post_tweet` | Post a tweet | Twitter/X |
| `send_linkedin_post` | Post to LinkedIn | LinkedIn |
| `send_sms` | Send SMS message | SMS |

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   MCP Client    │    │   MCP Server     │    │   Platforms     │
│   (AI System)  │◄──►│  (This Project)  │◄──►│   (APIs)        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   Configuration   │
                       │   & Credentials   │
                       └──────────────────┘
```

### Key Components

- **Protocol Handler**: JSON-RPC 2.0 MCP implementation
- **Authentication Manager**: OAuth 2.0 and API key management
- **Platform Adapters**: Pluggable integrations for each platform
- **Security Layer**: Encryption, rate limiting, and content filtering
- **Analytics Engine**: Usage tracking and insights

## 🛡️ Security

### Built-in Protections
- **Encrypted Credential Storage**: All API keys and tokens are encrypted at rest
- **Rate Limiting**: Prevents API abuse and quota exhaustion
- **Content Filtering**: Blocks inappropriate content and spam
- **Audit Logging**: Comprehensive activity logging for security monitoring
- **Input Validation**: All inputs are validated and sanitized

### Security Best Practices
- Use environment variables for sensitive configuration
- Regularly rotate API keys and tokens
- Monitor audit logs for suspicious activity
- Keep dependencies updated
- Use HTTPS for all communications

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...
```

## 🚧 Development Status

This project is currently in active development. See the [roadmap](ai_presence_mcp_roadmap.md) for detailed progress tracking.

### Current Phase: Foundation & Core Infrastructure
- [x] Project setup and structure
- [x] Basic MCP protocol implementation
- [ ] Authentication framework
- [ ] Configuration system

### Upcoming Phases
1. **Email Integration** (Weeks 4-6)
2. **Calendar Integration** (Weeks 7-8)
3. **Social Media Integration** (Weeks 9-12)
4. **Messaging & Communication** (Weeks 13-14)
5. **Safety & Security** (Weeks 15-16)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 📖 Documentation

- [Configuration Guide](docs/configuration.md)
- [API Reference](docs/api-reference.md)
- [Platform Integration Guides](docs/integrations/)
- [Troubleshooting](docs/troubleshooting.md)
- [Security Guide](docs/security.md)

## 🎯 Use Cases

### Personal Assistant
- Automatically respond to emails while you're away
- Schedule meetings based on calendar availability
- Monitor social media mentions and engage appropriately
- Send SMS reminders for important events

### Business Automation
- Manage customer communications across multiple channels
- Post scheduled content to social media
- Coordinate team communications
- Generate engagement analytics and reports

### AI Research
- Provide AI systems with real-world communication capabilities
- Enable autonomous social media presence
- Facilitate human-AI collaboration in communication tasks

## 🔮 Future Roadmap

- **Mobile Integration**: iOS and Android app support
- **Voice Assistants**: Integration with Siri, Google Assistant
- **Enterprise Features**: Team management and compliance tools
- **Advanced AI**: Personality customization and emotional intelligence
- **Plugin Marketplace**: Community-driven integrations

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) specification
- Go community for excellent libraries and tools
- Contributors and early adopters

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/ai-presence-mcp/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/ai-presence-mcp/discussions)
- **Email**: support@ai-presence-mcp.com
- **Discord**: [Join our community](https://discord.gg/ai-presence-mcp)

---

**⚠️ Disclaimer**: This software provides programmatic access to various online platforms. Users are responsible for complying with each platform's Terms of Service and applicable laws. Use responsibly and ethically.
