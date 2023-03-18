package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type
		true,         // durable?
		false,        // autoDeleted?
		false,        // internal?
		false,        // noWait?
		nil,          // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"rabbitQueue", // name
		false,         // durable
		false,         // delete when unused?
		true,          // exclusive
		false,         // noWait?
		nil,           // arguments
	)
}
