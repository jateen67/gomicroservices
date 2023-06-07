package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jateen67/authentication/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	// connect to db using our helper function defined below
	conn := connectToDB()
	if conn == nil {
		log.Panic("cant connect to postgres")
	}

	// set up some configuration from models.go
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	log.Printf("starting auth service on port %s\n", port)

	// define http server with stuff like the port number and the routes we will use
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// try to connect to db
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// were connected, so lets test
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// our docker postgres environment might not be ready before we try to connect to the db, so we need this function
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	// create infinite loop and stay in there until we connect to our database successfully
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready. retrying... ")
			count++
		} else {
			log.Println("connected to postgres")
			return conn
		}

		// try 20 times before failing
		if count > 20 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
