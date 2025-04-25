package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"encoding/json"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	reqDump, err := httputil.DumpRequest(r, true)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("\nAuthService received authenticate request, Payload:\n%s", string(reqDump))
	err = app.readJson(w, r, &requestPayload)
	if err != nil {
		log.Println("Auth payload read error:", err)
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

/*	//check DB connection by fetching all users
	allUsers, err := app.Models.User.GetAll()
	fmt.Printf("All users in table: %+v\n", allUsers)
	if err != nil {
		log.Println("\nError fetching all users from DB:", err)
		app.errorJson(w, errors.New("error fetching all users"), http.StatusBadRequest)
		return
	}
*/
	//validate user from DB
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	fmt.Printf("Check User in DB: %+v\n", user)
	if err != nil {
		log.Println("\nError finding user in DB:", err)
		app.errorJson(w, errors.New("error finding user"), http.StatusBadRequest)
		return
	}
	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	log.Println("Check if password matches:", valid)

	if err != nil{
		app.errorJson(w, errors.New("error matching password"), http.StatusBadRequest)
		return
	} else if !valid {
		app.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User Authenticated: %s", user.Email),
		Data:    user,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(("error calling log-service from auth-service"), err)
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request);
	if err != nil {
		return err
	}
	return nil
}