package service

import (
	"encoding/json"
	"fmt"
	"loan-mgt/uploader/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPConnection struct {
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	NotifQueue amqp.Queue
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
		"uploader", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	qNotif, err := ch.QueueDeclare(
		"notification", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	return &AMQPConnection{
		Conn:       conn,
		Channel:    ch,
		Queue:      q,
		NotifQueue: qNotif,
	}, nil
}

func (a *AMQPConnection) Consume() (<-chan amqp.Delivery, error) {
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

func (a *AMQPConnection) SendRequest(channelName string, jsonData []byte) error {

	return a.Channel.Publish(
		"",          // exchange
		channelName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
}

func (a *AMQPConnection) SendNotificationRequest(token, filename, userId string) error {

	data := map[string]string{
		"token":    token,
		"filename": filename,
		"userId":   userId,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return a.SendRequest(a.NotifQueue.Name, jsonData)

}
