package armadan

import (
	"context"
)

type PlayerService interface {
	All(context.Context) ([]Player, error)
	Get(context.Context, string) (*Player, error)
	GetByUserId(context.Context, string) (*Player, error)
	Create(context.Context, *Player) (*Player, error)
	Update(context.Context, *Player) (*Player, error)
	Delete(context.Context, string) error
}

type Player struct {
	ID         string
	FirstName  string
	LastName   string
	Points     int64
	UserID     string
	Email      string
	Hcp        float64
	NrOfRounds int64
}
