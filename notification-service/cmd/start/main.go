package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"notification-service-api/internal/notifications/delivery/queue"
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
	if err := dependencies.MultiCache.WarmFromRedis(); err != nil {
		dependencies.Logger.Error("failed to warm multi cache from redis", zap.Error(err))
	}
}

func main() {
	utils.InitMigrations(dependencies.DB)
	stopAutoFlush := dependencies.MultiCache.StartAutoFlush(5 * time.Minute)

	r := gin.Default()

	publicGroup := r.Group("")

	publicGroup.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File("html/docs.html")
	})

	rpcGroup := r.Group("")

	rpcGroup.Use(middlewares.LoggingContextMiddleware(dependencies.Logger))
	rpcGroup.Use(middlewares.AccessLogMiddleware())
	rpcGroup.Use(middlewares.AuthMiddleware(dependencies.Config, dependencies.Logger))
	rpcGroup.Use(middlewares.StatisticsMiddleware(dependencies.Influx, dependencies.Logger))

	systemRpc.InitSystemProcedures(dependencies)
	rpc2.InitNotificationProcedures(dependencies)

	rpcHandler := handlers.NewRPCHandler(dependencies.Registry)

	rpcGroup.POST("/rpc", rpc.Wrap(rpcHandler.MainRPCHandler))

	srv := &http.Server{Addr: ":8000", Handler: r}

	go queue.StartTelegramConsumers(dependencies)
	go queue.StartEmailConsumers(dependencies)

	srv.RegisterOnShutdown(func() {
		// flush multi cache to redis
		// in case of server crash, multi cache will be lost
		// so we need to flush it to redis before server shutdown
		//closeMultiCache()
	})

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

	closeMultiCache()
	stopAutoFlush()

	if dependencies.DB != nil {
		sqlDB, _ := dependencies.DB.DB()
		if err := sqlDB.Close(); err != nil {
			dependencies.Logger.Error("failed to close database", zap.Error(err))
		}
	}
	if dependencies.Logger != nil {
		if err := dependencies.Logger.Sync(); err != nil {
			dependencies.Logger.Error("failed to sync logger", zap.Error(err))
		}
	}

	if dependencies.Influx != nil {
		if err := dependencies.Influx.Close(); err != nil {
			dependencies.Logger.Error("failed to close influx", zap.Error(err))
		}
	}

	fmt.Println("Server exiting")
}

func closeMultiCache() {
	if dependencies.MultiCache != nil {
		if err := dependencies.MultiCache.FlushToRedisOnce(); err != nil {
			dependencies.Logger.Error("failed to flush multi cache to redis", zap.Error(err))
		}
	}
}
