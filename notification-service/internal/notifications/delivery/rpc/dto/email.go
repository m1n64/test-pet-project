package dto

import "encoding/base64"

type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"` // base64 encoded
}

type EmailRequestSendParams struct {
	To          string            `json:"to" validate:"required,email"`
	Subject     string            `json:"subject" validate:"required"`
	Body        string            `json:"body" validate:"required"`
	ContentType string            `json:"content_type" validate:"required"`
	ReplyTo     *string           `json:"reply_to" validate:"omitempty,email"`
	From        *string           `json:"from" validate:"omitempty,email"`
	CC          []string          `json:"cc"`
	BCC         []string          `json:"bcc"`
	Attachments []EmailAttachment `json:"attachments"`
}

type EmailRequestSendDTO struct {
	NotificationID string `json:"notification_id"`
	Queued         bool   `json:"queued"`
}

func (a EmailAttachment) Bytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(a.Data)
}
