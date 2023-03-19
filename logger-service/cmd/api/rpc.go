package main

import (
	"context"
	"log"
	"time"

	"github.com/celso-patiri/go-micro/logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (s *RPCServer) LogInfo(payload RPCPayload, response *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error writing to mongo", err)
		return err
	}

	*response = "Processed payload via RPC" + payload.Name
	return nil
}
