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
