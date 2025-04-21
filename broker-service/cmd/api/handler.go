package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// since the handler "HandleSubmission" will accept all requests coming into broker service,
// we need to define a struct that will hold the request payload
type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json: "password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.SendMail(w, requestPayload.Mail)
	default:
		app.errorJson(w, errors.New("unknown action"))
		
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "applicationm/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create payload to send to authenticate-service
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	// log.Println("Borker Service, sending auth request. Payload:", string(jsonData))
	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service:/authenticate", bytes.NewBuffer(jsonData))
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
	respDump, err := httputil.DumpResponse(response, true)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Broker Service, Auth RESPONSE:\n%s", string(respDump))

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling authentication service"))
		return
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

func (app *Config) SendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.Marshal(msg)

	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("\n broker service, send mail err:\n", err)
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJson(w, http.StatusAccepted, payload)
}