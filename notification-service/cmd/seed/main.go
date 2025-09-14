package main

import (
	"notification-service-api/pkg/di"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	dependencies.Logger.Info("Successfully seeded")

	defer dependencies.Redis.Close()
	defer dependencies.RabbitMQ.Close()
}
