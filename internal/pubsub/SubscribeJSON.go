package pubsub

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	delivery_ch, err := ch.Consume("", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for v := range delivery_ch {
			var gen_body T
			if err := json.Unmarshal(v.Body, &gen_body); err != nil {
				fmt.Println(err)
				continue
			}
			err = PublishJSON(ch, exchange, key, gen_body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			v.Ack(false)
		}
	}()

	if err != nil {
		return err
	}

	return nil
}
