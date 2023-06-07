package main

import (
	"net/http"

	"github.com/jateen67/log-service/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// method that will be called when we send a post request to "localhost:80/log" (will be mapped to 8080 through docker)
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// this is the json that the request will get decoded/fitted into
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// now we insert the data

	// create an event we will log
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	// insert it into the mongo db using the method defined in data/models.go
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	// write some json out to the user
	app.writeJSON(w, http.StatusAccepted, res)
}
