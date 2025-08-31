package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"notification-service-api/internal/notifications/delivery/rpc/dto"
	"notification-service-api/internal/notifications/domain/entity"
	"notification-service-api/internal/shared/queue"
	"notification-service-api/internal/shared/queue/notifications"
	"notification-service-api/pkg/utils"
	"time"
)

type TelegramPort interface {
	Send(ctx context.Context, payload []byte) error
}

type TelegramService struct {
	TG       TelegramPort
	rabbitMQ *utils.RabbitMQConnection
	logger   *zap.Logger
}

func NewTelegramService(t TelegramPort, rabbitMQ *utils.RabbitMQConnection) *TelegramService {
	return &TelegramService{
		TG:       t,
		rabbitMQ: rabbitMQ,
	}
}

func (s *TelegramService) WithLogger(logger *zap.Logger) *TelegramService {
	s.logger = logger
	return s
}

func (s *TelegramService) EnqueueTelegram(ctx context.Context, correlationID string, req dto.TelegramRequestSendParams) (uuid.UUID, error) {
	notificationID := uuid.New()

	s.logger.Info(fmt.Sprintf("Start sending tg notification to queue, ID: %s", notificationID.String()))

	parseMode := "markdown"
	if req.ParseMode != nil {
		parseMode = *req.ParseMode
	}

	tgEvent := entity.TelegramNotification{
		NotificationID: notificationID,
		CorrelationID:  correlationID,
		To:             req.To,
		Payload:        req.Message,
		ParseMode:      parseMode,
		CreatedAt:      time.Now(),
	}

	s.logger.Info(fmt.Sprintf("Telegram notification: %v, ID: %s", tgEvent, notificationID.String()))

	ch, err := s.rabbitMQ.Channel()
	if err != nil {
		s.logger.Error("failed to open channel", zap.Error(err))
		return uuid.Nil, err
	}
	defer ch.Close()

	eventBinary, err := msgpack.Marshal(tgEvent)
	if err != nil {
		s.logger.Error("failed to encode telegram notification", zap.Error(err))
		return uuid.Nil, err
	}

	err = queue.Publish(ch, notifications.RoutingTelegramSend, eventBinary, amqp.Table{}, &correlationID)
	if err != nil {
		s.logger.Error("failed to enqueue telegram notification", zap.Error(err))
		return uuid.Nil, err
	}

	s.logger.Info(fmt.Sprintf("Telegram notification sent successfully, ID: %s", notificationID.String()))

	return notificationID, nil
}
