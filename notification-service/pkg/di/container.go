package di

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"notification-service-api/pkg/utils"
	"os"
)

type Dependencies struct {
	Logger    *zap.Logger
	Redis     *redis.Client
	DB        *gorm.DB
	Validator *validator.Validate
}

func InitDependencies() *Dependencies {
	// Infrastructure
	utils.LoadEnv()
	logger := utils.InitLogs()

	logger.Info("Init redis")
	redisConn := utils.CreateRedisConn(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	logger.Info("Init DB")
	dbConn := utils.InitDBConnection(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	logger.Info("Init migrations")
	utils.InitMigrations(dbConn)

	validate := utils.InitValidator()

	logger.Info("Init dependencies successfully")

	return &Dependencies{
		Logger:    logger,
		Redis:     redisConn,
		DB:        dbConn,
		Validator: validate,
	}
}
