package pkg

type RMQueue string

const (
	NotificationQueue RMQueue = "notification.queue"
	SandboxQueue      RMQueue = "sandbox.queue"
	LoggerQueue       RMQueue = "log.queue"
)
