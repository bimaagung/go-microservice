package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ampq "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *ampq.Connection
	queueName string
}

func NewConsumer(conn *ampq.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	// setup consumer
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	// Open channel
	channel, err := consumer.conn.Channel()
	if err!= nil {
        return err
    }

	return declareExchange(channel)
}


// entity Payload
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	// open session in channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	// declare random queue
	q, err := declareRandomQueue(ch)
	if err!= nil {
        return err
    }

	// check queue binding
	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		
		if err != nil {
			return  err
		}
	}

	// get message from consume
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]", q.Name)
	<-forever

    return nil	
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Panicln(err)
		}
	case "auth":
		// authenticate
	// you can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Panicln(err)
		}
	}
}

func logEvent(entry Payload) error {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type","application/json")

	client := &http.Client{}
	
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}