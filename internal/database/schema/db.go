// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package schema

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createCourseStmt, err = db.PrepareContext(ctx, createCourse); err != nil {
		return nil, fmt.Errorf("error preparing query CreateCourse: %w", err)
	}
	if q.createHcpChangeStmt, err = db.PrepareContext(ctx, createHcpChange); err != nil {
		return nil, fmt.Errorf("error preparing query CreateHcpChange: %w", err)
	}
	if q.createHolesStmt, err = db.PrepareContext(ctx, createHoles); err != nil {
		return nil, fmt.Errorf("error preparing query CreateHoles: %w", err)
	}
	if q.createPlayerStmt, err = db.PrepareContext(ctx, createPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayer: %w", err)
	}
	if q.createPostStmt, err = db.PrepareContext(ctx, createPost); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePost: %w", err)
	}
	if q.createResultStmt, err = db.PrepareContext(ctx, createResult); err != nil {
		return nil, fmt.Errorf("error preparing query CreateResult: %w", err)
	}
	if q.createRoundStmt, err = db.PrepareContext(ctx, createRound); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRound: %w", err)
	}
	if q.createRoundDetailStmt, err = db.PrepareContext(ctx, createRoundDetail); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRoundDetail: %w", err)
	}
	if q.createScoresStmt, err = db.PrepareContext(ctx, createScores); err != nil {
		return nil, fmt.Errorf("error preparing query CreateScores: %w", err)
	}
	if q.createSessionStmt, err = db.PrepareContext(ctx, createSession); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSession: %w", err)
	}
	if q.createTeesStmt, err = db.PrepareContext(ctx, createTees); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTees: %w", err)
	}
	if q.createTokenStmt, err = db.PrepareContext(ctx, createToken); err != nil {
		return nil, fmt.Errorf("error preparing query CreateToken: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createWeekStmt, err = db.PrepareContext(ctx, createWeek); err != nil {
		return nil, fmt.Errorf("error preparing query CreateWeek: %w", err)
	}
	if q.deleteCourseStmt, err = db.PrepareContext(ctx, deleteCourse); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteCourse: %w", err)
	}
	if q.deletePostStmt, err = db.PrepareContext(ctx, deletePost); err != nil {
		return nil, fmt.Errorf("error preparing query DeletePost: %w", err)
	}
	if q.deleteResultStmt, err = db.PrepareContext(ctx, deleteResult); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteResult: %w", err)
	}
	if q.deleteRoundStmt, err = db.PrepareContext(ctx, deleteRound); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRound: %w", err)
	}
	if q.deleteSessionStmt, err = db.PrepareContext(ctx, deleteSession); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSession: %w", err)
	}
	if q.deleteTeeStmt, err = db.PrepareContext(ctx, deleteTee); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTee: %w", err)
	}
	if q.deleteTokenStmt, err = db.PrepareContext(ctx, deleteToken); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteToken: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.deleteWeekStmt, err = db.PrepareContext(ctx, deleteWeek); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteWeek: %w", err)
	}
	if q.getCourseStmt, err = db.PrepareContext(ctx, getCourse); err != nil {
		return nil, fmt.Errorf("error preparing query GetCourse: %w", err)
	}
	if q.getCoursesStmt, err = db.PrepareContext(ctx, getCourses); err != nil {
		return nil, fmt.Errorf("error preparing query GetCourses: %w", err)
	}
	if q.getLeaderboardStmt, err = db.PrepareContext(ctx, getLeaderboard); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeaderboard: %w", err)
	}
	if q.getLeaderboardSummaryStmt, err = db.PrepareContext(ctx, getLeaderboardSummary); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeaderboardSummary: %w", err)
	}
	if q.getManageResultViewStmt, err = db.PrepareContext(ctx, getManageResultView); err != nil {
		return nil, fmt.Errorf("error preparing query GetManageResultView: %w", err)
	}
	if q.getPlayerStmt, err = db.PrepareContext(ctx, getPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayer: %w", err)
	}
	if q.getPlayersStmt, err = db.PrepareContext(ctx, getPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayers: %w", err)
	}
	if q.getPostStmt, err = db.PrepareContext(ctx, getPost); err != nil {
		return nil, fmt.Errorf("error preparing query GetPost: %w", err)
	}
	if q.getPostsStmt, err = db.PrepareContext(ctx, getPosts); err != nil {
		return nil, fmt.Errorf("error preparing query GetPosts: %w", err)
	}
	if q.getRemainingPlayersByResultIdStmt, err = db.PrepareContext(ctx, getRemainingPlayersByResultId); err != nil {
		return nil, fmt.Errorf("error preparing query GetRemainingPlayersByResultId: %w", err)
	}
	if q.getResultByIdStmt, err = db.PrepareContext(ctx, getResultById); err != nil {
		return nil, fmt.Errorf("error preparing query GetResultById: %w", err)
	}
	if q.getRoundsByResultIdStmt, err = db.PrepareContext(ctx, getRoundsByResultId); err != nil {
		return nil, fmt.Errorf("error preparing query GetRoundsByResultId: %w", err)
	}
	if q.getSessionByTokenStmt, err = db.PrepareContext(ctx, getSessionByToken); err != nil {
		return nil, fmt.Errorf("error preparing query GetSessionByToken: %w", err)
	}
	if q.getTeesByCourseStmt, err = db.PrepareContext(ctx, getTeesByCourse); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeesByCourse: %w", err)
	}
	if q.getTokenStmt, err = db.PrepareContext(ctx, getToken); err != nil {
		return nil, fmt.Errorf("error preparing query GetToken: %w", err)
	}
	if q.getUserByEmailStmt, err = db.PrepareContext(ctx, getUserByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByEmail: %w", err)
	}
	if q.getUserByIdStmt, err = db.PrepareContext(ctx, getUserById); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserById: %w", err)
	}
	if q.getUsersStmt, err = db.PrepareContext(ctx, getUsers); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsers: %w", err)
	}
	if q.getWeekStmt, err = db.PrepareContext(ctx, getWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetWeek: %w", err)
	}
	if q.getWeeksStmt, err = db.PrepareContext(ctx, getWeeks); err != nil {
		return nil, fmt.Errorf("error preparing query GetWeeks: %w", err)
	}
	if q.updateCourseStmt, err = db.PrepareContext(ctx, updateCourse); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCourse: %w", err)
	}
	if q.updateHolesStmt, err = db.PrepareContext(ctx, updateHoles); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateHoles: %w", err)
	}
	if q.updatePlayerStmt, err = db.PrepareContext(ctx, updatePlayer); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayer: %w", err)
	}
	if q.updatePostStmt, err = db.PrepareContext(ctx, updatePost); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePost: %w", err)
	}
	if q.updateTeesStmt, err = db.PrepareContext(ctx, updateTees); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTees: %w", err)
	}
	if q.updateUserEmailStmt, err = db.PrepareContext(ctx, updateUserEmail); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserEmail: %w", err)
	}
	if q.updateUserPasswordStmt, err = db.PrepareContext(ctx, updateUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserPassword: %w", err)
	}
	if q.updateUserRoleStmt, err = db.PrepareContext(ctx, updateUserRole); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserRole: %w", err)
	}
	if q.updateWeekStmt, err = db.PrepareContext(ctx, updateWeek); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateWeek: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createCourseStmt != nil {
		if cerr := q.createCourseStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCourseStmt: %w", cerr)
		}
	}
	if q.createHcpChangeStmt != nil {
		if cerr := q.createHcpChangeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createHcpChangeStmt: %w", cerr)
		}
	}
	if q.createHolesStmt != nil {
		if cerr := q.createHolesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createHolesStmt: %w", cerr)
		}
	}
	if q.createPlayerStmt != nil {
		if cerr := q.createPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerStmt: %w", cerr)
		}
	}
	if q.createPostStmt != nil {
		if cerr := q.createPostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPostStmt: %w", cerr)
		}
	}
	if q.createResultStmt != nil {
		if cerr := q.createResultStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createResultStmt: %w", cerr)
		}
	}
	if q.createRoundStmt != nil {
		if cerr := q.createRoundStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRoundStmt: %w", cerr)
		}
	}
	if q.createRoundDetailStmt != nil {
		if cerr := q.createRoundDetailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRoundDetailStmt: %w", cerr)
		}
	}
	if q.createScoresStmt != nil {
		if cerr := q.createScoresStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createScoresStmt: %w", cerr)
		}
	}
	if q.createSessionStmt != nil {
		if cerr := q.createSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSessionStmt: %w", cerr)
		}
	}
	if q.createTeesStmt != nil {
		if cerr := q.createTeesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTeesStmt: %w", cerr)
		}
	}
	if q.createTokenStmt != nil {
		if cerr := q.createTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTokenStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.createWeekStmt != nil {
		if cerr := q.createWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createWeekStmt: %w", cerr)
		}
	}
	if q.deleteCourseStmt != nil {
		if cerr := q.deleteCourseStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteCourseStmt: %w", cerr)
		}
	}
	if q.deletePostStmt != nil {
		if cerr := q.deletePostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deletePostStmt: %w", cerr)
		}
	}
	if q.deleteResultStmt != nil {
		if cerr := q.deleteResultStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteResultStmt: %w", cerr)
		}
	}
	if q.deleteRoundStmt != nil {
		if cerr := q.deleteRoundStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRoundStmt: %w", cerr)
		}
	}
	if q.deleteSessionStmt != nil {
		if cerr := q.deleteSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSessionStmt: %w", cerr)
		}
	}
	if q.deleteTeeStmt != nil {
		if cerr := q.deleteTeeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTeeStmt: %w", cerr)
		}
	}
	if q.deleteTokenStmt != nil {
		if cerr := q.deleteTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTokenStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.deleteWeekStmt != nil {
		if cerr := q.deleteWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteWeekStmt: %w", cerr)
		}
	}
	if q.getCourseStmt != nil {
		if cerr := q.getCourseStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCourseStmt: %w", cerr)
		}
	}
	if q.getCoursesStmt != nil {
		if cerr := q.getCoursesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCoursesStmt: %w", cerr)
		}
	}
	if q.getLeaderboardStmt != nil {
		if cerr := q.getLeaderboardStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeaderboardStmt: %w", cerr)
		}
	}
	if q.getLeaderboardSummaryStmt != nil {
		if cerr := q.getLeaderboardSummaryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeaderboardSummaryStmt: %w", cerr)
		}
	}
	if q.getManageResultViewStmt != nil {
		if cerr := q.getManageResultViewStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getManageResultViewStmt: %w", cerr)
		}
	}
	if q.getPlayerStmt != nil {
		if cerr := q.getPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStmt: %w", cerr)
		}
	}
	if q.getPlayersStmt != nil {
		if cerr := q.getPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayersStmt: %w", cerr)
		}
	}
	if q.getPostStmt != nil {
		if cerr := q.getPostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPostStmt: %w", cerr)
		}
	}
	if q.getPostsStmt != nil {
		if cerr := q.getPostsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPostsStmt: %w", cerr)
		}
	}
	if q.getRemainingPlayersByResultIdStmt != nil {
		if cerr := q.getRemainingPlayersByResultIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRemainingPlayersByResultIdStmt: %w", cerr)
		}
	}
	if q.getResultByIdStmt != nil {
		if cerr := q.getResultByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getResultByIdStmt: %w", cerr)
		}
	}
	if q.getRoundsByResultIdStmt != nil {
		if cerr := q.getRoundsByResultIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRoundsByResultIdStmt: %w", cerr)
		}
	}
	if q.getSessionByTokenStmt != nil {
		if cerr := q.getSessionByTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSessionByTokenStmt: %w", cerr)
		}
	}
	if q.getTeesByCourseStmt != nil {
		if cerr := q.getTeesByCourseStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeesByCourseStmt: %w", cerr)
		}
	}
	if q.getTokenStmt != nil {
		if cerr := q.getTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTokenStmt: %w", cerr)
		}
	}
	if q.getUserByEmailStmt != nil {
		if cerr := q.getUserByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByEmailStmt: %w", cerr)
		}
	}
	if q.getUserByIdStmt != nil {
		if cerr := q.getUserByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByIdStmt: %w", cerr)
		}
	}
	if q.getUsersStmt != nil {
		if cerr := q.getUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsersStmt: %w", cerr)
		}
	}
	if q.getWeekStmt != nil {
		if cerr := q.getWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWeekStmt: %w", cerr)
		}
	}
	if q.getWeeksStmt != nil {
		if cerr := q.getWeeksStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWeeksStmt: %w", cerr)
		}
	}
	if q.updateCourseStmt != nil {
		if cerr := q.updateCourseStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCourseStmt: %w", cerr)
		}
	}
	if q.updateHolesStmt != nil {
		if cerr := q.updateHolesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateHolesStmt: %w", cerr)
		}
	}
	if q.updatePlayerStmt != nil {
		if cerr := q.updatePlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerStmt: %w", cerr)
		}
	}
	if q.updatePostStmt != nil {
		if cerr := q.updatePostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePostStmt: %w", cerr)
		}
	}
	if q.updateTeesStmt != nil {
		if cerr := q.updateTeesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTeesStmt: %w", cerr)
		}
	}
	if q.updateUserEmailStmt != nil {
		if cerr := q.updateUserEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserEmailStmt: %w", cerr)
		}
	}
	if q.updateUserPasswordStmt != nil {
		if cerr := q.updateUserPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserPasswordStmt: %w", cerr)
		}
	}
	if q.updateUserRoleStmt != nil {
		if cerr := q.updateUserRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserRoleStmt: %w", cerr)
		}
	}
	if q.updateWeekStmt != nil {
		if cerr := q.updateWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateWeekStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                DBTX
	tx                                *sql.Tx
	createCourseStmt                  *sql.Stmt
	createHcpChangeStmt               *sql.Stmt
	createHolesStmt                   *sql.Stmt
	createPlayerStmt                  *sql.Stmt
	createPostStmt                    *sql.Stmt
	createResultStmt                  *sql.Stmt
	createRoundStmt                   *sql.Stmt
	createRoundDetailStmt             *sql.Stmt
	createScoresStmt                  *sql.Stmt
	createSessionStmt                 *sql.Stmt
	createTeesStmt                    *sql.Stmt
	createTokenStmt                   *sql.Stmt
	createUserStmt                    *sql.Stmt
	createWeekStmt                    *sql.Stmt
	deleteCourseStmt                  *sql.Stmt
	deletePostStmt                    *sql.Stmt
	deleteResultStmt                  *sql.Stmt
	deleteRoundStmt                   *sql.Stmt
	deleteSessionStmt                 *sql.Stmt
	deleteTeeStmt                     *sql.Stmt
	deleteTokenStmt                   *sql.Stmt
	deleteUserStmt                    *sql.Stmt
	deleteWeekStmt                    *sql.Stmt
	getCourseStmt                     *sql.Stmt
	getCoursesStmt                    *sql.Stmt
	getLeaderboardStmt                *sql.Stmt
	getLeaderboardSummaryStmt         *sql.Stmt
	getManageResultViewStmt           *sql.Stmt
	getPlayerStmt                     *sql.Stmt
	getPlayersStmt                    *sql.Stmt
	getPostStmt                       *sql.Stmt
	getPostsStmt                      *sql.Stmt
	getRemainingPlayersByResultIdStmt *sql.Stmt
	getResultByIdStmt                 *sql.Stmt
	getRoundsByResultIdStmt           *sql.Stmt
	getSessionByTokenStmt             *sql.Stmt
	getTeesByCourseStmt               *sql.Stmt
	getTokenStmt                      *sql.Stmt
	getUserByEmailStmt                *sql.Stmt
	getUserByIdStmt                   *sql.Stmt
	getUsersStmt                      *sql.Stmt
	getWeekStmt                       *sql.Stmt
	getWeeksStmt                      *sql.Stmt
	updateCourseStmt                  *sql.Stmt
	updateHolesStmt                   *sql.Stmt
	updatePlayerStmt                  *sql.Stmt
	updatePostStmt                    *sql.Stmt
	updateTeesStmt                    *sql.Stmt
	updateUserEmailStmt               *sql.Stmt
	updateUserPasswordStmt            *sql.Stmt
	updateUserRoleStmt                *sql.Stmt
	updateWeekStmt                    *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                tx,
		tx:                                tx,
		createCourseStmt:                  q.createCourseStmt,
		createHcpChangeStmt:               q.createHcpChangeStmt,
		createHolesStmt:                   q.createHolesStmt,
		createPlayerStmt:                  q.createPlayerStmt,
		createPostStmt:                    q.createPostStmt,
		createResultStmt:                  q.createResultStmt,
		createRoundStmt:                   q.createRoundStmt,
		createRoundDetailStmt:             q.createRoundDetailStmt,
		createScoresStmt:                  q.createScoresStmt,
		createSessionStmt:                 q.createSessionStmt,
		createTeesStmt:                    q.createTeesStmt,
		createTokenStmt:                   q.createTokenStmt,
		createUserStmt:                    q.createUserStmt,
		createWeekStmt:                    q.createWeekStmt,
		deleteCourseStmt:                  q.deleteCourseStmt,
		deletePostStmt:                    q.deletePostStmt,
		deleteResultStmt:                  q.deleteResultStmt,
		deleteRoundStmt:                   q.deleteRoundStmt,
		deleteSessionStmt:                 q.deleteSessionStmt,
		deleteTeeStmt:                     q.deleteTeeStmt,
		deleteTokenStmt:                   q.deleteTokenStmt,
		deleteUserStmt:                    q.deleteUserStmt,
		deleteWeekStmt:                    q.deleteWeekStmt,
		getCourseStmt:                     q.getCourseStmt,
		getCoursesStmt:                    q.getCoursesStmt,
		getLeaderboardStmt:                q.getLeaderboardStmt,
		getLeaderboardSummaryStmt:         q.getLeaderboardSummaryStmt,
		getManageResultViewStmt:           q.getManageResultViewStmt,
		getPlayerStmt:                     q.getPlayerStmt,
		getPlayersStmt:                    q.getPlayersStmt,
		getPostStmt:                       q.getPostStmt,
		getPostsStmt:                      q.getPostsStmt,
		getRemainingPlayersByResultIdStmt: q.getRemainingPlayersByResultIdStmt,
		getResultByIdStmt:                 q.getResultByIdStmt,
		getRoundsByResultIdStmt:           q.getRoundsByResultIdStmt,
		getSessionByTokenStmt:             q.getSessionByTokenStmt,
		getTeesByCourseStmt:               q.getTeesByCourseStmt,
		getTokenStmt:                      q.getTokenStmt,
		getUserByEmailStmt:                q.getUserByEmailStmt,
		getUserByIdStmt:                   q.getUserByIdStmt,
		getUsersStmt:                      q.getUsersStmt,
		getWeekStmt:                       q.getWeekStmt,
		getWeeksStmt:                      q.getWeeksStmt,
		updateCourseStmt:                  q.updateCourseStmt,
		updateHolesStmt:                   q.updateHolesStmt,
		updatePlayerStmt:                  q.updatePlayerStmt,
		updatePostStmt:                    q.updatePostStmt,
		updateTeesStmt:                    q.updateTeesStmt,
		updateUserEmailStmt:               q.updateUserEmailStmt,
		updateUserPasswordStmt:            q.updateUserPasswordStmt,
		updateUserRoleStmt:                q.updateUserRoleStmt,
		updateWeekStmt:                    q.updateWeekStmt,
	}
}
