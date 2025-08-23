package http

import (
	"notification-service-api/internal/shared/httpx"
	"notification-service-api/internal/shared/httpx/respond"
	"notification-service-api/internal/system/delivery/http/dto"
)

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Ping(c *httpx.HttpCtx) {
	c.Logger().Info("ping")
	c.Logger().Warn("warn")
	c.Logger().Error("error")

	respond.OK(c, dto.PingDTO{Pong: true})
}
