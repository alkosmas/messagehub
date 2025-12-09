package domain

import "time"

type MessageType string

const (
    MessageTypeSMS   MessageType = "sms"
    MessageTypeEmail MessageType = "email"
)

type Message struct {
    ID        string      // Unique ID (UUID)
    Type      MessageType
    To        string
    Subject   string
    Body      string
    Status    string      // pending, sent, failed
    CreatedAt time.Time
}
