package main

import (
	"context"
	"fmt"

	"github.com/alkosmas/messagehub/internal/domain"
	"github.com/alkosmas/messagehub/internal/providers"
)

func main() {
	fmt.Println("🚀 MessageHub Starting...")

	smsProvider := providers.NewConsoleProvider("console-sms", domain.MessageTypeSMS)

	emailProvider := providers.NewConsoleProvider("console-email", domain.MessageTypeEmail)

	smsMessage := &domain.Message{
		Type: domain.MessageTypeSMS,
		To:   "+306947994665",
		Body: "Hello from MessageHub!",
	}

	fmt.Println("\n📱 Sending SMS...")
	err := smsProvider.Send(context.Background(), smsMessage)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}

	emailMessage := &domain.Message{
		Type:    domain.MessageTypeEmail,
		To:      "test@example.com",
		Subject: "Test Email",
		Body:    "This is a test email from MessageHub!",
	}

	fmt.Println("\n📧 Sending Email...")
	err = emailProvider.Send(context.Background(), emailMessage)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}

	fmt.Println("\n✅ Done!")
}
