package armadan

import (
	"context"
)

type UserService interface {
	All(context.Context) ([]User, error)
	Get(context.Context, string) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	UpdateRole(context.Context, string, string) (*User, error)
}

type Role string

const (
	AdminRole     = "admin"
	ModeratorRole = "moderator"
	UserRole      = "user"
)

type User struct {
	ID    string
	Email string
	Hash  string
	Role  Role
}
