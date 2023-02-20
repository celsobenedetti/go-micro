package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/celso-patiri/go-micro/helpers"
)

var tools = helpers.Tools{}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := tools.ReadJSON(w, *r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate the user agains the database
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

	payload := helpers.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func writeInvalidCredentials(w http.ResponseWriter) {
	tools.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
}
