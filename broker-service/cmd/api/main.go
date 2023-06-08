package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("starting broker service on port %s\n", port)

	// define http server with stuff like the port number and the routes we will use
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
