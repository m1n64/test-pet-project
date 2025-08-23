package rpc

import (
	"github.com/goccy/go-json"
	"notification-service-api/internal/shared/rpc/respond"
	"sync"
)

type Handler func(c *HttpCtx, params json.RawMessage) (any, *respond.RPCError)

type Registry struct {
	mu sync.RWMutex
	m  map[string]Handler
}
