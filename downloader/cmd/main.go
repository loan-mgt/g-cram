package main

import (
	"fmt"

	"loan-mgt/uploader/internal/config"
	"loan-mgt/uploader/internal/handler"
	"loan-mgt/uploader/internal/service"
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
	msgs, err := amqpConn.Consume()
	if err != nil {
		fmt.Println(err)
		return
	}

	for d := range msgs {
		handler.HandleDownload(d.Body, amqpConn)
	}
}
