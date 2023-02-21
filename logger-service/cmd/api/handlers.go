package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celso-patiri/go-micro/helpers"
	"github.com/celso-patiri/go-micro/logger/data"
)

type JSONPayload struct {
	name string
	data string
}

var tools = helpers.Tools{}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	//read json into var
	var requestPayload JSONPayload
	_ = tools.ReadJSON(w, *r, &requestPayload)

	//insert data
	event := data.LogEntry{
		Name:      requestPayload.name,
		Data:      requestPayload.data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		tools.ErrorJSON(w, err)
	}

	res := helpers.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("logged %s", event.Name),
	}

	tools.WriteJSON(w, http.StatusAccepted, res)
}
