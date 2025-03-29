package service

import (
	"context"
	"time"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/google/uuid"
)

func toSession(data any) *armadan.Session {
	switch s := data.(type) {
	case schema.GetSessionByTokenRow:
		return &armadan.Session{
			ID:        s.ID,
			UserID:    s.UserID,
			ExpiresAt: armadan.ParseTime(s.ExpiresAt),
			Token:     s.Token,
			Email:     s.Email,
			Role:      armadan.Role(s.UserRole),
		}
	case schema.Session:
		return &armadan.Session{
			ID:        s.ID,
			UserID:    s.UserID,
			ExpiresAt: armadan.ParseTime(s.ExpiresAt),
			Token:     s.Token,
		}
	}

	return &armadan.Session{}
}

type sessions struct {
	dbReader schema.Querier
	dbWriter schema.Querier
}

func NewSessionService(reader, writer schema.Querier) *sessions {
	return &sessions{
		dbReader: reader,
		dbWriter: writer,
	}
}

func (s *sessions) GetByToken(ctx context.Context, token string) (*armadan.Session, error) {
	session, err := s.dbReader.GetSessionByToken(ctx, token)

	if err != nil {
		return nil, err
	}

	return toSession(session), nil
}

func getExpiration(keepLoggedIn bool) string {
	if keepLoggedIn {
		return time.Now().UTC().Add(30 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	}

	return time.Now().UTC().Add(1 * time.Hour).Format("2006-01-02 15:04:05")
}

func (s *sessions) Create(ctx context.Context, id string, keepLoggedIn bool) (*armadan.Session, error) {
	session, err := s.dbWriter.CreateSession(ctx, &schema.CreateSessionParams{
		UserID:    id,
		Token:     uuid.NewString(),
		ExpiresAt: getExpiration(keepLoggedIn),
	})

	if err != nil {
		return nil, err
	}

	return toSession(session), nil
}

func (s *sessions) DeleteByToken(ctx context.Context, token string) error {
	return s.dbWriter.DeleteSession(ctx, token)
}
