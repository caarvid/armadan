package user

import (
	"context"

	"github.com/caarvid/armadan/internal/armadan"
)

func IsLoggedIn(ctx context.Context) bool {
	if val, ok := ctx.Value("isLoggedIn").(bool); ok {
		return val
	}

	return false
}

func IsModerator(ctx context.Context) bool {
	if val, ok := ctx.Value("role").(armadan.Role); ok {
		return val == armadan.ModeratorRole || val == armadan.AdminRole
	}

	return false
}

func IsAdmin(ctx context.Context) bool {
	if val, ok := ctx.Value("role").(armadan.Role); ok {
		return val == armadan.AdminRole
	}

	return false
}
