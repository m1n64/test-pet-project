package rpc

import (
	"fmt"
	"github.com/goccy/go-json"
	"notification-service-api/internal/shared/rpc/respond"
	"sync"
)

type Handler func(c *HttpCtx, params any) (any, *respond.RPCError)

type Registry struct {
	mu sync.RWMutex
	m  map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		m: make(map[string]Handler),
	}
}

func (r *Registry) Register(name string, h Handler) {
	r.mu.Lock()
	r.m[name] = h
	r.mu.Unlock() // Defer is overhead for this simple operation
}

func (r *Registry) Call(name string, c *HttpCtx, params any) (any, *respond.RPCError) {
	r.mu.RLock()
	h := r.m[name]
	r.mu.RUnlock()

	if h == nil {
		return nil, &respond.RPCError{Code: respond.MethodNotFound, Message: "method not found"}
	}

	return h(c, params)
}

func Typed[T any](fn func(c *HttpCtx, p T) (any, *respond.RPCError)) Handler {
	return func(c *HttpCtx, params any) (any, *respond.RPCError) {
		var p T

		switch v := params.(type) {
		case json.RawMessage:
			if len(v) > 0 && string(v) != "null" {
				if err := json.Unmarshal(v, &p); err != nil {
					return nil, &respond.RPCError{
						Code:    respond.InvalidParams,
						Message: "failed to bind params",
						Data:    &respond.ErrorDTO{Code: "invalid_params", Message: err.Error()},
					}
				}
			}
		case T:
			p = v
		case *T:
			if v != nil {
				p = *v
			}
		case nil:

		default:
			return nil, &respond.RPCError{
				Code:    respond.InvalidParams,
				Message: fmt.Sprintf("unsupported params type: %T", params),
				Data:    &respond.ErrorDTO{Code: "invalid_params"},
			}
		}

		return fn(c, p)
	}
}
