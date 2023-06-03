package main

import (
	"net/http"
)

// method that will be called when we send a post request to "localhost:80/" (will be mapped to 8080 through docker)
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	// jsonResponse comes from helpers.go
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}

	// write the data out using the writeJSON method defined in helpers.go
	_ = app.writeJSON(w, http.StatusOK, payload)
}
