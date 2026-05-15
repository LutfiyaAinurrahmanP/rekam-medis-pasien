package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
)

type Mailer interface {
	Send(to, subject, htmlBody string) error
	SendWithInlineImage(to, subject, htmlBody, contentID, filename, contentType string, imageData []byte) error
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

func (m *SMTPMailer) SendWithInlineImage(to, subject, htmlBody, contentID, filename, contentType string, imageData []byte) error {
	if !m.Enabled() {
		log.Printf("[MAILER] SMTP is disabled or incomplete config; skipping email to %s with subject %q", to, subject)
		return nil
	}

	var msg bytes.Buffer
	boundary := "mixed-boundary-" + strings.ReplaceAll(contentID, " ", "-")
	writer := multipart.NewWriter(&msg)
	if err := writer.SetBoundary(boundary); err != nil {
		return fmt.Errorf("smtp set boundary failed: %w", err)
	}

	headers := map[string]string{
		"From":         m.cfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type":  `multipart/related; boundary="` + boundary + `"`,
	}

	for k, v := range headers {
		msg.WriteString(k)
		msg.WriteString(": ")
		msg.WriteString(v)
		msg.WriteString("\r\n")
	}
	msg.WriteString("\r\n")

	htmlPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/html; charset=UTF-8"},
		"Content-Transfer-Encoding": {"7bit"},
	})
	if err != nil {
		return fmt.Errorf("smtp create html part failed: %w", err)
	}
	if _, err := htmlPart.Write([]byte(htmlBody)); err != nil {
		return fmt.Errorf("smtp write html part failed: %w", err)
	}

	imagePart, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {contentType},
		"Content-Transfer-Encoding": {"base64"},
		"Content-ID":                {"<" + contentID + ">"},
		"Content-Disposition":       {`inline; filename="` + filename + `"`},
	})
	if err != nil {
		return fmt.Errorf("smtp create image part failed: %w", err)
	}

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(imageData)))
	base64.StdEncoding.Encode(encoded, imageData)
	for len(encoded) > 76 {
		if _, err := imagePart.Write(encoded[:76]); err != nil {
			return fmt.Errorf("smtp write image part failed: %w", err)
		}
		if _, err := imagePart.Write([]byte("\r\n")); err != nil {
			return fmt.Errorf("smtp write image part failed: %w", err)
		}
		encoded = encoded[76:]
	}
	if len(encoded) > 0 {
		if _, err := imagePart.Write(encoded); err != nil {
			return fmt.Errorf("smtp write image part failed: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp close multipart writer failed: %w", err)
	}

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