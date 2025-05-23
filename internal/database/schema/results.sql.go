// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: results.sql

package schema

import (
	"context"
	"database/sql"
)

const createResult = `-- name: CreateResult :one
INSERT INTO results (id, week_id) VALUES (?, ?) RETURNING id, is_published, week_id
`

type CreateResultParams struct {
	ID     string `json:"id"`
	WeekID string `json:"weekId"`
}

func (q *Queries) CreateResult(ctx context.Context, arg *CreateResultParams) (Result, error) {
	row := q.queryRow(ctx, q.createResultStmt, createResult, arg.ID, arg.WeekID)
	var i Result
	err := row.Scan(&i.ID, &i.IsPublished, &i.WeekID)
	return i, err
}

const createRound = `-- name: CreateRound :one
INSERT INTO rounds (id, player_id, result_id) VALUES (?, ?, ?) RETURNING id, player_id, result_id
`

type CreateRoundParams struct {
	ID       string `json:"id"`
	PlayerID string `json:"playerId"`
	ResultID string `json:"resultId"`
}

func (q *Queries) CreateRound(ctx context.Context, arg *CreateRoundParams) (Round, error) {
	row := q.queryRow(ctx, q.createRoundStmt, createRound, arg.ID, arg.PlayerID, arg.ResultID)
	var i Round
	err := row.Scan(&i.ID, &i.PlayerID, &i.ResultID)
	return i, err
}

const createRoundDetail = `-- name: CreateRoundDetail :one
INSERT INTO round_details (net_in, net_out, gross_in, gross_out, round_id) VALUES (?, ?, ?, ?, ?) RETURNING id, net_in, net_out, net_total, gross_in, gross_out, gross_total, round_id
`

type CreateRoundDetailParams struct {
	NetIn    int64  `json:"netIn"`
	NetOut   int64  `json:"netOut"`
	GrossIn  int64  `json:"grossIn"`
	GrossOut int64  `json:"grossOut"`
	RoundID  string `json:"roundId"`
}

func (q *Queries) CreateRoundDetail(ctx context.Context, arg *CreateRoundDetailParams) (RoundDetail, error) {
	row := q.queryRow(ctx, q.createRoundDetailStmt, createRoundDetail,
		arg.NetIn,
		arg.NetOut,
		arg.GrossIn,
		arg.GrossOut,
		arg.RoundID,
	)
	var i RoundDetail
	err := row.Scan(
		&i.ID,
		&i.NetIn,
		&i.NetOut,
		&i.NetTotal,
		&i.GrossIn,
		&i.GrossOut,
		&i.GrossTotal,
		&i.RoundID,
	)
	return i, err
}

const createScores = `-- name: CreateScores :one
INSERT INTO scores (round_id, hole_id, strokes) VALUES (?, ?, ?) RETURNING id, strokes, hole_id, round_id
`

type CreateScoresParams struct {
	RoundID string `json:"roundId"`
	HoleID  string `json:"holeId"`
	Strokes int64  `json:"strokes"`
}

func (q *Queries) CreateScores(ctx context.Context, arg *CreateScoresParams) (Score, error) {
	row := q.queryRow(ctx, q.createScoresStmt, createScores, arg.RoundID, arg.HoleID, arg.Strokes)
	var i Score
	err := row.Scan(
		&i.ID,
		&i.Strokes,
		&i.HoleID,
		&i.RoundID,
	)
	return i, err
}

const createWinner = `-- name: CreateWinner :one
INSERT INTO winners (id, points, player_id, week_id) VALUES (?, ?, ?, ?) RETURNING id, points, week_id, player_id
`

type CreateWinnerParams struct {
	ID       string `json:"id"`
	Points   int64  `json:"points"`
	PlayerID string `json:"playerId"`
	WeekID   string `json:"weekId"`
}

func (q *Queries) CreateWinner(ctx context.Context, arg *CreateWinnerParams) (Winner, error) {
	row := q.queryRow(ctx, q.createWinnerStmt, createWinner,
		arg.ID,
		arg.Points,
		arg.PlayerID,
		arg.WeekID,
	)
	var i Winner
	err := row.Scan(
		&i.ID,
		&i.Points,
		&i.WeekID,
		&i.PlayerID,
	)
	return i, err
}

const deleteResult = `-- name: DeleteResult :exec
DELETE FROM results WHERE id = ?
`

func (q *Queries) DeleteResult(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.deleteResultStmt, deleteResult, id)
	return err
}

