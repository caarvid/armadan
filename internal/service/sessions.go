package service

import (
	"context"
	"time"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toSession(data any) *armadan.Session {
	switch s := data.(type) {
	case schema.GetSessionByTokenRow:
		return &armadan.Session{
			ID:        s.ID,
			UserID:    s.UserID,
			ExpiresAt: s.ExpiresAt.Time,
			Active:    s.IsActive,
			Token:     s.Token,
			Email:     s.Email,
			Role:      armadan.Role(s.Role.UsersRoleEnum),
		}
	case schema.UserSession:
		return &armadan.Session{
			ID:        s.ID,
			UserID:    s.UserID,
			ExpiresAt: s.ExpiresAt.Time,
			Active:    s.IsActive,
			Token:     s.Token,
		}
	}

	return &armadan.Session{}
}

type sessions struct {
	db schema.Querier
}

func NewSessionService(db schema.Querier) *sessions {
	return &sessions{
		db: db,
	}
}

func (s *sessions) GetByToken(ctx context.Context, token string) (*armadan.Session, error) {
	session, err := s.db.GetSessionByToken(ctx, token)

	if err != nil {
		return nil, err
	}

	return toSession(session), nil
}

func getExpiration(keepLoggedIn bool) pgtype.Timestamptz {
	if keepLoggedIn {
		return pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true}
	}

	return pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true}
}

func (s *sessions) Create(ctx context.Context, id uuid.UUID, keepLoggedIn bool) (*armadan.Session, error) {
	session, err := s.db.CreateSession(ctx, &schema.CreateSessionParams{
		UserID:    id,
		IsActive:  true,
		Token:     uuid.NewString(),
		ExpiresAt: getExpiration(keepLoggedIn),
	})

	if err != nil {
		return nil, err
	}

	return toSession(session), nil
}

func (s *sessions) DeleteByToken(ctx context.Context, token string) error {
	return s.db.DeleteSession(ctx, token)
}
