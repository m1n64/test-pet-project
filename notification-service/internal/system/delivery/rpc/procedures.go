package rpc

import (
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/system/delivery/rpc/dto"
	"notification-service-api/pkg/di"
)

func InitSystemProcedures(dependencies *di.Dependencies) {
	systemHandler := NewSystemHandler()

	dependencies.Registry.Register("system.ping", rpc.Typed[dto.PingParams](systemHandler.Ping))
}
