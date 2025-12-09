package providers

import (
	"context"
	"fmt"

	"github.com/alkosmas/messagehub/internal/domain"
)

type ConsoleProvider struct {
	name        string
	messageType domain.MessageType
}

func NewConsoleProvider(name string, msgType domain.MessageType) *ConsoleProvider {
	return &ConsoleProvider{
		name:        name,
		messageType: msgType,
	}
}

func (c *ConsoleProvider) Send(ctx context.Context, msg *domain.Message) error {
	fmt.Println("========================================")
	fmt.Printf("📤 CONSOLE PROVIDER: %s\n", c.name)
	fmt.Printf("   Type: %s\n", msg.Type)
	fmt.Printf("   To:   %s\n", msg.To)
	if msg.Subject != "" {
		fmt.Printf("   Subject: %s\n", msg.Subject)
	}
	fmt.Printf("   Body: %s\n", msg.Body)
	fmt.Println("========================================")
	return nil
}

func (c *ConsoleProvider) GetName() string {
	return c.name
}

func (c *ConsoleProvider) GetType() domain.MessageType {
	return c.messageType
}
