package pubsub

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	delivery_ch, err := ch.Consume(queueName, "", false, false, false, false, nil)
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
			ack_out := handler(gen_body)

			switch ack_out {
			case Ack:
				v.Ack(false)
				fmt.Println("Ack")
			case NackRequeue:
				v.Nack(false, true)
				fmt.Println("NackRequeue")
			case NackDiscard:
				v.Nack(false, false)
				fmt.Println("NackDiscard")
			default:
				continue
			}
		}
	}()

	return nil
}
