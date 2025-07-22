package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils/hcp"
	resultUtils "github.com/caarvid/armadan/internal/utils/result"
	"github.com/patrickmn/go-cache"
)

type results struct {
	dbReader schema.Querier
	dbWriter schema.Querier
	pool     *sql.DB
	cache    *cache.Cache
}

func NewResultService(reader, writer schema.Querier, pool *sql.DB, cache *cache.Cache) *results {
	return &results{
		dbReader: reader,
		dbWriter: writer,
		pool:     pool,
		cache:    cache,
	}
}

func toResult(entity any) *armadan.Result {
	switch e := entity.(type) {
	case schema.GetResultByIdRow:
		return &armadan.Result{
			ID:            e.ID,
			Slope:         e.Slope,
			Cr:            e.Cr,
			WeekNr:        e.WeekNr,
			WeekStartDate: armadan.ParseTime(e.WeekStartDate),
			WeekEndDate:   armadan.ParseTime(e.WeekEndDate),
			WeekID:        e.WeekID,
			CourseID:      e.CourseID,
		}
	case schema.Result:
		return &armadan.Result{
			ID:        e.ID,
			WeekID:    e.ID,
			Published: e.IsPublished == 1,
		}
	}

	return &armadan.Result{}
}

func (rs *results) Get(ctx context.Context, id string) (*armadan.Result, error) {
	result, err := rs.dbReader.GetResultById(ctx, id)

	if err != nil {
		return nil, err
	}

	return toResult(result), nil
}

func (rs *results) GetRounds(ctx context.Context, id string) ([]armadan.Round, error) {
	rounds, err := rs.dbReader.GetRoundsByResultId(ctx, id)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(rounds, func(a any) *armadan.Round {
		switch r := a.(type) {
		case schema.GetRoundsByResultIdRow:
			var scores []armadan.Score

			if r.Scores.Valid {
				json.Unmarshal([]byte(r.Scores.String), &scores)
			}

			return &armadan.Round{
				ID:         r.ID,
				NetIn:      r.NetIn,
				NetOut:     r.NetOut,
				NetTotal:   r.NetTotal.Int64,
				GrossIn:    r.GrossIn,
				GrossOut:   r.GrossOut,
				GrossTotal: r.GrossTotal.Int64,
				OldHcp:     r.OldHcp,
				NewHcp:     r.NewHcp,
				PlayerID:   r.PlayerID,
				ResultID:   r.ResultID,
				FirstName:  r.FirstName,
				LastName:   r.LastName,
				Scores:     scores,
			}
		}

		return &armadan.Round{}
	}), nil
}

func (rs *results) GetRound(ctx context.Context, id string) (*armadan.Round, error) {
	round, err := rs.dbReader.GetRoundById(ctx, id)
	if err != nil {
		return nil, err
	}
	var scores []armadan.Score

	if round.Scores.Valid {
		json.Unmarshal([]byte(round.Scores.String), &scores)
	}

	return &armadan.Round{
		ID:         round.ID,
		NetIn:      round.NetIn,
		NetOut:     round.NetOut,
		NetTotal:   round.NetTotal.Int64,
		GrossIn:    round.GrossIn,
		GrossOut:   round.GrossOut,
		GrossTotal: round.GrossTotal.Int64,
		OldHcp:     round.OldHcp,
		NewHcp:     round.NewHcp,
		PlayerID:   round.PlayerID,
		ResultID:   round.ResultID,
		Scores:     scores,
	}, nil
}

func (rs *results) LatestResult(ctx context.Context) (*armadan.Result, error) {
	if cachedLatest, found := rs.cache.Get("latest_result"); found {
		return cachedLatest.(*armadan.Result), nil
	}

	res, err := rs.dbReader.GetLatestResult(ctx)
	if err != nil {
		return nil, err
	}

	latestResult := &armadan.Result{
		ID:     res.ID,
		WeekNr: res.WeekNr,
	}

	rs.cache.Set("latest_result", latestResult, cache.NoExpiration)

	return latestResult, nil
}

