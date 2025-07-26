package queue

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"os"
	"sandbox-executor/pkg"
	"sandbox-executor/sandbox"
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

	msgs, err := ch.Consume(
		string(pkg.SandboxQueue),
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
			var payload pkg.SandboxPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				fmt.Println("Invalid payload:", err)
				continue
			}
			fmt.Println("Running job:", payload.JobName)
			err := sandbox.RunJobInDocker(payload)
			if err != nil {
				fmt.Println("Job failed:", err)
			} else {
				fmt.Println("Job finished successfully")
			}
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever // block forever
}
