package http

import (
	"github.com/gin-gonic/gin"
	"notification-service-api/internal/shared/httpx"
	"notification-service-api/pkg/di"
)

func InitNotificationRoutes(r *gin.Engine, dependencies *di.Dependencies) {
	notificationHandler := NewNotificationHandler()

	r.POST("/telegram/send", httpx.Wrap(notificationHandler.SendToTelegram))
}
