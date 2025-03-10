package main

import (
	"fmt"
	"loan-mgt/cramer/internal/config"
	"loan-mgt/cramer/internal/handler"
	"loan-mgt/cramer/internal/service"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Create AMQP connection
	amqpConn, err := service.NewAMQPConnection(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer amqpConn.Conn.Close()

	// Consume messages from queue
	msgs, err := amqpConn.ConsumeCramerQueue()
	if err != nil {
		fmt.Println(err)
		return
	}

	for d := range msgs {
		handler.HandleCompression(d.Body, amqpConn)
	}

}
