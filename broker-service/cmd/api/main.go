package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	webPort          = "8080"
	rabbitmqMaxTries = 5
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("broker-service: Connected to RabbitMQ")

	app := Config{
        Rabbit: rabbitConn,
    }

	log.Printf("Starting broker service at port %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection

	backoff := 1 * time.Second

	// dont continue until rabbimq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Printf("RabbitMQ not yet ready, trying again in: %s ...\n", backoff)
			counts++

		} else {
			connection = c
			return connection, nil
		}

		if counts > rabbitmqMaxTries {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backoff)
	}
}
