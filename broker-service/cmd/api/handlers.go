package main

import (
	"net/http"

	"github.com/celso-patiri/go-micro/helpers"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	tools := helpers.Tools{}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "Hello from Broker",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}
