package pkg

type Config struct {
	Version string         `yaml:"version" json:"version"`
	Jobs    map[string]Job `yaml:"jobs" json:"jobs"`
}
