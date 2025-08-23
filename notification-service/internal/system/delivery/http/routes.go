package http

import (
	"github.com/gin-gonic/gin"
	"notification-service-api/internal/shared/httpx"
)

func InitSystemRoutes(r *gin.Engine) {
	systemHandler := NewSystemHandler()

	r.GET("/ping", httpx.Wrap(systemHandler.Ping))
}