const deleteRound = `-- name: DeleteRound :exec
DELETE FROM rounds WHERE id = ?
`

func (q *Queries) DeleteRound(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.deleteRoundStmt, deleteRound, id)
	return err
}

const deleteWinnersByWeek = `-- name: DeleteWinnersByWeek :exec
DELETE FROM winners WHERE week_id = ?
`

func (q *Queries) DeleteWinnersByWeek(ctx context.Context, weekID string) error {
	_, err := q.exec(ctx, q.deleteWinnersByWeekStmt, deleteWinnersByWeek, weekID)
	return err
}

const getLatestResult = `-- name: GetLatestResult :one
SELECT
  r.id, r.is_published, r.week_id,
  w.nr as week_nr
FROM results r
JOIN weeks w on w.id = r.week_id
WHERE r.is_published = 1
ORDER BY week_nr DESC
LIMIT 1
`

type GetLatestResultRow struct {
	ID          string `json:"id"`
	IsPublished int64  `json:"isPublished"`
	WeekID      string `json:"weekId"`
	WeekNr      int64  `json:"weekNr"`
}

func (q *Queries) GetLatestResult(ctx context.Context) (GetLatestResultRow, error) {
	row := q.queryRow(ctx, q.getLatestResultStmt, getLatestResult)
	var i GetLatestResultRow
	err := row.Scan(
		&i.ID,
		&i.IsPublished,
		&i.WeekID,
		&i.WeekNr,
	)
	return i, err
}

const getLeaderboardSummary = `-- name: GetLeaderboardSummary :many
SELECT
  weeks.id, weeks.nr, weeks.is_finals, weeks.finals_date, weeks.start_date, weeks.end_date, weeks.course_id, weeks.tee_id,
  coalesce(w.points, 0) AS points,
  coalesce(r.is_published, FALSE) AS has_results
FROM weeks
JOIN results r ON r.week_id = weeks.id
JOIN winners w ON w.week_id = weeks.id AND w.player_id = ?
ORDER BY weeks.nr ASC
`

type GetLeaderboardSummaryRow struct {
	ID         string         `json:"id"`
	Nr         int64          `json:"nr"`
	IsFinals   int64          `json:"isFinals"`
	FinalsDate sql.NullString `json:"finalsDate"`
	StartDate  string         `json:"startDate"`
	EndDate    string         `json:"endDate"`
	CourseID   string         `json:"courseId"`
	TeeID      string         `json:"teeId"`
	Points     int64          `json:"points"`
	HasResults int64          `json:"hasResults"`
}

