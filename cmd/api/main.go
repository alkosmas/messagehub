package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alkosmas/messagehub/internal/api"
	"github.com/alkosmas/messagehub/internal/domain"
	"github.com/alkosmas/messagehub/internal/manager"
	"github.com/alkosmas/messagehub/internal/providers"
	"github.com/alkosmas/messagehub/internal/providers/sendgrid"
	"github.com/alkosmas/messagehub/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 MessageHub Starting...")

	// 1. Env
	godotenv.Load()

	// 2. Database Connection
	// user:password@localhost:5432/dbname?sslmode=disable
	connStr := "postgres://myadmin:secret123@localhost:5433/messagehub?sslmode=disable"
	fmt.Println("🔌 Connecting with:", connStr)

	repo, err := storage.NewPostgres(connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")

	// 3. Manager με DB
	mgr := manager.New(repo)

	// 4. Providers
	// Console SMS
	smsProvider := providers.NewConsoleProvider("console-sms", domain.MessageTypeSMS)
	mgr.RegisterProvider(smsProvider)

	// SendGrid Email
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	sendgridEmail := os.Getenv("SENDGRID_FROM_EMAIL")

	if sendgridKey != "" {
		sgProvider := sendgrid.New(sendgrid.Config{
			APIKey:    sendgridKey,
			FromEmail: sendgridEmail,
		})
		mgr.RegisterProvider(sgProvider)
		fmt.Println("✅ SendGrid provider registered")
	}

	// 5. API & Server
	handler := api.NewHandler(mgr)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", handler.SendMessage)
	mux.HandleFunc("/api/health", handler.Health)
	mux.Handle("/", http.FileServer(http.Dir("./web/static")))

	port := "8080"
	fmt.Printf("\n🌐 Server: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
