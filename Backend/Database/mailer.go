package database

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// accountEmailConfig contains the SMTP settings used to send account emails.
type accountEmailConfig struct {
	Host     string
	Port     string
	From     string
	FromName string
	Password string
	SiteURL  string
}

func accountEmailConfigFromEnv() accountEmailConfig {
	return accountEmailConfig{
		Host:     envOrDefault("SMTP_HOST", "smtp.gmail.com"),
		Port:     envOrDefault("SMTP_PORT", "587"),
		From:     strings.TrimSpace(os.Getenv("GMAIL_ADDRESS")),
		FromName: envOrDefault("MAILER_FROM_NAME", "KBF Premier League"),
		Password: strings.TrimSpace(os.Getenv("GMAIL_APP_PASSWORD")),
		SiteURL:  envOrDefault("SITE_URL", "http://localhost:4000/home.html"),
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func sanitizeEmailHeaderValue(value string) string {
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	return strings.TrimSpace(value)
}

func randomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendAccountCreatedEmail(toEmail, firstName, username string) error {
	cfg := accountEmailConfigFromEnv()
	if cfg.From == "" || cfg.Password == "" {
		return fmt.Errorf("mailer is not configured")
	}

	address, err := mail.ParseAddress(toEmail)
	if err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}

	firstName = sanitizeEmailHeaderValue(firstName)
	if firstName == "" {
		firstName = "there"
	}
	username = sanitizeEmailHeaderValue(username)
	if username == "" {
		username = "unknown"
	}

	msg := buildAccountCreatedMessage(cfg, address.Address, firstName, username)
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)
	return smtp.SendMail(addr, auth, cfg.From, []string{address.Address}, msg)
}

func buildAccountCreatedMessage(cfg accountEmailConfig, toEmail, firstName, username string) []byte {
	subject := "Your KBF Premier League account is ready"
	plainBody := fmt.Sprintf(
		"Hi %s,\r\n\r\nYour KBF Premier League account has been created successfully.\r\n\r\nUsername: %s\r\nStatus: Active\r\n\r\nSign in any time at: %s\r\n\r\nIf you didn't create this account, you can safely ignore this email.\r\n\r\n- KBF Premier League\r\n",
		firstName,
		username,
		cfg.SiteURL,
	)

	token, err := randomToken(8)
	if err != nil {
		token = "mailer"
	}

	var msg bytes.Buffer
	fromHeader := fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("UTF-8", cfg.FromName), cfg.From)
	fmt.Fprintf(&msg, "From: %s\r\n", fromHeader)
	fmt.Fprintf(&msg, "To: %s\r\n", toEmail)
	fmt.Fprintf(&msg, "Reply-To: %s\r\n", cfg.From)
	fmt.Fprintf(&msg, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject))
	fmt.Fprintf(&msg, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(&msg, "Message-ID: <%s.%d@%s>\r\n", token, time.Now().UnixNano(), domainOfEmail(cfg.From))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	msg.WriteString(plainBody)
	return msg.Bytes()
}

func domainOfEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return "localhost"
}
