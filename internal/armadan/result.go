package armadan

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ResultService interface {
	Get(context.Context, uuid.UUID) (*Result, error)
	GetRounds(context.Context, uuid.UUID) ([]Round, error)
	Create(context.Context, uuid.UUID) (*Result, error)
	CreateRound(context.Context, *Round, []Score) error
	Delete(context.Context, uuid.UUID) error
	DeleteRound(context.Context, uuid.UUID) error
	Leaderboard(context.Context) ([]Leader, error)
	LeaderboardSummary(context.Context, uuid.UUID) ([]LeaderSummary, error)
	ManagementView(context.Context) ([]ResultDetail, error)
	GetRemainingPlayers(context.Context, uuid.UUID) ([]Player, error)
}

type Round struct {
	ID         uuid.UUID
	NetIn      int32
	NetOut     int32
	NetTotal   int32
	GrossIn    int32
	GrossOut   int32
	GrossTotal int32
	OldHcp     decimal.Decimal
	NewHcp     decimal.Decimal
	PlayerID   uuid.UUID
	ResultID   uuid.UUID
	FirstName  string
	LastName   string
	Hcp        decimal.Decimal
}

type Score struct {
	ID      uuid.UUID
	HoleID  uuid.UUID
	Strokes int32
	Index   int32
	Par     int32
}

type Result struct {
	ID        uuid.UUID
	Slope     int32
	Cr        decimal.Decimal
	WeekNr    int32
	CourseID  uuid.UUID
	WeekID    uuid.UUID
	Published bool
}

type ResultDetail struct {
	ID                 uuid.UUID
	Nr                 int32
	IsFinals           bool
	CourseName         string
	TeeName            string
	ResultID           uuid.UUID
	Published          bool
	Participants       int64
	Winners            int64
	IsFirstUnpublished bool
}

type Leader struct {
	ID         uuid.UUID
	Name       string
	Points     int32
	NrOfRounds int32
}

type LeaderSummary struct {
	ID         uuid.UUID
	Nr         int32
	Points     int32
	HasResults bool
}
