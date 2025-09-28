package domain

type Channel string

const (
	ChannelTelegram Channel = "telegram"
	ChannelEmail    Channel = "email"
)

func (c Channel) String() string {
	return string(c)
}

type NotificationType string

const (
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeError   NotificationType = "error"
)

func (nt NotificationType) String() string {
	return string(nt)
}

type NotificationMonitoring interface {
	Send(channel Channel, notificationType NotificationType, delta int64)
	SendSuccess(channel Channel, count int64)
	SendError(channel Channel, count int64)
}
