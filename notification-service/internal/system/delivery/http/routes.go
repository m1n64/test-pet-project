package http

import (
	"github.com/gin-gonic/gin"
	"notification-service-api/internal/shared/rpc"
)

func InitSystemRoutes(r *gin.Engine) {
	systemHandler := NewSystemHandler()

	r.GET("/ping", rpc.Wrap(systemHandler.Ping))
}
