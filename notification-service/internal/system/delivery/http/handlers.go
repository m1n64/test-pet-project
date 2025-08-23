package http

import (
	"github.com/google/uuid"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/shared/rpc/respond"
	"notification-service-api/internal/system/delivery/http/dto"
)

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Ping(c *rpc.HttpCtx) {
	c.Logger().Info("ping")
	c.Logger().Warn("warn")
	c.Logger().Error("error")

	respond.OK(c, uuid.NewString(), dto.PingDTO{Pong: true})
}
