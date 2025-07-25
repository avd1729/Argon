package pkg

type SandboxPayload struct {
	RepoURL  string `json:"repo_url"`
	CommitID string `json:"commit_id"`
	JobName  string `json:"job_name"`
	Job      Job    `json:"job"`
}
