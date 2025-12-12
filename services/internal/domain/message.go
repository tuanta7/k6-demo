package domain

type MessageChannel string

const (
	EmailChannel MessageChannel = "email"
	PushChannel  MessageChannel = "push"
)

type Message struct {
	Channel MessageChannel `json:"channel"`
}
