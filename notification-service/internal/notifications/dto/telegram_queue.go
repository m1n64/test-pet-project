package dto

import "time"

type TelegramQueueDTO struct {
	NotificationID string            `msgpack:"notification_id"`
	To             string            `msgpack:"to"`
	Message        string            `msgpack:"message"`
	Vars           map[string]string `msgpack:"vars,omitempty"`
	RequestID      string            `msgpack:"request_id"`
	CreatedAt      time.Time         `msgpack:"created_at"`
	Retries        int               `msgpack:"retries"`
}
