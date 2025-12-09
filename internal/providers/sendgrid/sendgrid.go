package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/alkosmas/messagehub/internal/domain"
	"github.com/alkosmas/messagehub/internal/providers"
)

type SendGridProvider struct {
	apiKey    string
	fromEmail string
	client    *http.Client
	baseURL   string
}

type Config struct {
	APIKey    string
	FromEmail string
}

func New(cfg Config) *SendGridProvider {
	return &SendGridProvider{
		apiKey:    cfg.APIKey,
		fromEmail: cfg.FromEmail,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.sendgrid.com/v3",
	}
}

func (s *SendGridProvider) GetName() string {
	return "sendgrid"
}

func (s *SendGridProvider) GetType() domain.MessageType {
	return domain.MessageTypeEmail
}

type emailRequest struct {
	Personalizations []struct {
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
	} `json:"personalizations"`
	From struct {
		Email string `json:"email"`
	} `json:"from"`
	Subject string `json:"subject"`
	Content []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
}

func (s *SendGridProvider) Send(ctx context.Context, msg *domain.Message) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(msg.To) {
		return fmt.Errorf("%w: invalid email format", providers.ErrInvalidRecipient)
	}

	reqBody := emailRequest{
		Personalizations: []struct {
			To []struct {
				Email string `json:"email"`
			} `json:"to"`
		}{{To: []struct {
			Email string `json:"email"`
		}{{Email: msg.To}}}},
		From: struct {
			Email string `json:"email"`
		}{Email: s.fromEmail},
		Subject: msg.Subject,
		Content: []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{{Type: "text/plain", Value: msg.Body}},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/mail/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("📤 SendGrid: Sending email to %s...\n", msg.To)
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", providers.ErrProviderDown, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 202 {
		fmt.Println("✅ SendGrid: Email accepted for delivery!")
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("sendgrid error: status %d, body: %s", resp.StatusCode, string(body))
}
