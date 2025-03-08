package armadan

import (
	"context"
)

type ResultService interface {
	Get(context.Context, string) (*Result, error)
	GetRounds(context.Context, string) ([]Round, error)
	Create(context.Context, string) (*Result, error)
	CreateRound(context.Context, *Round, []Score) error
	Delete(context.Context, string) error
	DeleteRound(context.Context, string) error
	Leaderboard(context.Context) ([]Leader, error)
	LeaderboardSummary(context.Context, string) ([]LeaderSummary, error)
	ManagementView(context.Context) ([]ResultDetail, error)
	GetRemainingPlayers(context.Context, string) ([]Player, error)
}

type Round struct {
	ID         string
	NetIn      int64
	NetOut     int64
	NetTotal   int64
	GrossIn    int64
	GrossOut   int64
	GrossTotal int64
	OldHcp     float64
	NewHcp     float64
	PlayerID   string
	ResultID   string
	FirstName  string
	LastName   string
	Hcp        float64
}

type Score struct {
	ID      string
	HoleID  string
	Strokes int64
	Index   int64
	Par     int64
}

type Result struct {
	ID        string
	Slope     int64
	Cr        float64
	WeekNr    int64
	CourseID  string
	WeekID    string
	Published bool
}

type ResultDetail struct {
	ID                 string
	Nr                 int64
	IsFinals           bool
	CourseName         string
	TeeName            string
	ResultID           string
	Published          bool
	Participants       int64
	Winners            int64
	IsFirstUnpublished bool
}

type Leader struct {
	ID         string
	Name       string
	Points     int64
	NrOfRounds int64
}

type LeaderSummary struct {
	ID         string
	Nr         int64
	Points     int64
	HasResults bool
}
