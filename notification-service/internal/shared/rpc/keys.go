package rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CtxKeyLogger    = "logger"
	CtxKeyRequestID = "request_id"
	CtxKeyEndpoint  = "endpoint"
	CtxKeyMethod    = "method"
	CtxKeyClientIP  = "client_ip"
	CtxKeyUserAgent = "user_agent"
	CtxKeyHeaders   = "client_headers"
)

var baseLogger *zap.Logger

func SetBaseLogger(l *zap.Logger) { baseLogger = l }

func FromLogger(c *gin.Context) *zap.Logger {
	if v, ok := c.Get(CtxKeyLogger); ok {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}

	return zap.L()
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil || baseLogger == nil {
		return zap.L()
	}
	if v := ctx.Value(CtxKeyLogger); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return baseLogger
}

func FromGin(c *gin.Context) *zap.Logger {
	if c == nil {
		if baseLogger != nil {
			return baseLogger
		}
		return zap.L()
	}
	if v, ok := c.Get(CtxKeyLogger); ok {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	if baseLogger != nil {
		return baseLogger
	}
	return zap.L()
}
