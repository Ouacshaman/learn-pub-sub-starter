package pubsub

import (
	"fmt"
	"encoding/json"
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
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil{
		return err
	}

	delivery_ch, err := amqp.channel.Consume(nil, "", false, false, false, false, nil)
	if err != nil{
		return err
	}

	go func ()  {
		for v := range delivery_ch{
			var gen_body T
			if err := json.Unmarshal(v, &gen_body); err != nil{
				return err
			}
			err = PublishJSON(ch, exchange, key, gen_body)
			if err != nil{
				return err
			}
			delivery_ch.Ack(false)
		}
	}

	return nil
}