func (q *Queries) GetLeaderboardSummary(ctx context.Context, playerID string) ([]GetLeaderboardSummaryRow, error) {
	rows, err := q.query(ctx, q.getLeaderboardSummaryStmt, getLeaderboardSummary, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLeaderboardSummaryRow
	for rows.Next() {
		var i GetLeaderboardSummaryRow
		if err := rows.Scan(
			&i.ID,
			&i.Nr,
			&i.IsFinals,
			&i.FinalsDate,
			&i.StartDate,
			&i.EndDate,
			&i.CourseID,
			&i.TeeID,
			&i.Points,
			&i.HasResults,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManageResultView = `-- name: GetManageResultView :many
SELECT
  w.id,
  w.nr,
  w.is_finals,
  w.course_name,
  w.tee_name,
  r.id AS result_id,
  coalesce(r.is_published, FALSE) as is_published,
  coalesce(rd.participants, 0) AS participants,
  coalesce(win.winners, 0) as winners
FROM week_details w
LEFT JOIN results r ON r.week_id = w.id
LEFT JOIN  (
  SELECT result_id, COUNT(*) as participants FROM rounds GROUP BY result_id
) rd ON rd.result_id = r.id
LEFT JOIN (
    SELECT week_id, COUNT(*) AS winners FROM winners GROUP BY week_id
) win ON win.week_id = w.id
GROUP BY w.id, r.id, win.winners
ORDER BY w.nr ASC
`

type GetManageResultViewRow struct {
	ID           string         `json:"id"`
	Nr           int64          `json:"nr"`
	IsFinals     int64          `json:"isFinals"`
	CourseName   string         `json:"courseName"`
	TeeName      string         `json:"teeName"`
	ResultID     sql.NullString `json:"resultId"`
	IsPublished  int64          `json:"isPublished"`
	Participants int64          `json:"participants"`
	Winners      int64          `json:"winners"`
}

func (q *Queries) GetManageResultView(ctx context.Context) ([]GetManageResultViewRow, error) {
	rows, err := q.query(ctx, q.getManageResultViewStmt, getManageResultView)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetManageResultViewRow
	for rows.Next() {
		var i GetManageResultViewRow
		if err := rows.Scan(
			&i.ID,
			&i.Nr,
			&i.IsFinals,
			&i.CourseName,
			&i.TeeName,
			&i.ResultID,
			&i.IsPublished,
			&i.Participants,
			&i.Winners,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRemainingPlayersByResultId = `-- name: GetRemainingPlayersByResultId :many
SELECT 
  p.id, p.first_name, p.last_name, p.user_id,
  cast(coalesce((
      SELECT h.new_hcp FROM hcp_changes h
      WHERE h.player_id = p.id AND datetime(h.valid_from) < datetime(w.end_date)
      ORDER BY datetime(h.valid_from) DESC
      LIMIT 1
  ), (
      SELECT h.new_hcp FROM hcp_changes h 
      WHERE h.player_id = p.id 
      ORDER BY datetime(h.valid_from) ASC
      LIMIT 1
  ), 36.0) as real) AS hcp
FROM players p
LEFT JOIN rounds r ON r.player_id = p.id AND r.result_id = ?
LEFT JOIN results res ON res.id = r.result_id
LEFT JOIN weeks w ON w.id = res.week_id
WHERE r.player_id IS NULL
ORDER BY p.last_name ASC, p.first_name ASC
`

type GetRemainingPlayersByResultIdRow struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	UserID    string  `json:"userId"`
	Hcp       float64 `json:"hcp"`
}

func (q *Queries) GetRemainingPlayersByResultId(ctx context.Context, resultID string) ([]GetRemainingPlayersByResultIdRow, error) {
	rows, err := q.query(ctx, q.getRemainingPlayersByResultIdStmt, getRemainingPlayersByResultId, resultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRemainingPlayersByResultIdRow
	for rows.Next() {
		var i GetRemainingPlayersByResultIdRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.UserID,
			&i.Hcp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultById = `-- name: GetResultById :one
SELECT
  r.id, r.is_published, r.week_id,
  w.nr as week_nr,
  w.start_date as week_start_date,
  w.end_date as week_end_date,
  w.course_id,
  t.slope,
  t.cr
FROM results r
JOIN weeks w ON w.id = r.week_id
JOIN tees t ON t.id = w.tee_id
WHERE r.id = ?
`

type GetResultByIdRow struct {
	ID            string  `json:"id"`
	IsPublished   int64   `json:"isPublished"`
	WeekID        string  `json:"weekId"`
	WeekNr        int64   `json:"weekNr"`
	WeekStartDate string  `json:"weekStartDate"`
	WeekEndDate   string  `json:"weekEndDate"`
	CourseID      string  `json:"courseId"`
	Slope         int64   `json:"slope"`
	Cr            float64 `json:"cr"`
}

func (q *Queries) GetResultById(ctx context.Context, id string) (GetResultByIdRow, error) {
	row := q.queryRow(ctx, q.getResultByIdStmt, getResultById, id)
	var i GetResultByIdRow
	err := row.Scan(
		&i.ID,
		&i.IsPublished,
		&i.WeekID,
		&i.WeekNr,
		&i.WeekStartDate,
		&i.WeekEndDate,
		&i.CourseID,
		&i.Slope,
		&i.Cr,
	)
	return i, err
}

const getResultSummaryByWeek = `-- name: GetResultSummaryByWeek :one
SELECT
  w.id, w.nr, w.is_finals, w.finals_date, w.start_date, w.end_date, w.course_id, w.tee_id,
  c.name AS course_name,
  t.name AS tee_name,
  COALESCE((SELECT nr FROM weeks WHERE nr = w.nr - 1), 0) AS previous_week,
  COALESCE((SELECT nr FROM weeks JOIN results r2 ON r2.week_id = weeks.id WHERE nr = w.nr + 1 AND r2.is_published = 1), 0) AS next_week,
  json_group_array(
    json_object(
      'id', rd.id,
      'total', rd2.net_total,
      'playerName', CONCAT(p.first_name, ' ', p.last_name),
      'points', COALESCE(wi.points, 0)
    )
  ) AS rounds
FROM weeks w
JOIN courses c ON c.id = w.course_id
JOIN tees t ON t.id = w.tee_id
JOIN results r ON r.week_id = w.id
JOIN rounds rd ON rd.result_id = r.id
JOIN round_details rd2 ON rd2.round_id = rd.id
JOIN players p ON p.id = rd.player_id
LEFT JOIN winners wi ON wi.week_id = w.id AND wi.player_id = p.id
WHERE w.nr = ?
GROUP BY w.id, t.name, c.name
`

type GetResultSummaryByWeekRow struct {
	ID           string         `json:"id"`
	Nr           int64          `json:"nr"`
	IsFinals     int64          `json:"isFinals"`
	FinalsDate   sql.NullString `json:"finalsDate"`
	StartDate    string         `json:"startDate"`
	EndDate      string         `json:"endDate"`
	CourseID     string         `json:"courseId"`
	TeeID        string         `json:"teeId"`
	CourseName   string         `json:"courseName"`
	TeeName      string         `json:"teeName"`
	PreviousWeek interface{}    `json:"previousWeek"`
	NextWeek     interface{}    `json:"nextWeek"`
	Rounds       interface{}    `json:"rounds"`
}

func (q *Queries) GetResultSummaryByWeek(ctx context.Context, nr int64) (GetResultSummaryByWeekRow, error) {
	row := q.queryRow(ctx, q.getResultSummaryByWeekStmt, getResultSummaryByWeek, nr)
	var i GetResultSummaryByWeekRow
	err := row.Scan(
		&i.ID,
		&i.Nr,
		&i.IsFinals,
		&i.FinalsDate,
		&i.StartDate,
		&i.EndDate,
		&i.CourseID,
		&i.TeeID,
		&i.CourseName,
		&i.TeeName,
		&i.PreviousWeek,
		&i.NextWeek,
		&i.Rounds,
	)
	return i, err
}

const getRoundById = `-- name: GetRoundById :one
SELECT id, player_id, result_id, net_in, net_out, net_total, gross_in, gross_out, gross_total, old_hcp, new_hcp, scores FROM full_rounds WHERE id = ?
`

func (q *Queries) GetRoundById(ctx context.Context, id string) (FullRound, error) {
	row := q.queryRow(ctx, q.getRoundByIdStmt, getRoundById, id)
	var i FullRound
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.ResultID,
		&i.NetIn,
		&i.NetOut,
		&i.NetTotal,
		&i.GrossIn,
		&i.GrossOut,
		&i.GrossTotal,
		&i.OldHcp,
		&i.NewHcp,
		&i.Scores,
	)
	return i, err
}

const getRoundsByResultId = `-- name: GetRoundsByResultId :many
SELECT
  r.id, r.player_id, r.result_id, r.net_in, r.net_out, r.net_total, r.gross_in, r.gross_out, r.gross_total, r.old_hcp, r.new_hcp, r.scores,
  p.first_name,
  p.last_name
FROM full_rounds r
JOIN players p ON p.id = r.player_id
WHERE r.result_id = ?
ORDER BY r.net_total ASC
`

type GetRoundsByResultIdRow struct {
	ID         string         `json:"id"`
	PlayerID   string         `json:"playerId"`
	ResultID   string         `json:"resultId"`
	NetIn      int64          `json:"netIn"`
	NetOut     int64          `json:"netOut"`
	NetTotal   sql.NullInt64  `json:"netTotal"`
	GrossIn    int64          `json:"grossIn"`
	GrossOut   int64          `json:"grossOut"`
	GrossTotal sql.NullInt64  `json:"grossTotal"`
	OldHcp     float64        `json:"oldHcp"`
	NewHcp     float64        `json:"newHcp"`
	Scores     sql.NullString `json:"scores"`
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
}

func (q *Queries) GetRoundsByResultId(ctx context.Context, resultID string) ([]GetRoundsByResultIdRow, error) {
	rows, err := q.query(ctx, q.getRoundsByResultIdStmt, getRoundsByResultId, resultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRoundsByResultIdRow
	for rows.Next() {
		var i GetRoundsByResultIdRow
		if err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.ResultID,
			&i.NetIn,
			&i.NetOut,
			&i.NetTotal,
			&i.GrossIn,
			&i.GrossOut,
			&i.GrossTotal,
			&i.OldHcp,
			&i.NewHcp,
			&i.Scores,
			&i.FirstName,
			&i.LastName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const publishRound = `-- name: PublishRound :exec
UPDATE results SET is_published = TRUE WHERE id = ?
`

func (q *Queries) PublishRound(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.publishRoundStmt, publishRound, id)
	return err
}
