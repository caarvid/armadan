package armadan

import (
	"context"
	"time"
)

type SessionService interface {
	GetByToken(context.Context, string) (*Session, error)
	Create(context.Context, string, bool) (*Session, error)
	DeleteByToken(context.Context, string) error
}

type Session struct {
	ID        int64
	UserID    string
	ExpiresAt time.Time
	Token     string
	Email     string
	Role      Role
}

func (s *Session) IsValid() bool {
	return s.ExpiresAt.After(time.Now())
}
