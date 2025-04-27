package service

import (
	"encoding/json"
	"fmt"
	"loan-mgt/cramer/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPConnection struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	DoneQueue amqp.Queue
	Queue     amqp.Queue
}

type DoneMsg struct {
	MediaId   string `json:"mediaId"`
	UserId    string `json:"userId"`
	Timestamp int64  `json:"timestamp"`
	FileSize  int64  `json:"fileSize"`
}

func NewAMQPConnection(cfg *config.Config) (*AMQPConnection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQUser, cfg.RabbitMQPass, cfg.RabbitMQHost, cfg.RabbitMQPort))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"cramer_done", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, err
	}

	cramerQueue, err := ch.QueueDeclare(
		"cramer", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	return &AMQPConnection{
		Conn:      conn,
		Channel:   ch,
		DoneQueue: q,
		Queue:     cramerQueue,
	}, nil
}

func (a *AMQPConnection) SendRequest(mediaId, userId string, timestamp, fileSize int64) error {
	data := DoneMsg{
		MediaId:   mediaId,
		UserId:    userId,
		Timestamp: timestamp,
		FileSize:  fileSize,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return a.Channel.Publish(
		"",               // exchange
		a.DoneQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
}

func (a *AMQPConnection) ConsumeCramerQueue() (<-chan amqp.Delivery, error) {
	return a.Channel.Consume(
		a.Queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
}
