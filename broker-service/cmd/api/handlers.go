package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// agreed upon json format that all our microservices will adhere to. doesnt matter what were sending from our various services
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

// format of the json in our auth service's 'Authenticate' method
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

// method that will be called when we send a request to "localhost:80/handle" (will be mapped to 8080 through docker)
// this will handle requests from all our microservices as a way of being a single point of entry, so its very important
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	// try to read the json into the RequestPayload 'mold'
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// take a different action based on what kind of json we receive and its content
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json that well send to the auth microservice by encoding the email/password json we receive ('a')
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// prepare service to send a post request to the /authenicate endpoint defined in the auth-service routes.go file
	// we will prepare the recently encoded jsonData with the email/password as a request body
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// we will actually send the request now and get the response from the auth service
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	// make sure we get the correct status code from the auth service
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable that we will read response's Body (that we get from the auth service) into
	var jsonFromService jsonResponse
	// we call the writeJSON method in the Authenticate method in auth-service's handlers.go file
	// this means that we should receive back a json object that is of the same 'mold' as jsonFromService

	// define a json decoder so that we can decode the response's body (that we get from the auth service)
	dec := json.NewDecoder(res.Body)
	// check for error when decoding the jsons's body (res.Body) into 'jsonFromService'. if the 'mold' isnt the same its an error
	err = dec.Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// check if the response json contains some Error value in it
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// after all these checks, we know that we have a valid login, so we send back the user a payload with good info
	var payload jsonResponse
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonFromService.Data // as defined in the auth-service's Authenticate function, this will be our User

	app.writeJSON(w, http.StatusAccepted, payload)
}
