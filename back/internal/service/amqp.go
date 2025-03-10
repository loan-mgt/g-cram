package service

import (
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
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
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil
}

func (a *AMQPConnection) SendRequest(token, id, creationDate, filename, videoPath string) error {
	data := map[string]string{
		"token":        token,
		"id":           id,
		"creationDate": creationDate,
		"name":         filename,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return a.Channel.Publish(
		"",           // exchange
		a.Queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
}
