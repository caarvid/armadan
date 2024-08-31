package utils

import "github.com/caarvid/armadan/internal/schema"

func IsModerator(role schema.UsersRoleEnum) bool {
	return role == schema.UsersRoleEnumAdmin || role == schema.UsersRoleEnumModerator
}

func IsAdmin(role schema.UsersRoleEnum) bool {
	return role == schema.UsersRoleEnumAdmin
}
