package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct{
	Rabbit *ampq.Connection
}

func main() {

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err!= nil {
        log.Println(err)
		os.Exit(1)
    }

	defer rabbitConn.Close()
	
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)


	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}
}

func connect() (*ampq.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *ampq.Connection

	// don't continue until rabbit is ready
	for {
		c, err := ampq.Dial("amqp://guest:guest@rabbitmq")
		if err!= nil {
            fmt.Println("RabbitMQ not yest ready...")
			counts++
        } else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off ...")
		time.Sleep(backoff)

		continue
	}

	return connection, nil

}