package handlers

import (
	"bytes"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"net/http"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/shared/rpc/respond"
)

type RPCHandler struct {
	registry *rpc.Registry
}

func NewRPCHandler(registry *rpc.Registry) *RPCHandler {
	return &RPCHandler{
		registry: registry,
	}
}

func (h *RPCHandler) MainRPCHandler(c *rpc.HttpCtx) {
	raw, err := c.GetRawData()

	if err != nil {
		respond.Fail(c, nil, respond.ParseError, "parse_error", "failed to read body", err.Error())
		return
	}

	s := bytes.TrimSpace(raw)
	if len(s) == 0 {
		respond.Fail(c, nil, respond.InvalidRequest, "invalid_request", "empty body", nil)
		return
	}

	isBatch := false

	var requests []respond.Request
	switch s[0] {
	case '[':
		if err := json.Unmarshal(s, &requests); err != nil {
			respond.Fail(c, nil, respond.InvalidRequest, "invalid_request", "failed to read body", err.Error())
			return
		}

		if len(requests) == 0 {
			respond.Fail(c, nil, respond.InvalidRequest, "invalid_request", "empty body", nil)
			return
		}

		isBatch = true

		break
	case '{':
		var r respond.Request
		if err := json.Unmarshal(s, &r); err != nil {
			respond.Fail(c, nil, respond.InvalidRequest, "invalid_request", "failed to read body", err.Error())
			return
		}

		requests = []respond.Request{r}
		break
	default:
		respond.Fail(c, nil, respond.InvalidRequest, "invalid_request", "invalid request", nil)
		return
	}

	resps := make([]respond.Response[any], 0, len(requests))
	for _, req := range requests {
		if req.JSONRPC != respond.Version || req.Method == "" {
			if req.ID != nil {
				resps = append(resps, respond.BuildFail(req.ID, respond.InvalidRequest, "invalid_request", "invalid request", nil))
			}

			continue
		}

		result, rpcErr := h.registry.Call(req.Method, c, req.Params)

		if req.ID == nil {
			if rpcErr != nil {
				c.Logger().Warn("rpc notification error",
					zap.String("method", req.Method),
					zap.Any("error", rpcErr),
				)
			}
			continue
		}

		if rpcErr != nil {
			var appCode string
			var details any
			if rpcErr.Data != nil {
				appCode = rpcErr.Data.Code
				details = rpcErr.Data.Details
			}
			resps = append(resps, respond.BuildFail(req.ID, rpcErr.Code, appCode, rpcErr.Message, details))
			continue
		}

		resps = append(resps, respond.BuildOK(req.ID, result))
	}

	switch {
	case len(resps) == 0:
		c.Status(http.StatusNoContent)
	case len(requests) == 1 && isBatch:
		c.JSON(http.StatusOK, resps)
	case len(requests) == 1:
		c.JSON(http.StatusOK, resps[0])
	default:
		c.JSON(http.StatusOK, resps)
	}
}
