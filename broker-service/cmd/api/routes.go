// specify all the routes in our application
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is able to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// easily make sure that the service is running by hitting the endpoint to get a response
	mux.Use(middleware.Heartbeat("/ping"))

	// add routes that use handlers, which will be called when we access these routes
	// post request to localhost:80 will run the Broker method (will be mapped to 8080 through docker)
	mux.Post("/", app.Broker)

	// grpc route just for ease of reference
	mux.Post("/log-grpc", app.LogItemViaGRPC)

	// a single point of entry that will handle all requests from all other microservices
	mux.Post("/handle", app.HandleSubmission)

	return mux
}
