package email

import (
	"context"
	"notification-service-api/internal/notifications/domain/entity"
	"notification-service-api/pkg/utils"
)

type MailAPI struct {
	smtpClient *utils.SMTPClient
}

func NewEmailAPI(smtpClient *utils.SMTPClient) *MailAPI {
	return &MailAPI{
		smtpClient: smtpClient,
	}
}

func (e *MailAPI) SendEmailViaSMTP(ctx context.Context, message *entity.EmailNotification) error {
	attachments := make([]utils.MailAttachment, len(message.Attachments))
	for _, attachment := range message.Attachments {
		attachments = append(attachments, utils.MailAttachment{
			Filename:    attachment.Filename,
			Data:        attachment.Data,
			ContentType: attachment.ContentType,
		})
	}

	msg := &utils.MailMessage{
		To:          message.To,
		Subject:     message.Subject,
		Body:        message.Body,
		ContentType: message.ContentType,
		ReplyTo:     message.ReplyTo,
		CC:          message.CC,
		BCC:         message.BCC,
		Attachments: attachments,
		From:        message.From,
	}

	return e.smtpClient.Send(msg)
}
