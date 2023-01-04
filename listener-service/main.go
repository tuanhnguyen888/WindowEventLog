package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/database"
	"listener/event"
	"log"
	"math"
	"os"
	"time"
)

type logEvent struct {
	Level        string `json:"level" `
	ProviderName string `json:"provider_name" `
	Msg          string `json:"msg"`
	Created      int64  `json:"created"`
}

func main() {
	// ----- try connect

	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	redisClient, err := database.NewInitRedis()
	if err != nil {
		log.Panic(err)
	}

	elasticClient, err := database.NewInitEsClient()
	if err != nil {
		log.Panic(err)
	}

	// ----- listening
	// ------ create consumer

	fmt.Println("Starting...")
	consumer, err := event.NewConsumer(conn, redisClient, elasticClient)
	if err != nil {
		panic(err)
	}

	// ------  watch the queue ..
	err = consumer.Listen()
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Panic("not yet ready ....")
			counts++
		} else {
			connection = conn
			break
		}

		if counts > 5 {
			fmt.Println("can not connect")
			return nil, err
		}

		backOff := time.Duration(math.Pow(float64(counts), 2)) * time.Second

		log.Println("backing off ...")
		time.Sleep(backOff)
	}

	return connection, nil
}
