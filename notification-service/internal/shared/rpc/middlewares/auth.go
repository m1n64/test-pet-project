package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"notification-service-api/internal/shared/rpc/respond"
	"notification-service-api/pkg/utils"
)

func AuthMiddleware(configuration *utils.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !configuration.IsSecure {
			c.Next()
			return
		}

		token := c.GetHeader("X-API-KEY")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, returnUnauthorized())
			return
		}

		if configuration.MasterToken == token {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, returnUnauthorized())
		return
	}
}

func returnUnauthorized() respond.Response[any] {
	return respond.BuildFail(nil, respond.AuthError, "Unauthorized", "Invalid token", nil)
}
