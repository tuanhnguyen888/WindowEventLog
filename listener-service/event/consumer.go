package event

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"time"
)

type Consumer struct {
	conn          *amqp.Connection
	redisClient   *redis.Client
	elasticClient *elastic.Client
}

type logEvent struct {
	Level        string `json:"level" `
	ProviderName string `json:"provider_name" `
	Msg          string `json:"msg"`
	Created      int64  `json:"created"`
}

func NewConsumer(
	conn *amqp.Connection,
	redisClient *redis.Client,
	elasticClient *elastic.Client,
) (Consumer, error) {
	consumer := Consumer{
		conn:          conn,
		redisClient:   redisClient,
		elasticClient: elasticClient,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (consumer *Consumer) Listen() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareQueue(ch)
	if err != nil {
		return err
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			event := logEvent{}
			_ = json.Unmarshal(d.Body, &event)
			log.Printf("Received a message: %#v", event)

			go consumer.handlePayload(event)
		}
	}()

	fmt.Println("[*] Waiting for messages. To exit press CTRL+C ")

	<-forever

	return nil
}

func (consumer *Consumer) handlePayload(payload logEvent) {
	dataJson, err := json.Marshal(payload)
	//js := string(dataJson)
	if err != nil {
		log.Println(err)
	}

	//ind, err := consumer.elasticClient.Index().
	//	Index("logs").
	//	BodyJson(js).
	//	Do(context.Background())
	//if err != nil {
	//	log.Println(err, ind)
	//}

	err = consumer.redisClient.Set(strconv.FormatInt(payload.Created, 10), dataJson, 30*24*time.Hour).Err()
	if err != nil {
		log.Println(err)
	}

}
