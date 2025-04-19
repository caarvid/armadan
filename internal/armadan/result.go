package armadan

import (
	"context"
	"time"
)

type ResultService interface {
	Get(context.Context, string) (*Result, error)
	GetRounds(context.Context, string) ([]Round, error)
	GetRound(context.Context, string) (*Round, error)
	LatestResult(context.Context) (*Result, error)
	Create(context.Context, string) (*Result, error)
	CreateRound(context.Context, *Round, []Score) error
	UpdateRound(context.Context, *Round, []Score) error
	Delete(context.Context, string) error
	DeleteRound(context.Context, string) error
	Leaderboard(context.Context) ([]Leader, error)
	LeaderboardSummary(context.Context, string) ([]LeaderSummary, error)
	WeekSummary(context.Context, int64) (*WeeklyResult, error)
	ManagementView(context.Context) ([]ResultDetail, error)
	GetRemainingPlayers(context.Context, string) ([]Player, error)
	Publish(context.Context, string) error
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
	Scores     []Score
}

type RoundSummary struct {
	ID         string
	PlayerName string
	Total      int64
	Points     int64
}

type WeeklyResult struct {
	ID           string
	Nr           int64
	Course       string
	Tee          string
	NextWeek     int64
	PreviousWeek int64
	Rounds       []RoundSummary
}

type Score struct {
	ID      string
	HoleID  string
	Strokes int64
	Index   int64
	Par     int64
	Hole    Hole
}

type Result struct {
	ID            string
	Slope         int64
	Cr            float64
	WeekNr        int64
	WeekStartDate time.Time
	WeekEndDate   time.Time
	CourseID      string
	WeekID        string
	Published     bool
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

type Winner struct {
	ID       string
	Points   int64
	PlayerID string
	WeekID   string
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