func (rs *results) WeekSummary(ctx context.Context, nr int64) (*armadan.WeeklyResult, error) {
	cacheKey := fmt.Sprintf("/results/%d", nr)

	if cachedSummary, found := rs.cache.Get(cacheKey); found {
		return cachedSummary.(*armadan.WeeklyResult), nil
	}

	res, err := rs.dbReader.GetResultSummaryByWeek(ctx, nr)
	if err != nil {
		return nil, err
	}

	r, ok := res.Rounds.(string)

	if !ok {
		return nil, errors.New("could not convert rounds to string")
	}

	var rounds []armadan.RoundSummary
	err = json.Unmarshal([]byte(r), &rounds)

	if err != nil {
		return nil, err
	}

	// TODO: Should be able to sort this in the sql query but sqlc is complaining
	slices.SortFunc(rounds, func(a, b armadan.RoundSummary) int {
		return int(a.Total - b.Total)
	})

	summary := &armadan.WeeklyResult{
		ID:           res.ID,
		Nr:           res.Nr,
		Course:       res.CourseName,
		Tee:          res.TeeName,
		NextWeek:     res.NextWeek.(int64),
		PreviousWeek: res.PreviousWeek.(int64),
		Rounds:       rounds,
	}

	rs.cache.Set(cacheKey, summary, cache.NoExpiration)

	return summary, nil
}

func (rs *results) GetRemainingPlayers(ctx context.Context, id string) ([]armadan.Player, error) {
	players, err := rs.dbReader.GetRemainingPlayersByResultId(ctx, id)
	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(players, toPlayer), nil
}

func (rs *results) Create(ctx context.Context, weekId string) (*armadan.Result, error) {
	result, err := rs.dbWriter.CreateResult(ctx, &schema.CreateResultParams{
		ID:     armadan.GetId(),
		WeekID: weekId,
	})

	if err != nil {
		return nil, err
	}

	return toResult(result), nil
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

	tx, err := rs.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	var par int64
	for _, h := range scores {
		par += h.Par
	}

	strokes := hcp.GetStrokes(round.Hcp, result.Cr, int(result.Slope), int(par))
	roundSummary := resultUtils.GetRoundSummary(scores, int64(strokes))
	newHcp := hcp.GetNewHcp(round.Hcp, par, roundSummary.NetIn+roundSummary.NetOut)

	roundId := armadan.GetId()
	_, err = qtx.CreateRound(ctx, &schema.CreateRoundParams{
		ID:       roundId,
		PlayerID: round.PlayerID,
		ResultID: result.ID,
	})

	if err != nil {
		return err
	}

	_, err = qtx.CreateRoundDetail(ctx, &schema.CreateRoundDetailParams{
		RoundID:  roundId,
		NetIn:    roundSummary.NetIn,
		NetOut:   roundSummary.NetOut,
		GrossIn:  roundSummary.GrossIn,
		GrossOut: roundSummary.GrossOut,
	})

	if err != nil {
		return err
	}

	_, err = qtx.CreateHcpChange(ctx, &schema.CreateHcpChangeParams{
		RoundID:   sql.NullString{String: roundId, Valid: true},
		PlayerID:  sql.NullString{String: round.PlayerID, Valid: true},
		ValidFrom: result.WeekEndDate.Format(armadan.DEFAULT_TIME_FORMAT),
		NewHcp:    newHcp,
		OldHcp:    round.Hcp,
	})

	if err != nil {
		return err
	}

	for _, s := range scores {
		_, err = qtx.CreateScores(ctx, &schema.CreateScoresParams{
			Strokes: s.Strokes,
			RoundID: roundId,
			HoleID:  s.HoleID,
		})

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (rs *results) UpdateRound(ctx context.Context, round *armadan.Round, scores []armadan.Score) error {
	result, err := rs.Get(ctx, round.ResultID)
	if err != nil {
		return err
	}

	tx, err := rs.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	if err := qtx.DeleteRound(ctx, round.ID); err != nil {
		return err
	}

	var par int64
	for _, h := range scores {
		par += h.Par
	}

	strokes := hcp.GetStrokes(round.Hcp, result.Cr, int(result.Slope), int(par))
	roundSummary := resultUtils.GetRoundSummary(scores, int64(strokes))
	newHcp := hcp.GetNewHcp(round.Hcp, par, roundSummary.NetIn+roundSummary.NetOut)

	roundId := armadan.GetId()
	_, err = qtx.CreateRound(ctx, &schema.CreateRoundParams{
		ID:       roundId,
		PlayerID: round.PlayerID,
		ResultID: result.ID,
	})

	if err != nil {
		return err
	}

	_, err = qtx.CreateRoundDetail(ctx, &schema.CreateRoundDetailParams{
		RoundID:  roundId,
		NetIn:    roundSummary.NetIn,
		NetOut:   roundSummary.NetOut,
		GrossIn:  roundSummary.GrossIn,
		GrossOut: roundSummary.GrossOut,
	})

	if err != nil {
		return err
	}

	_, err = qtx.CreateHcpChange(ctx, &schema.CreateHcpChangeParams{
		RoundID:   sql.NullString{String: roundId, Valid: true},
		PlayerID:  sql.NullString{String: round.PlayerID, Valid: true},
		ValidFrom: result.WeekEndDate.Format(armadan.DEFAULT_TIME_FORMAT),
		NewHcp:    newHcp,
		OldHcp:    round.Hcp,
	})

	if err != nil {
		return err
	}

	for _, s := range scores {
		_, err = qtx.CreateScores(ctx, &schema.CreateScoresParams{
			Strokes: s.Strokes,
			RoundID: roundId,
			HoleID:  s.HoleID,
		})

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (rs *results) Delete(ctx context.Context, id string) error {
	return rs.dbWriter.DeleteResult(ctx, id)
}

func (rs *results) DeleteRound(ctx context.Context, id string) error {
	return rs.dbWriter.DeleteRound(ctx, id)
}

func (rs *results) Leaderboard(ctx context.Context) ([]armadan.Leader, error) {
	board, err := rs.dbReader.GetLeaderboard(ctx)

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
				NrOfRounds: l.NrOfRounds,
			}
		}

		return &armadan.Leader{}
	}), nil
}

