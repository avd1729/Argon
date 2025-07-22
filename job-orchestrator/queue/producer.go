package queue

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"job-orchestrator/pkg"
	"job-orchestrator/utils"
	"log"
	"os"
)

func SendToSandbox(config pkg.Config) {
	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	url := os.Getenv("RABBIT_MQ_LISTENER_URL")
	conn, err := amqp.Dial(url)

	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"sandbox.queue", // New queue
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare queue")

	// Convert config struct to JSON
	body, err := json.Marshal(config)
	utils.FailOnError(err, "Failed to marshal config")

	err = ch.Publish(
		"",     // Default exchange
		q.Name, // Queue name
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	utils.FailOnError(err, "Failed to publish message")

	log.Println("Sent job to sandbox.queue")
}
