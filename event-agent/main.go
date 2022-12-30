package main

import (
	"agent/event"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	winlog "github.com/ofcoursedude/gowinlog"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

type Connection struct {
	conn *amqp.Connection
}

type logEvent struct {
	Level        string `json:"level" `
	ProviderName string `json:"provider_name" `
	Msg          string `json:"msg"`
	Created      int64  `json:"created"`
}

func main() {

	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	var c = Connection{
		conn: rabbitConn,
	}
	forever := make(chan bool)
	go c.getWindowsEventLog()
	<-forever
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

func (b *Connection) pushToQueue(body logEvent, ctx context.Context) error {

	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(body)
	if err != nil {
		return err
	}

	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}

	q, err := event.DeclareQueue(ch)
	if err != nil {
		return err
	}
	//log.Printf("break point")

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        reqBodyBytes.Bytes(),
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %#v\n", body)

	return nil
}

func getWatcher() (*winlog.WinLogWatcher, error) {
	log.Println("Starting listen win log event ...")
	watcher, err := winlog.NewWinLogWatcher()
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

// getWindowsEventLog
func (b *Connection) getWindowsEventLog() {
	watcher, err := getWatcher()
	if err != nil {
		panic(err)
	}

	err = watcher.SubscribeFromNow("Application", "*")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case evt := <-watcher.Event():
			// Print the event struct
			// fmt.Printf("\nEvent: %v\n", evt)
			// or print basic output

			timeUnix := evt.Created.UnixMicro()
			if evt.LevelText == "" {
				evt.LevelText = "Information"
			}
			body := logEvent{
				Level:        evt.LevelText,
				ProviderName: evt.ProviderName,
				Msg:          evt.Msg,
				Created:      timeUnix,
			}

			//log.Printf(" [x] Sent %#v\n", body)

			err = b.pushToQueue(body, ctx)
			if err != nil {
				panic(err)
			}
		case err := <-watcher.Error():
			fmt.Printf("\nError: %v\n\n", err)
		default:
			// If no event is waiting, need to wait or do something else, otherwise
			// the app fails on deadlock.
			<-time.After(1 * time.Millisecond)
		}
	}
}
