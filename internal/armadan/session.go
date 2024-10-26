package armadan

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionService interface {
	GetByToken(context.Context, string) (*Session, error)
	Create(context.Context, uuid.UUID, bool) (*Session, error)
	DeleteByToken(context.Context, string) error
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	Active    bool
	Token     string
	Email     string
	Role      Role
}

func (s *Session) IsValid() bool {
	return s.Active && s.ExpiresAt.After(time.Now())
}
