package queue

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"job-orchestrator/git"
	"job-orchestrator/pkg"
	"job-orchestrator/utils"
	"log"
	"os"
)

func Listen() {

	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	url := os.Getenv("RABBIT_MQ_LISTENER_URL")

	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		string(pkg.WebhookQueue), "", true, false, false, false, nil,
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println("\n Received message:")
			var payload pkg.WebhookPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Println("Error unmarshaling payload:", err)
				continue
			}
			config, meta := git.CloneAndReadRunnerCI(payload)

			// For now, assume we only handle "build" job
			jobName := "build"
			job := config.Jobs[jobName]

			sandboxJob := pkg.SandboxPayload{
				RepoURL:  meta.RepositoryUrl,
				CommitID: meta.CommitId,
				JobName:  jobName,
				Job:      job,
			}

			SendToSandbox(sandboxJob)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
