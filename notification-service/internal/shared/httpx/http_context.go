package httpx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HttpCtx http context (implements gin.Context and respond.Ctx)
type HttpCtx struct {
	*gin.Context
}

func Wrap(h func(*HttpCtx)) gin.HandlerFunc {
	return func(c *gin.Context) { h(&HttpCtx{c}) }
}

func (h *HttpCtx) RequestID() string {
	return h.GetString(CtxKeyRequestID)
}

func (h *HttpCtx) Endpoint() string {
	if ep := h.GetString(CtxKeyEndpoint); ep != "" {
		return ep
	}

	if ep := h.FullPath(); ep != "" {
		return ep
	}

	return h.Request.URL.Path
}

func (h *HttpCtx) Logger() *zap.Logger {
	return FromLogger(h.Context)
}
