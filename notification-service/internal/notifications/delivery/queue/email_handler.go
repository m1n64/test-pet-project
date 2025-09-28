package queue

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"notification-service-api/internal/notifications/app"
	"notification-service-api/internal/notifications/domain/entity"
)

type EmailHandler struct {
	logger       *zap.Logger
	emailService *app.EmailService
}

func NewEmailHandler(logger *zap.Logger, emailService *app.EmailService) *EmailHandler {
	return &EmailHandler{
		logger:       logger,
		emailService: emailService,
	}
}

func (h *EmailHandler) Handle(ctx context.Context, d amqp.Delivery) error {
	logger := h.logger.With(zap.String("request_id", d.CorrelationId))
	h.emailService.WithLogger(logger)

	logger.Info("Handling email...")

	var email *entity.EmailNotification
	if err := msgpack.Unmarshal(d.Body, &email); err != nil {
		logger.Error("failed to unmarshal telegram message", zap.Error(err))
		return err
	}

	return h.emailService.SendEmail(ctx, email)
}
