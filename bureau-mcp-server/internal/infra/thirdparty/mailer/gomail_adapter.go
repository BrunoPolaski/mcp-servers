package mailer

import (
	"context"
	"crypto/tls"
	"io"
	"os"
	"strconv"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"gopkg.in/gomail.v2"
)

type goMailAdapter struct {
	dialer   *gomail.Dialer
	from     string
	fromName string
}

func NewGoMailAdapter() *goMailAdapter {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("SMTP_FROM_NAME")

	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587
	}

	dialer := gomail.NewDialer(host, port, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &goMailAdapter{
		dialer:   dialer,
		from:     username,
		fromName: fromName,
	}
}

func (s *goMailAdapter) Send(ctx context.Context, message *entities.Mail) *rest_err.RestErr {
	m := gomail.NewMessage()

	from := s.from
	if s.fromName != "" {
		from = m.FormatAddress(s.from, s.fromName)
	}
	m.SetHeader("From", from)

	if len(message.To) > 0 {
		m.SetHeader("To", message.To)
	}

	m.SetHeader("Subject", message.Subject)

	if message.IsHTML {
		m.SetBody("text/html", message.Body)
	} else {
		m.SetBody("text/plain", message.Body)
	}

	if len(message.Cc) > 0 {
		m.SetHeader("Cc", message.Cc...)
	}

	if len(message.Bcc) > 0 {
		m.SetHeader("Bcc", message.Bcc...)
	}

	// Add attachments
	for _, attachment := range message.Attachments {
		m.Attach(attachment.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Content)
			return err
		}))
	}

	// Add embedded images
	for _, embedded := range message.Embedded {
		m.Embed(embedded.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(embedded.Content)
			return err
		}), gomail.SetHeader(map[string][]string{
			"Content-ID": {"<" + embedded.CID + ">"},
		}))
	}

	err := s.dialer.DialAndSend(m)
	if err != nil {
		return rest_err.NewInternalServerError("Error sending email: %s", err.Error()).WithCause(err)
	}

	return nil
}
