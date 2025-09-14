package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	rpc2 "notification-service-api/internal/notifications/delivery/rpc"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/shared/rpc/handlers"
	"notification-service-api/internal/shared/rpc/middlewares"
	systemRpc "notification-service-api/internal/system/delivery/rpc"
	"notification-service-api/pkg/di"
	"notification-service-api/pkg/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	utils.InitMigrations(dependencies.DB)

	r := gin.Default()

	r.Use(middlewares.LoggingContextMiddleware(dependencies.Logger))
	r.Use(middlewares.AccessLogMiddleware())
	r.Use(middlewares.AuthMiddleware(dependencies.Config, dependencies.Logger))
	r.Use(middlewares.StatisticsMiddleware(dependencies.Influx, dependencies.Logger))

	systemRpc.InitSystemProcedures(dependencies)
	rpc2.InitNotificationProcedures(dependencies)

	rpcHandler := handlers.NewRPCHandler(dependencies.Registry)

	r.POST("/rpc", rpc.Wrap(rpcHandler.MainRPCHandler))

	srv := &http.Server{Addr: ":8000", Handler: r}

	go func() {
		fmt.Println("Server started on port 8000")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("listen error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	if dependencies.DB != nil {
		sqlDB, _ := dependencies.DB.DB()
		_ = sqlDB.Close()
	}
	if dependencies.Logger != nil {
		_ = dependencies.Logger.Sync()
	}

	if dependencies.Influx != nil {
		_ = dependencies.Influx.Close()
	}

	fmt.Println("Server exiting")
}
