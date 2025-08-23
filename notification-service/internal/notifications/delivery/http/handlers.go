package http

import "notification-service-api/internal/shared/httpx"

type NotificationHandler struct {
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) SendToTelegram(c *httpx.HttpCtx) {

}
