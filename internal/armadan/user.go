package armadan

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	All(context.Context) ([]User, error)
	Get(context.Context, uuid.UUID) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	UpdateRole(context.Context, uuid.UUID, Role) (*User, error)
}

type Role string

const (
	AdminRole     = "admin"
	ModeratorRole = "moderator"
	UserRole      = "user"
)

type User struct {
	ID    uuid.UUID
	Email string
	Hash  string
	Role  Role
}
