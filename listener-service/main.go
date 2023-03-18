package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/celso-patiri/go-micro/listener/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connet to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("listener-service: Connected to RabbitMQ")

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
        log.Println(err)
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

		if counts > max_counts {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backoff)
	}
}

const (
	max_counts = 5
)
