package armadan

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PlayerService interface {
	All(context.Context) ([]Player, error)
	Get(context.Context, uuid.UUID) (*Player, error)
	Create(context.Context, *Player) (*Player, error)
	Update(context.Context, *Player) (*Player, error)
	Delete(context.Context, uuid.UUID) error
}

type Player struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Points    int32
	UserID    uuid.UUID
	Email     string
	Hcp       decimal.Decimal
}
