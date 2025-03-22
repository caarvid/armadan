package armadan

import (
	"context"
	"time"
)

type ResetPasswordService interface {
	Get(context.Context, string) (*ResetPasswordToken, error)
	Create(context.Context, string) (*ResetPasswordToken, error)
	UpdateUserPassword(context.Context, *ResetPasswordToken, string) error
}

type ResetPasswordToken struct {
	ID        int64
	UserId    string
	Token     string
	ExpiresAt time.Time
}

func (t *ResetPasswordToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}
