package queue

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"notification-service-api/internal/notifications/app"
	"notification-service-api/internal/notifications/domain/entity"
)

type TelegramHandler struct {
	logger          *zap.Logger
	TelegramService *app.TelegramService
}

func NewTelegramHandler(logger *zap.Logger, telegramService *app.TelegramService) *TelegramHandler {
	return &TelegramHandler{
		logger:          logger,
		TelegramService: telegramService,
	}
}

func (h *TelegramHandler) Handle(ctx context.Context, d amqp.Delivery) error {
	logger := h.logger.With(zap.String("request_id", d.CorrelationId))
	h.TelegramService.WithLogger(logger)

	logger.Info("Handling telegram message...")

	var notification *entity.TelegramNotification
	if err := msgpack.Unmarshal(d.Body, &notification); err != nil {
		logger.Error("failed to unmarshal telegram message", zap.Error(err))
		return err
	}

	return h.TelegramService.SendNotification(ctx, notification)
}
