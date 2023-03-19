package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/celso-patiri/go-micro/broker/event"
	"github.com/celso-patiri/go-micro/broker/logs"
	"github.com/celso-patiri/go-micro/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tools = helpers.Tools{}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RPCPayload struct {
	Name string
	Data string
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
		// app.logEventViaRabbit(w, RequestPayload.Log)
		app.LogItemViaRPC(w, RequestPayload.Log)
	case "mail":
		app.sendMail(w, RequestPayload.Mail)
	default:
		tools.ErrorJSON(w, errors.New("Unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, reqPayload AuthPayload) {
	// create some json well send to the auth microservice
	jsonData, _ := json.MarshalIndent(reqPayload, "", "\t")

	// call service
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

	// make sure we get back the correct statuscode
	if res.StatusCode == http.StatusUnauthorized {
		tools.ErrorJSON(w, errors.New("Invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("Error calling auth service, not accepted"))
		return
	}

	// create a variable we'll read res.Body into
	var jsonFromService helpers.JSONResponse

	// decode the json from auth service
	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	// auth service returned StatusUnauthorized
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
	req, err := http.NewRequest("POST", loggerServiceUrl, bytes.NewBuffer(jsonData))
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
		Error:   false,
		Message: "logged",
	}

	tools.WriteJSON(w, http.StatusAccepted, resPayload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceUrl := "http://mail-service/send"

	// post to mail service
	req, err := http.NewRequest(http.MethodPost, mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	defer res.Body.Close()

	// make sure we get back the right status code
	if res.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("Error calling mail service"))
		return
	}

	// send back json response
	var payload helpers.JSONResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	var payload helpers.JSONResponse
	payload.Error = false
	payload.Message = "Logged via RabbitMQ"

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")

	if err != nil {
		return err
	}
	return nil
}

func (app *Config) LogItemViaRPC(w http.ResponseWriter, payload LogPayload) {
	client, err := rpc.Dial("tcp", loggerServerRPC)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: payload.Name,
		Data: payload.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	jsonPayload := helpers.JSONResponse{
		Error:   false,
		Message: result,
	}

	tools.WriteJSON(w, http.StatusAccepted, jsonPayload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(
		loggerServerGRPC,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	defer conn.Close()

	client := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	var payload helpers.JSONResponse
	payload.Error = false
	payload.Message = "Logged with gRPC"

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

const (
	authenticateUrl  = "http://authentication-service/authenticate"
	loggerServiceUrl = "http://logger-service/log"
	loggerServerRPC  = "logger-service:5001"
	loggerServerGRPC = "logger-service:50001"
)
