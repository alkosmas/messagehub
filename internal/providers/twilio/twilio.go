package twilio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alkosmas/messagehub/internal/domain"
	"github.com/alkosmas/messagehub/internal/providers"
)

// TwilioProvider στέλνει SMS μέσω Twilio API
type TwilioProvider struct {
	accountSID string       // Το Account SID από το Twilio Console
	authToken  string       // Το Auth Token από το Twilio Console
	fromNumber string       // Ο αριθμός που στέλνει (π.χ. +1234567890)
	client     *http.Client // HTTP client για τα requests
	baseURL    string       // Το URL του Twilio API
}

// Config - ρυθμίσεις για τον Twilio Provider
type Config struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// New δημιουργεί έναν νέο Twilio Provider
func New(cfg Config) *TwilioProvider {
	return &TwilioProvider{
		accountSID: cfg.AccountSID,
		authToken:  cfg.AuthToken,
		fromNumber: cfg.FromNumber,
		client: &http.Client{
			Timeout: 30 * time.Second, // Timeout μετά από 30 sec
		},
		baseURL: "https://api.twilio.com/2010-04-01",
	}
}

// GetName επιστρέφει το όνομα του provider
func (t *TwilioProvider) GetName() string {
	return "twilio"
}

// GetType επιστρέφει τον τύπο μηνυμάτων που χειρίζεται
func (t *TwilioProvider) GetType() domain.MessageType {
	return domain.MessageTypeSMS
}

// Send στέλνει ένα SMS μέσω Twilio
func (t *TwilioProvider) Send(ctx context.Context, msg *domain.Message) error {
	// Βήμα 1: Έλεγχος ότι ο αριθμός είναι σωστός
	if !strings.HasPrefix(msg.To, "+") {
		return fmt.Errorf("%w: phone must start with +", providers.ErrInvalidRecipient)
	}

	// Βήμα 2: Ετοιμασία των δεδομένων για το Twilio API
	data := url.Values{}
	data.Set("To", msg.To)         // Προς ποιον
	data.Set("From", t.fromNumber) // Από ποιον
	data.Set("Body", msg.Body)     // Το μήνυμα

	// Βήμα 3: Δημιουργία του HTTP request
	reqURL := fmt.Sprintf("%s/Accounts/%s/Messages.json", t.baseURL, t.accountSID)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Βήμα 4: Προσθήκη headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Basic Auth: username = accountSID, password = authToken
	req.SetBasicAuth(t.accountSID, t.authToken)

	// Βήμα 5: Αποστολή του request
	fmt.Printf("📤 Twilio: Sending SMS to %s...\n", msg.To)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", providers.ErrProviderDown, err)
	}
	defer resp.Body.Close()

	// Βήμα 6: Διάβασε την απάντηση
	body, _ := io.ReadAll(resp.Body)

	// Βήμα 7: Έλεγχος αν πέτυχε
	if resp.StatusCode == 201 {
		// Επιτυχία!
		fmt.Printf("✅ Twilio: SMS sent successfully!\n")
		return nil
	}

	// Αποτυχία - ας δούμε γιατί
	return t.handleError(resp.StatusCode, body)
}

// handleError διαχειρίζεται τα errors από το Twilio
func (t *TwilioProvider) handleError(statusCode int, body []byte) error {
	// Parse το error response από το Twilio
	var twilioErr struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(body, &twilioErr)

	fmt.Printf("❌ Twilio Error: %s (code: %d)\n", twilioErr.Message, twilioErr.Code)

	// Επιστροφή κατάλληλου error
	switch statusCode {
	case 401:
		return providers.ErrAuthFailed
	case 429:
		return providers.ErrRateLimited
	default:
		return fmt.Errorf("twilio error: %s", twilioErr.Message)
	}
}
