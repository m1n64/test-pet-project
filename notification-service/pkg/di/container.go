package di

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"notification-service-api/internal/notifications/app"
	"notification-service-api/internal/notifications/infra/telegram"
	"notification-service-api/internal/shared/queue"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/pkg/utils"
	"os"
)

type Dependencies struct {
	Logger          *zap.Logger
	Redis           *redis.Client
	DB              *gorm.DB
	RabbitMQ        *utils.RabbitMQConnection
	Validator       *validator.Validate
	Registry        *rpc.Registry
	TelegramService *app.TelegramService
	Config          *utils.Config
	Influx          *utils.InfluxDB
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

	logger.Info("Init RabbitMQ")
	rabbitmqConn := utils.ConnectRabbitMQ(os.Getenv("RABBITMQ_URL"), logger)

	ch, err := rabbitmqConn.Channel()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to open channel: %v", err))
	}
	defer ch.Close()

	logger.Info("Init queues")
	if err := queue.InitTopology(ch); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to open channel: %v", err))
	}

	logger.Info("Init configuration")
	config := utils.LoadConfig()

	logger.Info("Init InfluxDB")
	influx, err := utils.InitInfluxUDP(os.Getenv("INFLUX_UDP_HOST"))
	if err != nil {
		logger.Fatal("Error InfluxDB connection")
	}

	validate := utils.InitValidator()

	registry := rpc.NewRegistry()

	tgApi := telegram.NewTGApiClient()
	tgService := app.NewTelegramService(tgApi, rabbitmqConn)

	logger.Info("Init dependencies successfully")

	return &Dependencies{
		Logger:          logger,
		Redis:           redisConn,
		DB:              dbConn,
		RabbitMQ:        rabbitmqConn,
		Validator:       validate,
		Registry:        registry,
		TelegramService: tgService,
		Config:          config,
		Influx:          influx,
	}
}
