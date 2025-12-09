package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alkosmas/messagehub/internal/domain"
	"github.com/alkosmas/messagehub/internal/providers"
	"github.com/alkosmas/messagehub/internal/storage"
	"github.com/google/uuid"
)

var ErrNoProvider = errors.New("no provider available for this message type")

type Manager struct {
	providers map[domain.MessageType]providers.Provider
	repo      *storage.Repository // Προσθήκη DB repository
}

// New δέχεται τώρα και το repository
func New(repo *storage.Repository) *Manager {
	return &Manager{
		providers: make(map[domain.MessageType]providers.Provider),
		repo:      repo,
	}
}

func (m *Manager) RegisterProvider(p providers.Provider) {
	m.providers[p.GetType()] = p
}

func (m *Manager) Send(ctx context.Context, msg *domain.Message) error {
	// 1. Προετοιμασία Message (ID, Timestamps)
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()
	msg.Status = "PENDING"

	// 2. Αποθήκευση στη βάση ως PENDING
	if m.repo != nil {
		fmt.Println("💾 Saving message to DB...")
		if err := m.repo.Save(msg); err != nil {
			fmt.Printf("⚠️ Failed to save to DB: %v\n", err)
		}
	}

	// 3. Εύρεση Provider
	provider, exists := m.providers[msg.Type]
	if !exists {
		// Καταγραφή αποτυχίας
		if m.repo != nil {
			m.repo.UpdateStatus(msg.ID, "FAILED", "", "No provider found")
		}
		return fmt.Errorf("%w: %s", ErrNoProvider, msg.Type)
	}

	// 4. Αποστολή
	err := provider.Send(ctx, msg)

	// 5. Ενημέρωση βάσης με το αποτέλεσμα
	if m.repo != nil {
		status := "SENT"
		errorMsg := ""
		if err != nil {
			status = "FAILED"
			errorMsg = err.Error()
		}
		m.repo.UpdateStatus(msg.ID, status, provider.GetName(), errorMsg)
	}

	return err
}

func (m *Manager) ListProviders() {
	for t, p := range m.providers {
		fmt.Printf("   %s → %s\n", t, p.GetName())
	}
}
