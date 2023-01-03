package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"

	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err!= nil {
        log.Println(err)
		os.Exit(1)
    }

	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ!")

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err!= nil {
        log.Println(err)
    }
}

// connect to rabbitmq
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