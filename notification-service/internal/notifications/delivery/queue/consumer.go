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
		Queue:           notifications.QueueTelegram,
		Workers:         5,
		Prefetch:        5,
		Args:            amqp.Table{},
		RetryBackoff:    30 * time.Second,
		RetryMax:        3,
		RetryRoutingKey: notifications.RoutingTelegramSendRetry,
		DLQRoutingKey:   notifications.RoutingTelegramSendDLQ,
	}, handler.Handle)
	if err != nil {
		dependencies.Logger.Error("failed to register telegram consumer", zap.Error(err))
		return
	}

	dependencies.Logger.Info("Telegram consumer registered")
}

func StartEmailConsumers(dependencies *di.Dependencies) {
	ctx := context.Background()

	handler := NewEmailHandler(dependencies.Logger, dependencies.EmailService)

	err := dependencies.RabbitMQ.Consume(ctx, utils.ConsumeOptions{
		Queue:           notifications.QueueEmail,
		Workers:         5,
		Prefetch:        5,
		Args:            amqp.Table{},
		RetryBackoff:    30 * time.Second,
		RetryMax:        3,
		RetryRoutingKey: notifications.RoutingEmailSendRetry,
		DLQRoutingKey:   notifications.RoutingEmailSendDLQ,
	}, handler.Handle)
	if err != nil {
		dependencies.Logger.Error("failed to register email consumer", zap.Error(err))
		return
	}

	dependencies.Logger.Info("Email consumer registered")
}
