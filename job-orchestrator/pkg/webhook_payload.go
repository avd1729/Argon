package pkg

type WebhookPayload struct {
	RepositoryUrl string `json:"repositoryUrl"`
	RepoName      string `json:"repoName"`
	Branch        string `json:"branch"`
	CommitId      string `json:"commitId"`
}
