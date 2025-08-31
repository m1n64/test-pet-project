package rpc

import (
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/shared/rpc/respond"
	"notification-service-api/internal/system/delivery/rpc/dto"
)

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Ping(c *rpc.HttpCtx, params dto.PingParams) (any, *respond.RPCError) {
	c.Logger().Info("ping")
	c.Logger().Warn("warn")
	c.Logger().Error("error")

	return dto.PingDTO{Pong: true, Timestamp: params.Timestamp}, nil
}
