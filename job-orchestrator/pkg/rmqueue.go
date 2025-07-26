package pkg

type RMQueue string

const (
	WebhookQueue RMQueue = "webhook.queue"
	SandboxQueue RMQueue = "sandbox.queue"
	LoggerQueue  RMQueue = "log.queue"
)
