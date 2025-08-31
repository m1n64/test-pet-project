package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/shared/rpc/handlers"
	"notification-service-api/internal/shared/rpc/middlewares"
	"notification-service-api/internal/system/delivery/http"
	"notification-service-api/pkg/di"
	"notification-service-api/pkg/utils"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	utils.InitMigrations(dependencies.DB)

	go func() {
		fmt.Println("Server started on port 8000")

		r := gin.Default()
		r.Use(middlewares.LoggingContextMiddleware(dependencies.Logger))
		r.Use(middlewares.AccessLogMiddleware())

		http.InitSystemRoutes(dependencies.Registry)

		rpcHandler := handlers.NewRPCHandler(dependencies.Registry)

		r.POST("/rpc", rpc.Wrap(rpcHandler.MainRPCHandler))

		r.Run(":8000")
	}()

	select {}
}
