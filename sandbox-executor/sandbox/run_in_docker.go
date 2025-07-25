package sandbox

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sandbox-executor/pkg"
	"sandbox-executor/utils"
	"strings"
	"time"
)

func RunJobInDocker(payload pkg.SandboxPayload) error {
	job := payload.Job
	repoPath := "/tmp/" + randomString(8)
	containerName := "sandbox_" + randomString(8)

	notify := func(msg string) {
		sendNotification(fmt.Sprintf("Job '%s': %s", job, msg))
	}

	// Step 0: Clone the repo
	fmt.Println("Cloning repo:", payload.RepoURL)
	cloneCmd := exec.Command("git", "clone", payload.RepoURL, repoPath)
	if err := runCmd(cloneCmd); err != nil {
		notify("Failed to clone repo")
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	// Step 1: Pull the image
	fmt.Println("Pulling image:", job.Image)
	if err := runCmd(exec.Command("docker", "pull", job.Image)); err != nil {
		notify("Failed to pull image")
		return fmt.Errorf("failed to pull image: %w", err)
	}

	// Step 2: Create container
	cmd := exec.Command("docker", "create", "--name", containerName, "-it", job.Image, "sh")
	if err := runCmd(cmd); err != nil {
		notify("Failed to create container")
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Step 3: Copy repo into container
	copyCmd := exec.Command("docker", "cp", repoPath+"/.", containerName+":/app")
	if err := runCmd(copyCmd); err != nil {
		notify("Failed to copy repo into container")
		return fmt.Errorf("failed to copy repo into container: %w", err)
	}

	// Step 4: Start container
	if err := runCmd(exec.Command("docker", "start", containerName)); err != nil {
		notify("Failed to start container")
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Step 5: Execute steps
	for _, step := range job.Steps {
		fmt.Printf("Running step: %s\n", step.Name)
		execCmd := []string{"exec", containerName, "sh", "-c", fmt.Sprintf("cd /app && %s", step.Run)}
		if err := runCmd(exec.Command("docker", execCmd...)); err != nil {
			notify(fmt.Sprintf("Step '%s' failed", step.Name))
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}
	}

	defer exec.Command("docker", "rm", "-f", containerName).Run()
	defer exec.Command("rm", "-rf", repoPath).Run()

	notify("Job completed successfully")
	fmt.Println("Job completed successfully.")
	return nil
}

func runCmd(cmd *exec.Cmd) error {
	fmt.Println("Executing:", strings.Join(cmd.Args, " "))

	// Stream command output to terminal in real time
	cmd.Stdout = io.MultiWriter(os.Stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr)

	return cmd.Run()
}

func sendNotification(message string) {
	type NotifyPayload struct {
		Message string `json:"message"`
	}
	body, _ := json.Marshal(NotifyPayload{Message: message})

	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	url := os.Getenv("RABBIT_MQ_LISTENER_URL")
	conn, err := amqp091.Dial(url)

	if err != nil {
		fmt.Println("Notification Error: failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Notification Error: failed to open channel:", err)
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"",                   // exchange
		"notification.queue", // routing key
		false,                // mandatory
		false,                // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		fmt.Println("Notification Error: failed to publish message:", err)
	}
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	var result strings.Builder
	for i := 0; i < n; i++ {
		result.WriteByte(letters[rand.Intn(len(letters))])
	}

	return result.String()
}
