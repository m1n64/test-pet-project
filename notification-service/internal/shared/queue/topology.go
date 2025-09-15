package queue

import (
	"github.com/streadway/amqp"
	"log"
	"notification-service-api/internal/shared/queue/notifications"
	"time"
)

var defaultRetryTTL = 5 * time.Second

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

	if err := ch.ExchangeDeclare(notifications.ExchangeRetry, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(notifications.ExchangeDLX, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	bindings := []struct {
		queue string
		key   string
		dlq   string
	}{
		{notifications.QueueEmail, notifications.RoutingEmailSend, notifications.DeadQueueEmail},
		{notifications.QueueSMS, notifications.RoutingSMSSend, notifications.DeadQueueSMS},
		{notifications.QueueTelegram, notifications.RoutingTelegramSend, notifications.DeadQueueTelegram},
	}

	for _, b := range bindings {
		if err := declareWithRetryAndDLQ(ch, b.queue, b.key, b.dlq, defaultRetryTTL); err != nil {
			return err
		}
		log.Printf("Queue declared: %s (rk=%s)", b.queue, b.key)
	}

	return nil
}

func declareWithRetryAndDLQ(ch *amqp.Channel, mainQueue, routingKey, dlqQueue string, retryTTL time.Duration) error {
	retryRouting := routingKey + ".retry"
	dlqRouting := routingKey + ".dlq"
	retryQueue := mainQueue + ".retry"

	mainArgs := amqp.Table{
		"x-dead-letter-exchange":    notifications.ExchangeRetry,
		"x-dead-letter-routing-key": retryRouting,
	}
	if _, err := ch.QueueDeclare(mainQueue, true, false, false, false, mainArgs); err != nil {
		return err
	}
	if err := ch.QueueBind(mainQueue, routingKey, notifications.ExchangeNotifications, false, nil); err != nil {
		return err
	}

	retryArgs := amqp.Table{
		"x-message-ttl":             int32(retryTTL / time.Millisecond),
		"x-dead-letter-exchange":    notifications.ExchangeNotifications,
		"x-dead-letter-routing-key": routingKey,
	}
	if _, err := ch.QueueDeclare(retryQueue, true, false, false, false, retryArgs); err != nil {
		return err
	}
	if err := ch.QueueBind(retryQueue, retryRouting, notifications.ExchangeRetry, false, nil); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(dlqQueue, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(dlqQueue, dlqRouting, notifications.ExchangeDLX, false, nil); err != nil {
		return err
	}

	return nil
}
