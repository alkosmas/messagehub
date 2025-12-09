package domain

type MessageType string

const (
	MessageTypeSMS   MessageType = "sms"
	MessageTypeEmail MessageType = "email"
)

type Message struct {
	Type    MessageType // sms ή email
	To      string      // +30699... ή email@example.com
	Subject string
	Body    string
}
