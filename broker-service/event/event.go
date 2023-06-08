package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// when we call this function in consumer.go, we want it to return nil
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type of the exchange
		true,         // is this exchange durable?
		false,        // do you get rid of it when you are done with it?
		false,        // exchange just used internally?
		false,        // no wait?
		nil,          // any specific arguments?
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name of queue
		false, // is the queue durable?
		false, // do you get rid of it when you are done with it?
		true,  // is the queue exclusive?
		false, // no wait?
		nil,   // any specific arguments?
	)
}
