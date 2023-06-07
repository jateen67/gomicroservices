package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// method that will be called when we send a post request to "localhost:80/authenticate" (will be mapped to 8080 through docker)
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// this is the json that the request will get decoded/fitted into
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// check if the login attempt can be decoded into that requestPayload 'mold'
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// if the json is in the correct format and can be decoded, we now want to validate the email against the db
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// now we validate the password
	validPass, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !validPass {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication to logger-service
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// helper function that logs to the logger-service anytime we try to authenticate
func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	// create json that well send to the log microservice by encoding the name/data json we receive ('entry')
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// setup the request to the logger service
	req, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// we will actually send the request now and get the response from the auth service
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
