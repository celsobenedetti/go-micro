package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/celso-patiri/go-micro/helpers"
)

var tools = helpers.Tools{}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	tools := helpers.Tools{}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "Hello from Broker",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var RequestPayload RequestPayload

	err := tools.ReadJSON(w, r, &RequestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	case "log":
		app.logItem(w, RequestPayload.Log)
	default:
		tools.ErrorJSON(w, errors.New("Unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, reqPayload AuthPayload) {
	//create some json well send to the auth microservice
	jsonData, _ := json.MarshalIndent(reqPayload, "", "\t")

	//call service
	req, err := http.NewRequest(http.MethodPost, authenticateUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	defer res.Body.Close()

	//make sure we get back the correct statuscode
	if res.StatusCode == http.StatusUnauthorized {
		tools.ErrorJSON(w, errors.New("Invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("Error calling auth service, not accepted"))
		return
	}

	//create a variable we'll read res.Body into
	var jsonFromService helpers.JSONResponse

	//decode the json from auth service
	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	//auth service returned StatusUnauthorized
	if jsonFromService.Error {
		tools.ErrorJSON(w, err, http.StatusUnauthorized)
	}

	resPayload := helpers.JSONResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	tools.WriteJSON(w, http.StatusAccepted, resPayload)
}

func (app *Config) logItem(w http.ResponseWriter, reqPayload LogPayload) {
	jsonData, _ := json.MarshalIndent(reqPayload, "", "\t")

	// call the service
    req, err := http.NewRequest("POST", logUrl, bytes.NewBuffer(jsonData))
    if err != nil {
        tools.ErrorJSON(w, err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}

    res, err := client.Do(req)
    if err != nil {
        tools.ErrorJSON(w, err)
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusAccepted {
        tools.ErrorJSON(w, err)
    }

    resPayload := helpers.JSONResponse{
        Error: false,
        Message: "logged",
    }

    tools.WriteJSON(w, http.StatusAccepted, resPayload)
}

const (
	authenticateUrl = "http://authentication-service/authenticate"
	logUrl          = "http://logger-service/log"
)
