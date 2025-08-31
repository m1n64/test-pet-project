package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"notification-service-api/internal/shared/rpc"
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
			zap.String(rpc.CtxKeyRequestID, rid),
			zap.String(rpc.CtxKeyMethod, c.Request.Method),
			zap.String(rpc.CtxKeyEndpoint, ep),
			zap.String(rpc.CtxKeyClientIP, clientIP),
			zap.String(rpc.CtxKeyUserAgent, ua),
			//zap.Any(rpc.CtxKeyHeaders, c.Request.Header),
		)

		c.Set(rpc.CtxKeyLogger, reqLogger)
		c.Set(rpc.CtxKeyRequestID, rid)
		c.Writer.Header().Set("X-Request-ID", rid)

		c.Next()
	}
}
