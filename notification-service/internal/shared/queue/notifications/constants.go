package notifications

const (
	ExchangeNotifications = "notifications"

	ExchangeDLX   = "notifications.dlx"
	ExchangeRetry = "notifications.retry"

	RoutingEmailSend    = "email.send"
	RoutingSMSSend      = "sms.send"
	RoutingTelegramSend = "telegram.send"

	QueueEmail    = "notifications.email"
	QueueSMS      = "notifications.sms"
	QueueTelegram = "notifications.telegram"

	DeadQueueEmail    = "notifications.email.dlq"
	DeadQueueSMS      = "notifications.sms.dlq"
	DeadQueueTelegram = "notifications.telegram.dlq"
)
