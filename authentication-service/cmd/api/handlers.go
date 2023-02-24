package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/celso-patiri/go-micro/helpers"
)

var tools = helpers.Tools{}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user agains the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		writeInvalidCredentials(w)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		writeInvalidCredentials(w)
		return
	}

	// log authentication event
	err = app.logRequest("authentication", fmt.Sprintf("%.v - %s logged in", time.Now(), user.Email))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	req, err := http.NewRequest(http.MethodPost, logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil
	}

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func writeInvalidCredentials(w http.ResponseWriter) {
	tools.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
}

const (
	logServiceUrl = "http://logger-service/log"
)
