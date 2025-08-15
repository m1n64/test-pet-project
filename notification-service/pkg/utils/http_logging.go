package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
)

const (
	CtxKeyLogger    = "logger"
	ctxKeyRequestID = "request_id"
	ctxKeyEndpoint  = "endpoint"
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

		reqLogger := base.With(
			zap.String("request_id", rid),
			zap.String("endpoint", ep),
		)

		c.Set(CtxKeyLogger, reqLogger)
		c.Set(ctxKeyRequestID, rid)
		c.Set(ctxKeyEndpoint, ep)
		c.Writer.Header().Set("X-Request-ID", rid)

		c.Next()
	}
}
