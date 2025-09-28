package entity

import (
	"github.com/google/uuid"
	"time"
)

type TelegramNotification struct {
	NotificationID uuid.UUID `msgpack:"notification_id"`
	CorrelationID  string    `msgpack:"request_id"`
	To             string    `msgpack:"to"`
	Payload        string    `msgpack:"payload"`
	ParseMode      string    `msgpack:"parse_mode"`
	CreatedAt      time.Time `msgpack:"created_at"`
}

type EmailAttachment struct {
	Filename    string `msgpack:"filename"`
	Data        []byte `msgpack:"data"`
	ContentType string `msgpack:"content_type"`
}

type EmailNotification struct {
	NotificationID uuid.UUID         `msgpack:"notification_id"`
	CorrelationID  string            `msgpack:"request_id"`
	To             string            `msgpack:"to"`
	Subject        string            `msgpack:"subject"`
	Body           string            `msgpack:"body"`
	ContentType    string            `msgpack:"content_type"`
	CC             []string          `msgpack:"cc"`
	BCC            []string          `msgpack:"bcc"`
	From           *string           `msgpack:"from"`
	ReplyTo        *string           `msgpack:"reply_to"`
	Attachments    []EmailAttachment `msgpack:"attachments"`
	CreatedAt      time.Time         `msgpack:"created_at"`
	SentAt         time.Time         `msgpack:"sent_at"`
}
