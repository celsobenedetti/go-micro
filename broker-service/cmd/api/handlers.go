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
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	err := tools.ReadJSON(w, *r, &RequestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	default:
		tools.ErrorJSON(w, errors.New("Unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, reqPayload AuthPayload) {
	//create some json well send to the auth microservice
	jsonData, _ := json.MarshalIndent(reqPayload, "", "\t")

	//call service
	req, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
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
