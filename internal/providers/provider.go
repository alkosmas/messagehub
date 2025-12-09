package providers

import (
	"context"
	"errors"

	"github.com/alkosmas/messagehub/internal/domain"
)

var (
	ErrInvalidRecipient = errors.New("invalid recipient")
	ErrRateLimited      = errors.New("rate limited by provider")
	ErrAuthFailed       = errors.New("authentication failed")
	ErrProviderDown     = errors.New("provider unavailable")
)

type Provider interface {
	Send(ctx context.Context, msg *domain.Message) error

	GetName() string

	GetType() domain.MessageType
}
