package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// since the handler "HandleSubmission" will accept all requests coming into broker service,
// we need to define a struct that will hold the request payload
type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}
type AuthPayload struct {
	Email string `json:"email"`
	Password string `json: "password"`
}
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(w,r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJson(w, errors.New("unknown action"))
		
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create payload to send to authenticate-service
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", "http://localhost:8081/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w,err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	// This is to make sure that after we are done with processing the response.Body, we do not leave it open
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling authentication service"))
	}

	var jsonFromService jsonResponse
	// checking if response has same structure as defined in jsonFromService 
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	// checking if we got error in the reponse
	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data
	app.writeJson(w, http.StatusAccepted, payload) 
}