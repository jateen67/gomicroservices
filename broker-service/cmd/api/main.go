package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("starting broker service on port %s\n", port)

	// define http server with stuff like the port number and the routes we will use
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
