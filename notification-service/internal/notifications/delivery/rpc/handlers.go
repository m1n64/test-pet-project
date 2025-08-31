package rpc

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"notification-service-api/internal/notifications/app"
	"notification-service-api/internal/notifications/delivery/rpc/dto"
	"notification-service-api/internal/shared/rpc"
	"notification-service-api/internal/shared/rpc/respond"
)

type NotificationHandler struct {
	validator       *validator.Validate
	telegramService *app.TelegramService
}

func NewNotificationHandler(validator *validator.Validate, telegramService *app.TelegramService) *NotificationHandler {
	return &NotificationHandler{
		validator:       validator,
		telegramService: telegramService,
	}
}

func (h *NotificationHandler) SendToTelegram(c *rpc.HttpCtx, params dto.TelegramRequestSendParams) (any, *respond.RPCError) {
	if err := h.validator.Struct(params); err != nil {
		return nil, respond.NewRPCError(respond.InvalidParams, "invalid_params", "invalid params", err.Error())
	}

	h.telegramService.WithLogger(c.Logger())

	id, err := h.telegramService.EnqueueTelegram(c.Context, c.RequestID(), params)
	if err != nil {
		c.Logger().Error("enqueue_telegram", zap.Error(err))
		return nil, respond.NewRPCError(respond.InternalError, "enqueue_telegram", "enqueue_telegram", err.Error())
	}

	return dto.TelegramResponseSendDTO{NotificationID: id.String(), Queued: true}, nil
}
