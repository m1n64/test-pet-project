package rpc

import (
	"notification-service-api/internal/notifications/delivery/rpc/dto"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/pkg/di"
)

func InitNotificationProcedures(dependencies *di.Dependencies) {
	notificationHandler := NewNotificationHandler(dependencies.Validator, dependencies.TelegramService)

	dependencies.Registry.Register("telegram.send", rpc.Typed[dto.TelegramRequestSendParams](notificationHandler.SendToTelegram))
}
