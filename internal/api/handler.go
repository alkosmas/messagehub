package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alkosmas/messagehub/internal/domain"
	"github.com/alkosmas/messagehub/internal/manager"
)

// Handler χειρίζεται τα HTTP requests
type Handler struct {
	manager *manager.Manager
}

// NewHandler δημιουργεί έναν νέο Handler
func NewHandler(mgr *manager.Manager) *Handler {
	return &Handler{
		manager: mgr,
	}
}

// SendMessageRequest - τι περιμένουμε από τον client
type SendMessageRequest struct {
	Type    string `json:"type"`    // "sms" ή "email"
	To      string `json:"to"`      // παραλήπτης
	Subject string `json:"subject"` // θέμα (για email)
	Body    string `json:"body"`    // το μήνυμα
}

// SendMessageResponse - τι επιστρέφουμε στον client
type SendMessageResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Provider string `json:"provider"`
}

// ErrorResponse - για errors
type ErrorResponse struct {
	Error string `json:"error"`
}

// SendMessage χειρίζεται το POST /api/messages
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Βήμα 1: Έλεγχος ότι είναι POST
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed. Use POST.")
		return
	}

	// Βήμα 2: Διάβασε το request body
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Βήμα 3: Validation
	if req.Type == "" || req.To == "" || req.Body == "" {
		h.writeError(w, http.StatusBadRequest, "Missing required fields: type, to, body")
		return
	}

	// Βήμα 4: Μετατροπή σε domain.Message
	msg := &domain.Message{
		Type:    domain.MessageType(req.Type),
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
	}

	// Βήμα 5: Στείλε μέσω του Manager
	log.Printf("📨 API: Received request to send %s to %s", req.Type, req.To)

	err := h.manager.Send(r.Context(), msg)
	if err != nil {
		log.Printf("❌ API: Failed to send: %v", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Βήμα 6: Επιτυχία!
	log.Printf("✅ API: Message sent successfully")
	h.writeJSON(w, http.StatusOK, SendMessageResponse{
		ID:       "msg-001", // Προσωρινό - θα το φτιάξουμε με UUID αργότερα
		Status:   "sent",
		Provider: string(msg.Type), // Προσωρινό
	})
}

// Health χειρίζεται το GET /api/health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed. Use GET.")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "MessageHub is running",
	})
}

// writeJSON γράφει JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError γράφει error response
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{Error: message})
}
