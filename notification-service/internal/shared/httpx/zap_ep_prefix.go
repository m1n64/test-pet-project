package httpx

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

type epPrefixCore struct {
	zapcore.Core
	method    string
	endpoint  string
	requestID string
}

func WrapWithEndpointPrefix(core zapcore.Core) zapcore.Core {
	return &epPrefixCore{Core: core}
}

func (c *epPrefixCore) With(fields []zapcore.Field) zapcore.Core {
	nc := &epPrefixCore{
		Core:      c.Core.With(fields),
		method:    c.method,
		endpoint:  c.endpoint,
		requestID: c.requestID,
	}
	for i := range fields {
		if fields[i].Key == CtxKeyEndpoint && fields[i].Type == zapcore.StringType {
			nc.endpoint = fields[i].String
		}

		if fields[i].Key == CtxKeyMethod && fields[i].Type == zapcore.StringType {
			nc.method = fields[i].String
		}

		if fields[i].Key == CtxKeyRequestID && fields[i].Type == zapcore.StringType {
			nc.requestID = fields[i].String
		}
	}
	return nc
}

func (c *epPrefixCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

func (c *epPrefixCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	if !strings.HasPrefix(e.Message, "[") {
		parts := make([]string, 0, 3)
		if c.method != "" {
			parts = append(parts, c.method)
		}
		if c.endpoint != "" {
			parts = append(parts, c.endpoint)
		}
		if c.requestID != "" {
			parts = append(parts, c.requestID)
		}
		if len(parts) > 0 {
			e.Message = "[" + strings.Join(parts, "] [") + "] " + e.Message
		}
	}

	return c.Core.Write(e, fields)
}

func (c *epPrefixCore) Sync() error {
	return c.Core.Sync()
}
