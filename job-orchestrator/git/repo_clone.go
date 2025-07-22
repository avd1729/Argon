package git

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"job-orchestrator/pkg"
	"job-orchestrator/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneAndReadRunnerCI(payload pkg.WebhookPayload) pkg.Config {
	// Step 1: Create temp dir
	tempDir, err := os.MkdirTemp("", "repo-*")
	utils.FailOnError(err, "Failed to create temp directory")

	log.Printf("📦 Cloning repo %s into %s", payload.RepositoryUrl, tempDir)

	// Step 2: Clone the repo
	cmd := exec.Command("git", "clone", payload.RepositoryUrl, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	utils.FailOnError(err, "Failed to clone repo")

	// Step 3: Checkout to commit
	cmd = exec.Command("git", "checkout", payload.CommitId)
	cmd.Dir = tempDir
	err = cmd.Run()
	utils.FailOnError(err, "Failed to checkout commit")

	// Step 4: Read .runnerci.yml
	ymlPath := filepath.Join(tempDir, ".runnerci.yml")
	content, err := ioutil.ReadFile(ymlPath)
	utils.FailOnError(err, "Failed to read .runnerci.yml")

	fmt.Println("\n.runnerci.yml content:")
	fmt.Println(string(content))

	// Step 5: Parse YAML into Go struct
	var config pkg.Config
	err = yaml.Unmarshal(content, &config)
	utils.FailOnError(err, "Failed to parse .runnerci.yml")

	return config
}
