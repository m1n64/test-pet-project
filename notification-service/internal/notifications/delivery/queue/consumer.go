package queue

import (
	"context"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"notification-service-api/internal/shared/queue/notifications"
	"notification-service-api/pkg/di"
	"notification-service-api/pkg/utils"
	"time"
)

func StartTelegramConsumers(dependencies *di.Dependencies) {
	ctx := context.Background()

	handler := NewTelegramHandler(dependencies.Logger, dependencies.TelegramService)

	err := dependencies.RabbitMQ.Consume(ctx, utils.ConsumeOptions{
		Queue:        notifications.QueueTelegram,
		Workers:      5,
		Prefetch:     5,
		Args:         amqp.Table{},
		RetryBackoff: 30 * time.Second,
		RetryMax:     3,
	}, handler.Handle)
	if err != nil {
		dependencies.Logger.Error("failed to register telegram consumer", zap.Error(err))
		return
	}

	dependencies.Logger.Info("Telegram consumer registered")
}
