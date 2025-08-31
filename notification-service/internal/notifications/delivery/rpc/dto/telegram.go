package dto

type TelegramRequestSendParams struct {
	To        string  `json:"to" validate:"required"`
	Message   string  `json:"message" validate:"required"`
	ParseMode *string `json:"parse_mode,omitempty"`
}

type TelegramResponseSendDTO struct {
	NotificationID string `json:"notification_id"`
	Queued         bool   `json:"queued"`
}
