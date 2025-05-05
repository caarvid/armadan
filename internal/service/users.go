package service

import (
	"context"
	"strings"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
)

func toUser(entity any) *armadan.User {
	switch u := entity.(type) {
	case schema.User:
		return &armadan.User{
			ID:    u.ID,
			Role:  armadan.Role(u.UserRole),
			Hash:  u.Password,
			Email: u.Email,
		}
	}

	return &armadan.User{}
}

type users struct {
	dbReader schema.Querier
	dbWriter schema.Querier
}

func NewUserService(reader, writer schema.Querier) *users {
	return &users{
		dbReader: reader,
		dbWriter: writer,
	}
}

func (s *users) All(ctx context.Context) ([]armadan.User, error) {
	users, err := s.dbReader.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(users, toUser), nil
}

func (s *users) Get(ctx context.Context, id string) (*armadan.User, error) {
	user, err := s.dbReader.GetUserById(ctx, id)

	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func (s *users) GetByEmail(ctx context.Context, email string) (*armadan.User, error) {
	user, err := s.dbReader.GetUserByEmail(ctx, strings.ToLower(email))

	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func (s *users) UpdateRole(
	ctx context.Context,
	id string,
	role string,
) (*armadan.User, error) {
	user, err := s.dbWriter.UpdateUserRole(ctx, &schema.UpdateUserRoleParams{
		ID:       id,
		UserRole: role,
	})

	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}
