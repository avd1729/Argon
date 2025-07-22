package git

import (
	"fmt"
	"io/ioutil"
	"job-orchestrator/pkg"
	"job-orchestrator/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneAndReadRunnerCI(payload pkg.WebhookPayload) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "repo-*")
	utils.FailOnError(err, "Failed to create temp directory")

	log.Printf("Cloning repo %s into %s", payload.RepositoryUrl, tempDir)

	// Clone the repo
	cmd := exec.Command("git", "clone", payload.RepositoryUrl, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	utils.FailOnError(err, "Failed to clone repo")

	// Checkout the specific commit
	cmd = exec.Command("git", "checkout", payload.CommitId)
	cmd.Dir = tempDir
	err = cmd.Run()
	utils.FailOnError(err, "Failed to checkout commit")

	// Read the .runnerci.yml
	ymlPath := filepath.Join(tempDir, ".runnerci.yml")
	content, err := ioutil.ReadFile(ymlPath)
	utils.FailOnError(err, "Failed to read .runnerci.yml")

	fmt.Println("\n .runnerci.yml content:")
	fmt.Println(string(content))
}
