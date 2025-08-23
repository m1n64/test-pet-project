package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"notification-service-api/internal/shared/httpx"
	"strings"
)

func LoggingContextMiddleware(base *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(c.GetHeader("X-Request-ID"))
		if rid == "" {
			rid = uuid.NewString()
		}

		ep := c.FullPath()
		if ep == "" {
			ep = c.Request.URL.Path
		}

		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()

		reqLogger := base.With(
			zap.String(httpx.CtxKeyRequestID, rid),
			zap.String(httpx.CtxKeyMethod, c.Request.Method),
			zap.String(httpx.CtxKeyEndpoint, ep),
			zap.String(httpx.CtxKeyClientIP, clientIP),
			zap.String(httpx.CtxKeyUserAgent, ua),
		)

		c.Set(httpx.CtxKeyLogger, reqLogger)
		c.Set(httpx.CtxKeyRequestID, rid)
		c.Writer.Header().Set("X-Request-ID", rid)

		c.Next()
	}
}
