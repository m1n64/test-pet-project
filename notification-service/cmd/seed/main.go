package main

import (
	"golang-service-template/pkg/di"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	dependencies.Logger.Info("Successfully seeded")
}
