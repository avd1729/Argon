package queue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/streadway/amqp"
)

type WebhookPayload struct {
	RepositoryUrl string `json:"repositoryUrl"`
	RepoName      string `json:"repoName"`
	Branch        string `json:"branch"`
	CommitId      string `json:"commitId"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func cloneAndReadRunnerCI(payload WebhookPayload) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "repo-*")
	failOnError(err, "Failed to create temp directory")

	log.Printf("Cloning repo %s into %s", payload.RepositoryUrl, tempDir)

	// Clone the repo
	cmd := exec.Command("git", "clone", payload.RepositoryUrl, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	failOnError(err, "Failed to clone repo")

	// Checkout the specific commit
	cmd = exec.Command("git", "checkout", payload.CommitId)
	cmd.Dir = tempDir
	err = cmd.Run()
	failOnError(err, "Failed to checkout commit")

	// Read the .runnerci.yml
	ymlPath := filepath.Join(tempDir, ".runnerci.yml")
	content, err := ioutil.ReadFile(ymlPath)
	failOnError(err, "Failed to read .runnerci.yml")

	fmt.Println("\n .runnerci.yml content:")
	fmt.Println(string(content))
}

func Listen() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"webhook.queue", true, false, false, false, nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println("\n📦 Received message:")
			var payload WebhookPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Println("Error unmarshaling payload:", err)
				continue
			}
			cloneAndReadRunnerCI(payload)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
