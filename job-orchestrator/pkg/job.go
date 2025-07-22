package pkg

type Job struct {
	Image string `yaml:"image" json:"image"`
	Steps []Step `yaml:"steps" json:"steps"`
}
