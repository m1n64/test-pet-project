package queue

import (
	"github.com/streadway/amqp"
	"log"
	"notification-service-api/internal/shared/queue/notifications"
)

func InitTopology(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(
		notifications.ExchangeNotifications,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return err
	}

	bindings := []struct {
		queue string
		key   string
	}{
		{notifications.QueueEmail, notifications.RoutingEmailSend},
		{notifications.QueueSMS, notifications.RoutingSMSSend},
		{notifications.QueueTelegram, notifications.RoutingTelegramSend},
	}

	for _, b := range bindings {
		_, err := ch.QueueDeclare(
			b.queue,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,
		)
		if err != nil {
			return err
		}

		if err := ch.QueueBind(
			b.queue,
			b.key,
			notifications.ExchangeNotifications,
			false,
			nil,
		); err != nil {
			return err
		}

		log.Printf("Queue declared: %s (rk=%s)", b.queue, b.key)
	}

	return nil
}
