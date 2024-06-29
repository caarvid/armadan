package schema

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (sess *UserSession) IsValid() bool {
	return sess.IsActive && sess.ExpiresAt.Time.After(time.Now())
}

func getExpiration(keepLoggedIn bool) pgtype.Timestamptz {
	if keepLoggedIn {
		return pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true}
	}

	return pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true}
}

func NewSession(userId uuid.UUID, keepLoggedIn bool) *CreateSessionParams {
	return &CreateSessionParams{
		UserID:    userId,
		IsActive:  true,
		Token:     uuid.NewString(),
		ExpiresAt: getExpiration(keepLoggedIn),
	}
}
