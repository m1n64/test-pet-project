package http

import (
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/system/delivery/http/dto"
)

func InitSystemRoutes(registry *rpc.Registry) {
	systemHandler := NewSystemHandler()

	registry.Register("ping", rpc.Typed[dto.PingParams](systemHandler.Ping))
}
