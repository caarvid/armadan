// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package schema

import (
	"context"
)

type Querier interface {
	CreateCourse(ctx context.Context, arg *CreateCourseParams) (Course, error)
	CreateHcpChange(ctx context.Context, arg *CreateHcpChangeParams) (HcpChange, error)
	CreateHoles(ctx context.Context, arg *CreateHolesParams) (Hole, error)
	CreatePlayer(ctx context.Context, arg *CreatePlayerParams) (Player, error)
	CreatePost(ctx context.Context, arg *CreatePostParams) (Post, error)
	CreateResult(ctx context.Context, arg *CreateResultParams) (Result, error)
	CreateRound(ctx context.Context, arg *CreateRoundParams) (Round, error)
	CreateRoundDetail(ctx context.Context, arg *CreateRoundDetailParams) (RoundDetail, error)
	CreateScores(ctx context.Context, arg *CreateScoresParams) (Score, error)
	CreateSession(ctx context.Context, arg *CreateSessionParams) (Session, error)
	CreateTees(ctx context.Context, arg *CreateTeesParams) (Tee, error)
	CreateToken(ctx context.Context, arg *CreateTokenParams) (PasswordResetToken, error)
	CreateUser(ctx context.Context, arg *CreateUserParams) (User, error)
	CreateWeek(ctx context.Context, arg *CreateWeekParams) (Week, error)
	DeleteCourse(ctx context.Context, id string) error
	DeletePost(ctx context.Context, id string) error
	DeleteResult(ctx context.Context, id string) error
	DeleteRound(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, token string) error
	DeleteTee(ctx context.Context, id string) error
	DeleteToken(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, id string) error
	DeleteWeek(ctx context.Context, id string) error
	GetCourse(ctx context.Context, id string) (CourseDetail, error)
	GetCourses(ctx context.Context) ([]CourseDetail, error)
	GetLeaderboard(ctx context.Context) ([]GetLeaderboardRow, error)
	GetLeaderboardSummary(ctx context.Context, playerID string) ([]GetLeaderboardSummaryRow, error)
	GetManageResultView(ctx context.Context) ([]GetManageResultViewRow, error)
	GetPlayer(ctx context.Context, id string) (PlayersWithPoint, error)
	GetPlayers(ctx context.Context) ([]PlayersWithPoint, error)
	GetPost(ctx context.Context, id string) (Post, error)
	GetPosts(ctx context.Context) ([]Post, error)
	GetRemainingPlayersByResultId(ctx context.Context, resultID string) ([]Player, error)
	GetResultById(ctx context.Context, id string) (GetResultByIdRow, error)
	GetRoundsByResultId(ctx context.Context, resultID string) ([]GetRoundsByResultIdRow, error)
	GetSessionByToken(ctx context.Context, token string) (GetSessionByTokenRow, error)
	GetTeesByCourse(ctx context.Context, courseID string) ([]Tee, error)
	GetToken(ctx context.Context, token string) (PasswordResetToken, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetWeek(ctx context.Context, id string) (WeekDetail, error)
	GetWeeks(ctx context.Context) ([]WeekDetail, error)
	UpdateCourse(ctx context.Context, arg *UpdateCourseParams) (Course, error)
	UpdateHoles(ctx context.Context, arg *UpdateHolesParams) error
	UpdatePlayer(ctx context.Context, arg *UpdatePlayerParams) (Player, error)
	UpdatePost(ctx context.Context, arg *UpdatePostParams) (Post, error)
	UpdateTees(ctx context.Context, arg *UpdateTeesParams) error
	UpdateUserEmail(ctx context.Context, arg *UpdateUserEmailParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg *UpdateUserPasswordParams) (User, error)
	UpdateUserRole(ctx context.Context, arg *UpdateUserRoleParams) (User, error)
	UpdateWeek(ctx context.Context, arg *UpdateWeekParams) (Week, error)
}

var _ Querier = (*Queries)(nil)
