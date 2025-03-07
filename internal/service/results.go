package service

import (
	"context"
	"fmt"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils/hcp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type results struct {
	db   schema.Querier
	pool *pgxpool.Pool
}

func NewResultService(db schema.Querier, pool *pgxpool.Pool) *results {
	return &results{db: db, pool: pool}
}

func toResult(entity any) *armadan.Result {
	switch e := entity.(type) {
	case schema.GetResultByIdRow:
		return &armadan.Result{
			ID:       e.ID,
			Slope:    e.Slope.Int32,
			Cr:       decimal.New(e.Cr.Int.Int64(), e.Cr.Exp),
			WeekNr:   e.WeekNr,
			CourseID: e.CourseID,
		}
	case schema.Result:
		return &armadan.Result{
			ID:        e.ID,
			WeekID:    e.ID,
			Published: e.Published,
		}
	}

	return &armadan.Result{}
}

func (rs *results) Get(ctx context.Context, id uuid.UUID) (*armadan.Result, error) {
	result, err := rs.db.GetResultById(ctx, id)

	if err != nil {
		return nil, err
	}

	return toResult(result), nil
}

func (rs *results) GetRounds(ctx context.Context, id uuid.UUID) ([]armadan.Round, error) {
	rounds, err := rs.db.GetRoundsByResultId(ctx, id)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(rounds, func(a any) *armadan.Round {
		switch r := a.(type) {
		case schema.GetRoundsByResultIdRow:
			return &armadan.Round{
				ID:         r.ID,
				NetIn:      r.NetIn,
				NetOut:     r.NetOut,
				NetTotal:   r.NetTotal,
				GrossIn:    r.GrossIn,
				GrossOut:   r.GrossOut,
				GrossTotal: r.GrossTotal,
				OldHcp:     r.OldHcp,
				NewHcp:     r.NewHcp,
				PlayerID:   r.PlayerID,
				ResultID:   r.ResultID,
				FirstName:  r.FirstName,
				LastName:   r.LastName,
				Hcp:        decimal.New(r.Hcp.Int.Int64(), r.Hcp.Exp),
			}
		}

		return &armadan.Round{}
	}), nil
}

func (rs *results) GetRemainingPlayers(ctx context.Context, id uuid.UUID) ([]armadan.Player, error) {
	players, err := rs.db.GetRemainingPlayersByResultId(ctx, id)
	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(players, toPlayer), nil
}

func (rs *results) Create(ctx context.Context, weekId uuid.UUID) (*armadan.Result, error) {
	result, err := rs.db.CreateResult(ctx, weekId)

	if err != nil {
		return nil, err
	}

	return toResult(result), nil
}

type roundResults struct {
	NetIn    int32
	NetOut   int32
	GrossIn  int32
	GrossOut int32
}

func getRoundSummary(scores []armadan.Score, strokes int) roundResults {
	results := roundResults{}

	for i, s := range scores {
		if i < 9 {
			results.GrossOut += s.Strokes
			results.NetOut += s.Strokes

			if s.Index <= int32(strokes) {
				results.NetOut -= 1
			}
		} else {
			results.GrossIn += s.Strokes
			results.NetIn += s.Strokes

			if s.Index <= int32(strokes) {
				results.NetIn -= 1
			}
		}
	}

	return results
}

func (rs *results) CreateRound(
	ctx context.Context,
	round *armadan.Round,
	scores []armadan.Score,
) error {
	result, err := rs.Get(ctx, round.ResultID)
	if err != nil {
		return err
	}

	tx, err := rs.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	qtx := schema.New(tx)

	var par int32
	for _, h := range scores {
		par += h.Par
	}

	strokes := hcp.GetStrokes(round.Hcp.InexactFloat64(), result.Cr.InexactFloat64(), int(result.Slope), int(par))
	roundSummary := getRoundSummary(scores, strokes)
	newHcp := decimal.NewFromFloat32(
		float32(hcp.GetNewHcp(round.Hcp.InexactFloat64(), par, roundSummary.NetIn+roundSummary.NetOut)))

	newRound, err := qtx.CreateRound(ctx, &schema.CreateRoundParams{
		PlayerID: round.PlayerID,
		ResultID: result.ID,
		NewHcp:   newHcp,
		OldHcp:   round.Hcp,
		NetIn:    roundSummary.NetIn,
		NetOut:   roundSummary.NetOut,
		GrossIn:  roundSummary.GrossIn,
		GrossOut: roundSummary.GrossOut,
	})

	var newScores []*schema.CreateScoresParams

	for _, s := range scores {
		newScores = append(newScores, &schema.CreateScoresParams{
			Strokes: s.Strokes,
			RoundID: newRound.ID,
			HoleID:  s.HoleID,
		})
	}

	if _, err = qtx.CreateScores(ctx, newScores); err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil
}

func (rs *results) Delete(ctx context.Context, id uuid.UUID) error {
	return rs.db.DeleteResult(ctx, id)
}

func (rs *results) DeleteRound(ctx context.Context, id uuid.UUID) error {
	return rs.db.DeleteRound(ctx, id)
}

func (rs *results) Leaderboard(ctx context.Context) ([]armadan.Leader, error) {
	board, err := rs.db.GetLeaderboard(ctx)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(board, func(a any) *armadan.Leader {
		switch l := a.(type) {
		case schema.GetLeaderboardRow:
			return &armadan.Leader{
				ID:         l.ID,
				Name:       fmt.Sprintf("%s %s", l.FirstName, l.LastName),
				Points:     l.Points,
				NrOfRounds: int32(l.NrOfRounds),
			}
		}

		return &armadan.Leader{}
	}), nil
}

func (rs *results) LeaderboardSummary(ctx context.Context, id uuid.UUID) ([]armadan.LeaderSummary, error) {
	summary, err := rs.db.GetLeaderboardSummary(ctx, id)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(summary, func(a any) *armadan.LeaderSummary {
		switch l := a.(type) {
		case schema.GetLeaderboardSummaryRow:
			return &armadan.LeaderSummary{
				ID:         l.ID,
				Nr:         l.Nr,
				Points:     l.Points,
				HasResults: l.HasResults,
			}
		}

		return &armadan.LeaderSummary{}
	}), nil
}

func (rs *results) ManagementView(ctx context.Context) ([]armadan.ResultDetail, error) {
	details, err := rs.db.GetManageResultView(ctx)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(details, func(a any) *armadan.ResultDetail {
		switch d := a.(type) {
		case schema.GetManageResultViewRow:
			detail := &armadan.ResultDetail{
				ID:           d.ID,
				Nr:           d.Nr,
				IsFinals:     d.IsFinals.Bool,
				CourseName:   d.CourseName,
				TeeName:      d.TeeName,
				ResultID:     d.ResultID,
				Published:    d.Published,
				Participants: d.Participants,
				Winners:      d.Winners,
				// TODO: Fix this
				IsFirstUnpublished: true,
			}

			return detail
		}

		return &armadan.ResultDetail{}
	}), nil
}
