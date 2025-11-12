package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	durable = iota
	transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return nil, amqp.Queue{}, err
	}

	switch queueType {
	case durable:
		queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		err = ch.QueueBind(queueName, key, exchange, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}

		return ch, queue, nil
	case transient:
		queue, err := ch.QueueDeclare(queueName, false, true, true, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		err = ch.QueueBind(queueName, key, exchange, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}

		return ch, queue, nil

	default:
		return ch, amqp.Queue{}, nil
	}

}
