package queue

import (
	"context"
	"github.com/streadway/amqp"
	"log"
)

type HandlerFunc func(ctx context.Context, d amqp.Delivery) error

func Consume(ctx context.Context, ch *amqp.Channel, queue string, workers int, handler HandlerFunc) error {
	if err := ch.Qos(workers, 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue,
		"",    // consumer tag
		false, // auto-ack = false
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case d, ok := <-msgs:
					if !ok {
						return
					}

					if err := handler(ctx, d); err != nil {
						log.Printf("Handler error: %v, nack msg id=%s", err, d.MessageId)
						_ = d.Nack(false, true)
						continue
					}

					_ = d.Ack(false)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return nil
}
