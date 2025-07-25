package queue

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"os"
	"sandbox-executor/utils"
)

func Listen() {
	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	url := os.Getenv("RABBIT_MQ_LISTENER_URL")

	conn, err := amqp091.Dial(url)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"sandbox.queue",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever // block forever
}
