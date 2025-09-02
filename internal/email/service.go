package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"

	"ai-presence-mcp/pkg/types"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Service struct {
	configs []types.EmailConfig
}

func NewService(configs []types.EmailConfig) *Service {
	return &Service{
		configs: configs,
	}
}

func (s *Service) getConfig(account string) (*types.EmailConfig, error) {
	if account == "" && len(s.configs) > 0 {
		return &s.configs[0], nil
	}

	for i := range s.configs {
		if s.configs[i].Username == account {
			return &s.configs[i], nil
		}
	}

	return nil, fmt.Errorf("email account not found: %s", account)
}

func (s *Service) ReadEmails(account, folder string, limit int, unreadOnly bool) ([]types.EmailMessage, error) {
	config, err := s.getConfig(account)
	if err != nil {
		return nil, err
	}

	// Connect to IMAP server
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", config.IMAPServer, config.IMAPPort), &tls.Config{
		ServerName: config.IMAPServer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer c.Close()

	// Login
	if err := c.Login(config.Username, config.Password); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	// Select folder
	if folder == "" {
		folder = "INBOX"
	}
	
	mbox, err := c.Select(folder, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select folder %s: %w", folder, err)
	}

	if mbox.Messages == 0 {
		return []types.EmailMessage{}, nil
	}

	// Determine message range
	var seqSet *imap.SeqSet
	if limit <= 0 {
		limit = 10 // Default limit
	}

	start := uint32(1)
	if mbox.Messages > uint32(limit) {
		start = mbox.Messages - uint32(limit) + 1
	}

	seqSet = new(imap.SeqSet)
	seqSet.AddRange(start, mbox.Messages)

	// Fetch messages
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid}
	messages := make(chan *imap.Message, limit)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var emails []types.EmailMessage
	for msg := range messages {
		if unreadOnly {
			hasUnreadFlag := true
			for _, flag := range msg.Flags {
				if flag == imap.SeenFlag {
					hasUnreadFlag = false
					break
				}
			}
			if !hasUnreadFlag {
				continue
			}
		}

		email := types.EmailMessage{
			ID:      msg.Uid,
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date.Format(time.RFC3339),
			Unread:  !hasFlag(msg.Flags, imap.SeenFlag),
			Folder:  folder,
		}

		if len(msg.Envelope.From) > 0 {
			email.From = msg.Envelope.From[0].Address()
		}

		for _, addr := range msg.Envelope.To {
			email.To = append(email.To, addr.Address())
		}

		emails = append(emails, email)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	return emails, nil
}

func (s *Service) SendEmail(to, subject, body, account string) error {
	config, err := s.getConfig(account)
	if err != nil {
		return err
	}

	// Setup authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPServer)

	// Compose message
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Send email
	addr := fmt.Sprintf("%s:%d", config.SMTPServer, config.SMTPPort)
	if err := smtp.SendMail(addr, auth, config.Username, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func hasFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}