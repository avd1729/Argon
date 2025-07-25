package sandbox

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sandbox-executor/pkg"
	"strings"
	"time"
)

func RunJobInDocker(payload pkg.SandboxPayload) error {
	job := payload.Job
	repoPath := "/tmp/" + randomString(8)
	containerName := "sandbox_" + randomString(8)

	// Step 0: Clone the repo
	fmt.Println("Cloning repo:", payload.RepoURL)
	cloneCmd := exec.Command("git", "clone", payload.RepoURL, repoPath)
	if err := runCmd(cloneCmd); err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	// Step 1: Pull the image
	fmt.Println("Pulling image:", job.Image)
	if err := runCmd(exec.Command("docker", "pull", job.Image)); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	// Step 2: Create container
	cmd := exec.Command("docker", "create", "--name", containerName, "-it", job.Image, "sh")
	if err := runCmd(cmd); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Step 3: Copy repo into container
	copyCmd := exec.Command("docker", "cp", repoPath+"/.", containerName+":/app")
	if err := runCmd(copyCmd); err != nil {
		return fmt.Errorf("failed to copy repo into container: %w", err)
	}

	// Step 4: Start container
	if err := runCmd(exec.Command("docker", "start", containerName)); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Step 5: Execute each step in container (inside /app)
	for _, step := range job.Steps {
		fmt.Printf("Running step: %s\n", step.Name)
		execCmd := []string{"exec", containerName, "sh", "-c", fmt.Sprintf("cd /app && %s", step.Run)}
		if err := runCmd(exec.Command("docker", execCmd...)); err != nil {
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}
	}

	// Step 6: Cleanup
	defer exec.Command("docker", "rm", "-f", containerName).Run()
	defer exec.Command("rm", "-rf", repoPath).Run()

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

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	var result strings.Builder
	for i := 0; i < n; i++ {
		result.WriteByte(letters[rand.Intn(len(letters))])
	}

	return result.String()
}
