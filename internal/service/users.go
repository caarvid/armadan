package service

import (
	"context"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/google/uuid"
)

func toUser(entity any) *armadan.User {
	switch u := entity.(type) {
	case schema.User:
		return &armadan.User{
			ID:    u.ID,
			Role:  armadan.Role(u.Role),
			Hash:  u.Password,
			Email: u.Email,
		}
	}

	return &armadan.User{}
}

type users struct {
	db schema.Querier
}

func NewUserService(db schema.Querier) *users {
	return &users{
		db: db,
	}
}

func (s *users) All(ctx context.Context) ([]armadan.User, error) {
	users, err := s.db.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(users, toUser), nil
}

func (s *users) Get(ctx context.Context, id uuid.UUID) (*armadan.User, error) {
	user, err := s.db.GetUserById(ctx, id)

	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func (s *users) GetByEmail(ctx context.Context, email string) (*armadan.User, error) {
	user, err := s.db.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func (s *users) UpdateRole(
	ctx context.Context,
	id uuid.UUID,
	role armadan.Role,
) (*armadan.User, error) {
	user, err := s.db.UpdateUserRole(ctx, &schema.UpdateUserRoleParams{
		ID:   id,
		Role: schema.UsersRoleEnum(role),
	})

	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}
