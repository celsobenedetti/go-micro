package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/celso-patiri/go-micro/authentication/data"
)

const (
	webPort  = "80"
	maxConnectionTries = 10
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	conn := connectToDb()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	// set up config
	app := &Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDb() *sql.DB {
	dsn := os.Getenv("DSN")

	connectionTries := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready, ", err)
			connectionTries++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if connectionTries > maxConnectionTries {
			log.Println(err)
			return nil
		}

		wait := time.Duration(2)
		log.Printf("Waiting for %d seconds...\n", wait)
		time.Sleep(wait * time.Second)
	}
}
