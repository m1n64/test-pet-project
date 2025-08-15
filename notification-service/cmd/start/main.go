package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		r.Use(utils.LoggingContextMiddleware(dependencies.Logger))
		//v1Group := r.Group("/v1")

		r.GET("/ping", func(c *gin.Context) {
			logger := c.MustGet(utils.CtxKeyLogger).(*zap.Logger)
			logger.Info("Pong")

			c.JSON(200, gin.H{"pong": true, "request_id": c.GetString("request_id")})
		})

		r.Run(":8000")
	}()

	select {}
}
