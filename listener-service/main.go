package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jateen67/listener/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// then start listening for messages
	log.Println("listening for and consuming rabbitmq messages...")

	// create consumer to consume messages from the queue
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume the events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	// attempt to connect a fixed number of times
	var count int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq") // name specified in our docker-compose file
		if err != nil {
			fmt.Println("rabbitmq not yet ready...")
			count++
		} else {
			log.Println("connected to rabbitmq successfully")
			connection = c
			break
		}

		// if we didnt connect after 5 tries something is wrong
		if count > 5 {
			fmt.Println(err)
			return nil, err
		}

		// increase the delay each time i back off
		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
	}

	return connection, nil
}
