package queue

import (
	"github.com/streadway/amqp"
	"notification-service-api/internal/shared/queue/notifications"
	"time"
)

// Publish deprecated
func Publish(ch *amqp.Channel, routingKey string, body []byte, headers amqp.Table, correlationID *string) error {
	return ch.Publish(
		notifications.ExchangeNotifications,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/x-msgpack",
			Body:          body,
			Headers:       headers,
			Timestamp:     time.Now(),
			CorrelationId: *correlationID,
		},
	)
}
