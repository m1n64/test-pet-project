package rpc

import (
	"github.com/gin-gonic/gin"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/pkg/di"
)

func InitNotificationRoutes(r *gin.Engine, dependencies *di.Dependencies) {
	notificationHandler := NewNotificationHandler()

	r.POST("/telegram/send", rpc.Wrap(notificationHandler.SendToTelegram))
}
