package dto

type TelegramRequestSendDTO struct {
	To      string            `json:"to" validate:"required"`
	Message string            `json:"message" validate:"required"`
	Vars    map[string]string `json:"vars,omitempty"`
}

type TelegramResponseSendDTO struct {
	NotificationID string `json:"notification_id"`
	Queued         bool   `json:"queued"`
}
