package utils

import (
	"bytes"
	"context"
	mail "github.com/wneessen/go-mail"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type SMTPClient struct {
	client *mail.Client
	from   string
}

type MailAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

type MailMessage struct {
	To          string
	Subject     string
	Body        string
	ContentType string
	ReplyTo     *string
	From        *string
	CC          []string
	BCC         []string
	Attachments []MailAttachment
}

func NewSMTPClient(host string, port string, tlsMode string, username string, password string, from string, logger *zap.Logger) *SMTPClient {
	smtpPort, _ := strconv.Atoi(port)

	client, err := mail.NewClient(
		host,
		mail.WithPort(smtpPort),
		mail.WithTimeout(10*time.Second),
		mail.WithTLSPolicy(tlsPolicy(tlsMode)),
	)
	if err != nil {
		logger.Fatal("Failed to create SMTP client", zap.Error(err))
		return nil
	}

	if username != "" && password != "" {
		client.SetSMTPAuth(mail.SMTPAuthLogin)
		client.SetUsername(username)
		client.SetPassword(password)
	}

	return &SMTPClient{
		client: client,
		from:   from,
	}
}

func (s *SMTPClient) Send(message *MailMessage) error {
	msg := mail.NewMsg()

	fromEmail := s.from
	if message.From != nil {
		fromEmail = *message.From
	}

	if err := msg.From(fromEmail); err != nil {
		return err
	}

	if err := msg.To(message.To); err != nil {
		return err
	}

	if message.ReplyTo != nil {
		if err := msg.ReplyTo(*message.ReplyTo); err != nil {
			return err
		}
	}

	if len(message.CC) > 0 {
		if err := msg.Cc(message.CC...); err != nil {
			return err
		}
	}

	if len(message.BCC) > 0 {
		if err := msg.Bcc(message.BCC...); err != nil {
			return err
		}
	}

	msg.Subject(message.Subject)
	msg.SetBodyString(parseBodyContentType(message.ContentType), message.Body)

	for _, a := range message.Attachments {
		if err := msg.AttachReader(a.Filename, bytes.NewReader(a.Data), mail.WithFileContentType(mail.ContentType(a.ContentType))); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.client.DialAndSendWithContext(ctx, msg); err != nil {
		return err
	}

	return nil
}

func tlsPolicy(mode string) mail.TLSPolicy {
	switch mode {
	case "required":
		return mail.TLSMandatory
	case "starttls":
		return mail.TLSOpportunistic
	default:
		return mail.NoTLS
	}
}

func parseBodyContentType(ct string) mail.ContentType {
	switch ct {
	case "text/html":
		return mail.TypeTextHTML
	case "text/plain":
		fallthrough
	default:
		return mail.TypeTextPlain
	}
}
