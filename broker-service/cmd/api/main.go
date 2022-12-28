package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

const port = "80"

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	var brokerStruct = brokerHandler{
		conn: rabbitConn,
	}
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	app.Use(CORSMiddleware())

	log.Printf("Starting broker service on port 80\n")

	app.POST("/", broker)

	go brokerStruct.getWindowsEventLog()
	//app.POST("/handleLog", handleLog)

	err = app.Run(":" + port)
	if err != nil {
		panic("can not run port ")
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
