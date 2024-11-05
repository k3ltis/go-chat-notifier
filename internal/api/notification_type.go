package api

type NotificationType int

const (
	StateInfo = iota
	StateWarning
)

var notificationStateName = map[NotificationType]string{
	StateInfo:    "info",
	StateWarning: "warning",
}
