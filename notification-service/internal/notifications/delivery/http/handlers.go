package http

import "notification-service-api/internal/shared/rpc"

type NotificationHandler struct {
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) SendToTelegram(c *rpc.HttpCtx) {

}
