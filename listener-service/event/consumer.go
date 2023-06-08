package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// type used for receiving events from the queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// type used for pushing events to the queue
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	// declare consumer
	consumer := Consumer{
		conn: conn,
	}

	// set up the consumer by opening up a channel and declaring an exchange
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// helper function for setting up the channel
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	// function defined in event.go used to declare our exchange
	return declareExchange(channel)
}

// listens to the queue for specific topics
func (consumer *Consumer) Listen(topics []string) error {
	// go to our consumer channel and get things from it
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// we have our channel, now we need to get a random queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// we have our channel and our queue now

	// go through our list of topics
	for _, s := range topics {
		// bind our channel to each of these topics
		ch.QueueBind(
			q.Name,       // name of queue
			s,            // topic
			"logs_topic", // name of the exchange
			false,        // no wait?
			nil,          // any specific arguments?
		)

		// after we try to bind each topic, we check for error
		if err != nil {
			return err
		}
	}

	// look for messages
	messages, err := ch.Consume(
		q.Name, // name of queue
		"",     // name of consumer
		true,   // auto acknowledge?
		false,  // exclusive?
		false,  // internal>
		false,  // no wait?
		nil,    // any specific arguments
	)
	if err != nil {
		return err
	}

	// i want to do this forever; i want to consume all the things that come from rabbitmq until i exit the app
	// do this by declaring a new channel
	forever := make(chan bool)
	// will run in background
	go func() {
		for d := range messages {
			// decode d.Body into this json payload variable
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message [exchange, queue] [logs_topic, %s]\n", q.Name)
	// keep the consumption going forever by making this blocking
	<-forever

	return nil
}

// take an action based on the name of an event that we get pushed to us from the queue
func handlePayload(payload Payload) {
	// switch on the 'Name' value from the payload variable we received as a call to this function
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// logic to log an event to the logger service once we get it from rabbitmq
func logEvent(entry Payload) error {

	// create json that well send to the logger microservice by encoding the name/data json we receive ('entry')
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// prepare service to send a post request to the /log endpoint defined in the logger-service routes.go file
	// we will prepare the recently encoded jsonData with the name/data as a request body
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	// we will actually send the request now and get the response from the logger service
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// make sure we get the correct status code from the logger service
	if res.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
