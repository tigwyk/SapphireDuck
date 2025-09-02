package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// ValidateEmail validates email address format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

// SanitizeInput removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove null bytes and control characters
	input = strings.ReplaceAll(input, "\x00", "")
	input = regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")
	return strings.TrimSpace(input)
}

// GenerateRandomToken generates a cryptographically secure random token
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// IsValidSubject checks if email subject is reasonable
func IsValidSubject(subject string) error {
	subject = strings.TrimSpace(subject)
	if len(subject) > 998 { // RFC 5322 line length limit
		return fmt.Errorf("subject too long (max 998 characters)")
	}
	return nil
}

// IsValidBody checks if email body is reasonable
func IsValidBody(body string) error {
	if len(body) > 50000 { // Reasonable limit
		return fmt.Errorf("body too long (max 50,000 characters)")
	}
	return nil
}