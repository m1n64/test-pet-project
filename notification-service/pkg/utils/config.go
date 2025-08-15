package utils

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		GetLogger().Sugar().Error("Error loading .env file")
	}
}

func IsLocal() bool {
	return os.Getenv("SERVICE_ENV") == "local"
}

func IsDev() bool {
	return os.Getenv("SERVICE_ENV") == "dev" || os.Getenv("SERVICE_ENV") == "development"
}

func IsProd() bool {
	return os.Getenv("SERVICE_ENV") == "prod" || os.Getenv("SERVICE_ENV") == "production"
}
