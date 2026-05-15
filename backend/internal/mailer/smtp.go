package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
)

type Mailer interface {
	Send(to, subject, htmlBody string) error
}

type SMTPMailer struct {
	cfg *config.SMTPConfig
}

func NewSMTPMailer(cfg *config.SMTPConfig) *SMTPMailer {
	if cfg == nil {
		return nil
	}
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Enabled() bool {
	return m != nil && m.cfg != nil && m.cfg.Enabled && m.cfg.Host != "" && m.cfg.Port != "" && m.cfg.From != ""
}

func (m *SMTPMailer) Send(to, subject, htmlBody string) error {
	if !m.Enabled() {
		log.Printf("[MAILER] SMTP is disabled or incomplete config; skipping email to %s with subject %q", to, subject)
		return nil
	}

	headers := map[string]string{
		"From":         m.cfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(k)
		msg.WriteString(": ")
		msg.WriteString(v)
		msg.WriteString("\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	hostPort := netJoinHostPort(m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	if err := smtp.SendMail(hostPort, auth, m.cfg.From, []string{to}, msg.Bytes()); err != nil {
		return fmt.Errorf("smtp send mail failed: %w", err)
	}

	return nil
}

func netJoinHostPort(host, port string) string {
	if strings.Contains(host, ":") {
		return host
	}
	if _, err := strconv.Atoi(port); err == nil {
		return host + ":" + port
	}
	return host + ":" + port
}