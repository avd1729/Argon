package utils

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"os"
	"sandbox-executor/pkg"
)

func SendLog(level string, message string, context any) {
	type LogPayload struct {
		Level   string      `json:"level"`   // e.g., "INFO", "ERROR", "WARN"
		Message string      `json:"message"` // your log message
		Context interface{} `json:"context"` // optional: step, container, etc.
	}

	logPayload := LogPayload{
		Level:   level,
		Message: message,
		Context: context,
	}

	body, _ := json.Marshal(logPayload)

	err := godotenv.Load()
	FailOnError(err, "Error loading .env file")

	url := os.Getenv("RABBIT_MQ_LISTENER_URL")
	conn, err := amqp091.Dial(url)
	if err != nil {
		fmt.Println("Logger Error: failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Logger Error: failed to open channel:", err)
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"",
		string(pkg.LoggerQueue),
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		fmt.Println("Logger Error: failed to publish log:", err)
	}
}
