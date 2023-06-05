package main

import (
	"errors"
	"fmt"
	"net/http"
)

// method that will be called when we send a post request to "localhost:80/authenticate" (will be mapped to 8081 through docker)
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

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
