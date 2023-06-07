package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// json type created that we will send to make sure that the broker service works
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// function to read json any time i want to
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// specify max size of uploaded json that is being read
	maxBytes := 1048576 // 1mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// define a json decoder so that we can decode the json's body
	dec := json.NewDecoder(r.Body)
	// check for error when decoding the jsons's body (r.Body) into 'data'
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// one more check to make: make sure there is only one single json value in the file being received
	err = dec.Decode(&struct{}{})
	// if there is only one single json value being received, it will throw an eof, since the json wont be too big
	// if it returns an error other than the eof, it means that there is indeed more than one json value being sent
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// function to write json any time i want to
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// encode the json we want to write so that we can send it in the correct format
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// check if any headers were included as parameters to this function
	if len(headers) > 0 {
		for key, value := range headers[0] {
			// add the header to the json we are writing
			w.Header()[key] = value
		}
	}

	// write the data out
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// function to write error json
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	// default status code if the status is not defined in the function parameter
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	// write the error using the writeJSON method defined above
	return app.writeJSON(w, statusCode, payload)
}
