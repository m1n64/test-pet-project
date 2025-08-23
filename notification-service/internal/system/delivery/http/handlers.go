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

func (h *SystemHandler) PingHandler(c *httpx.HttpCtx) {
	c.Logger().Info("ping")

	respond.OK(c, dto.PingDTO{Pong: true})
}
