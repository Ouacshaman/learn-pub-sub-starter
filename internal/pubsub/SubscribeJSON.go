package pubsub

import (
	"encoding/json"
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
	err := Subscribe(conn, exchange, queueName, key, queueType, handler, jsonUnmarshalHelp)
	if err != nil {
		return err
	}
	return nil
}

func jsonUnmarshalHelp[T any](data []byte) (T, error) {
	var out T
	err := json.Unmarshal(data, &out)
	return out, err
}
