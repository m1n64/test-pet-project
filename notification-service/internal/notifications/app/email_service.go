package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"notification-service-api/internal/notifications/delivery/rpc/dto"
	"notification-service-api/internal/notifications/domain"
	"notification-service-api/internal/notifications/domain/entity"
	"notification-service-api/internal/shared/queue/notifications"
	"notification-service-api/pkg/utils"
	"time"
)

type EmailPort interface {
	SendEmailViaSMTP(ctx context.Context, message *entity.EmailNotification) error
}

type EmailService struct {
	emailAPI   EmailPort
	rabbitMQ   *utils.RabbitMQConnection
	logger     *zap.Logger
	monitoring domain.NotificationMonitoring
}

func NewEmailService(emailAPI EmailPort, rabbitMQ *utils.RabbitMQConnection, monitoring domain.NotificationMonitoring) *EmailService {
	return &EmailService{
		emailAPI:   emailAPI,
		rabbitMQ:   rabbitMQ,
		monitoring: monitoring,
	}
}

func (s *EmailService) WithLogger(logger *zap.Logger) *EmailService {
	s.logger = logger
	return s
}

func (s *EmailService) EnqueueEmail(ctx context.Context, correlationID string, req dto.EmailRequestSendParams) (uuid.UUID, error) {
	notificationID := uuid.New()

	s.logger.Info(fmt.Sprintf("Start sending email to queue, ID: %s", notificationID.String()))

	attachments := make([]entity.EmailAttachment, len(req.Attachments))
	for _, attachment := range req.Attachments {
		data, err := attachment.Bytes()
		if err != nil {
			s.logger.Error("failed to decode attachment", zap.Error(err))
			return uuid.Nil, err
		}

		attachments = append(attachments, entity.EmailAttachment{
			Filename:    attachment.Filename,
			Data:        data,
			ContentType: attachment.ContentType,
		})
	}

	emailEvent := entity.EmailNotification{
		NotificationID: notificationID,
		CorrelationID:  correlationID,
		To:             req.To,
		Subject:        req.Subject,
		Body:           req.Body,
		ContentType:    req.ContentType,
		CC:             req.CC,
		BCC:            req.BCC,
		From:           req.From,
		ReplyTo:        req.ReplyTo,
		Attachments:    attachments,
		CreatedAt:      time.Now(),
	}

	s.logger.Info(fmt.Sprintf("Email: %v, ID: %s", emailEvent, notificationID.String()))

	eventBinary, err := msgpack.Marshal(emailEvent)
	if err != nil {
		s.logger.Error("failed to encode email", zap.Error(err))
		return uuid.Nil, err
	}

	err = s.rabbitMQ.PublishMsgpack(ctx, notifications.ExchangeNotifications, notifications.RoutingEmailSend, eventBinary, amqp.Table{}, &correlationID)
	if err != nil {
		s.logger.Error("failed to enqueue email", zap.Error(err))
		return uuid.Nil, err
	}

	s.logger.Info(fmt.Sprintf("Email sent successfully, ID: %s", notificationID.String()))

	return notificationID, nil
}

func (s *EmailService) SendEmail(ctx context.Context, email *entity.EmailNotification) error {
	s.logger.Info(fmt.Sprintf("Sending email, ID: %s", email.NotificationID.String()))
	err := s.emailAPI.SendEmailViaSMTP(ctx, email)
	if err != nil {
		s.monitoring.SendError(domain.ChannelEmail, 1)
		s.logger.Error("failed to send email", zap.Error(err))
		return err
	}

	s.monitoring.SendSuccess(domain.ChannelEmail, 1)
	s.logger.Info(fmt.Sprintf("Email sent successfully, ID: %s", email.NotificationID.String()))
	return nil
}
