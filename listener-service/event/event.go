package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"logEvent",
		false,
		false,
		false,
		false,
		nil,
	)
}
