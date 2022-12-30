package event

import amqp "github.com/rabbitmq/amqp091-go"

//func DeclareExchange(ch *amqp.Channel) error {
//	return ch.ExchangeDeclare(
//		"logs",
//		"topic",
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//}

func DeclareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"logEvent",
		false,
		false,
		false,
		false,
		nil,
	)
}
