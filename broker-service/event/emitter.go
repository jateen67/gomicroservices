package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

// similar to the setup() function in the listener service consumer.go file
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

// push event to queue
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("pushing to channel...")

	err = channel.Publish(
		"logs_topic", // name of the exchange
		severity,     // either "log.INFO", "log.WARNING", or "log.ERROR"
		false,        // is mandatory?
		false,        // is immediate?
		amqp.Publishing{ // type amqp.Publishing
			ContentType: "text/plain",
			Body:        []byte(event), // payload of the message
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	// create new emitter
	emitter := Emitter{
		connection: conn,
	}

	// set it up
	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