func (rs *results) LeaderboardSummary(ctx context.Context, id string) ([]armadan.LeaderSummary, error) {
	summary, err := rs.dbReader.GetLeaderboardSummary(ctx, id)

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
				HasResults: l.HasResults == 1,
			}
		}

		return &armadan.LeaderSummary{}
	}), nil
}

func (rs *results) ManagementView(ctx context.Context) ([]armadan.ResultDetail, error) {
	details, err := rs.dbReader.GetManageResultView(ctx)

	if err != nil {
		return nil, err
	}

	foundFirstUnpublished := false

	return armadan.MapEntities(details, func(a any) *armadan.ResultDetail {
		switch d := a.(type) {
		case schema.GetManageResultViewRow:
			isPublished := d.IsPublished == 1
			detail := &armadan.ResultDetail{
				ID:                 d.ID,
				Nr:                 d.Nr,
				IsFinals:           d.IsFinals == 1,
				CourseName:         d.CourseName,
				TeeName:            d.TeeName,
				ResultID:           d.ResultID.String,
				Published:          isPublished,
				Participants:       d.Participants,
				Winners:            d.Winners,
				IsFirstUnpublished: !isPublished && !foundFirstUnpublished,
			}

			if !isPublished && !foundFirstUnpublished {
				foundFirstUnpublished = true
			}

			return detail
		}

		return &armadan.ResultDetail{}
	}), nil
}

func (rs *results) Publish(ctx context.Context, id string) error {
	rs.cache.Flush()

	result, err := rs.Get(ctx, id)
	if err != nil {
		return err
	}

	rounds, err := rs.GetRounds(ctx, id)
	if err != nil {
		return err
	}

	tx, err := rs.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	winners := resultUtils.GetWinners(rounds)

	err = qtx.DeleteWinnersByWeek(ctx, result.WeekID)
	if err != nil {
		return err
	}

	for _, w := range winners {
		_, err = qtx.CreateWinner(ctx, &schema.CreateWinnerParams{
			ID:       w.ID,
			Points:   w.Points,
			PlayerID: w.PlayerID,
			WeekID:   result.WeekID,
		})
		if err != nil {
			return err
		}
	}

	err = qtx.PublishRound(ctx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
