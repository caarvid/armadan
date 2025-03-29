package server

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/handler"
	h "github.com/caarvid/armadan/internal/handler"
	mw "github.com/caarvid/armadan/internal/middleware"
)

func setupRoutes(
	main *http.ServeMux,
	postService armadan.PostService,
	weekService armadan.WeekService,
	userService armadan.UserService,
	playerService armadan.PlayerService,
	sessionService armadan.SessionService,
	courseService armadan.CourseService,
	resultService armadan.ResultService,
	resetPasswordService armadan.ResetPasswordService,
	emailService armadan.EmailService,
	validator armadan.Validator,
) {
	api := http.NewServeMux()
	admin := http.NewServeMux()
	auth := http.NewServeMux()

	protected := mw.Protected(sessionService, armadan.UserRole, h.LoginView())
	moderatorOnly := mw.Protected(sessionService, armadan.ModeratorRole, h.LoginView())

	// Views
	main.Handle("GET /{$}", protected(mw.HxPushPath("", h.HomeView(postService))))
	main.Handle("GET /schedule", protected(mw.HxPushPath("", h.ScheduleView(weekService))))
	main.Handle("GET /leaderboard", protected(mw.HxPushPath("", h.LeaderboardView(resultService))))
	main.Handle("GET /login", mw.HxPushPath("", h.LoginView()))
	main.Handle("GET /forgot-password", mw.HxPushPath("", h.ForgotPasswordView()))
	main.Handle("GET /reset-password", mw.HxPushPath("", h.ResetPasswordView()))
	admin.Handle("GET /{$}", moderatorOnly(h.AdminView()))

	// Posts
	admin.Handle("GET /posts", mw.HxPushPath("/admin", h.ManagePostsView(postService)))
	admin.Handle("GET /posts/new", mw.HxPushPath("/admin", h.CreatePostView()))
	admin.Handle("GET /posts/{id}/edit", h.EditPost(postService, validator))
	api.Handle("POST /posts", h.InsertPost(postService, validator))
	api.Handle("POST /posts/preview", h.PreviewPost(validator))
	api.Handle("PUT /posts/{id}", h.UpdatePost(postService, validator))
	api.Handle("DELETE /posts/{id}", h.DeletePost(postService, validator))

	// Weeks
	admin.Handle("GET /weeks", mw.HxPushPath("/admin", h.ManageWeeksView(weekService, courseService)))
	admin.Handle("GET /weeks/tees", h.CourseTees(courseService, validator))
	admin.Handle("GET /weeks/{id}/edit", h.EditWeek(weekService, courseService, validator))
	admin.Handle("GET /weeks/{id}/edit/cancel", h.CancelEditWeek(weekService, validator))
	api.Handle("POST /weeks", h.InsertWeek(weekService, validator))
	api.Handle("PUT /weeks/{id}", h.UpdateWeek(weekService, validator))
	api.Handle("DELETE /weeks/{id}", h.DeleteWeek(weekService, validator))

	// Courses
	admin.Handle("GET /courses", mw.HxPushPath("/admin", h.ManageCoursesView(courseService)))
	admin.Handle("GET /courses/new", mw.HxPushPath("/admin", h.CreateCourseView()))
	admin.Handle("GET /courses/tee/new", h.GetEmptyTeeForm())
	admin.Handle("GET /courses/tee/remove", h.RemoveEmptyTeeForm())
	admin.Handle("GET /courses/{id}/edit", h.EditCourse(courseService, validator))
	api.Handle("POST /courses", h.InsertCourse(courseService, validator))
	api.Handle("PUT /courses/{id}", h.UpdateCourse(courseService, validator))
	api.Handle("DELETE /courses/{id}", h.DeleteCourse(courseService, validator))
	api.Handle("DELETE /courses/tee/{id}", h.RemoveTee(courseService, validator))

	// Players
	admin.Handle("GET /players", mw.HxPushPath("/admin", h.ManagePlayersView(playerService)))
	admin.Handle("GET /players/new", h.NewPlayer())
	admin.Handle("GET /players/{id}/edit", h.EditPlayer(playerService, validator))
	admin.Handle("GET /players/{id}/edit/cancel", h.CancelEditPlayer(playerService, validator))
	api.Handle("POST /players", h.InsertPlayer(playerService, validator))
	api.Handle("PUT /players/{id}", h.UpdatePlayer(playerService, validator))
	api.Handle("DELETE /players/{id}", h.DeletePlayer(playerService, validator))

	// Users
	admin.Handle("GET /users", mw.HxPushPath("/admin", h.ManageUsersView(userService)))
	admin.Handle("GET /users/{id}/edit", h.EditUser(userService, validator))
	admin.Handle("GET /users/{id}/edit/cancel", h.CancelEditUser(userService, validator))
	api.Handle("PUT /users/{id}", h.UpdateUser(userService, validator))

	// Results
	main.Handle("GET /leaderboard/{id}", protected(h.GetLeaderboardSummary(resultService, validator)))
	admin.Handle("GET /results", mw.HxPushPath("/admin", handler.ManageResultsView(resultService)))
	admin.Handle("GET /results/{id}", mw.HxPushPath("/admin", handler.EditResultView(resultService, validator)))
	admin.Handle("GET /results/{id}/publish", mw.HxPushPath("/admin", handler.PublishResultView(resultService, validator)))
	admin.Handle("GET /results/week/{id}/new", handler.AddNewResult(resultService, playerService, validator))
	admin.Handle("GET /results/{id}/form", handler.NewRoundForm(resultService, courseService, playerService, validator))
	admin.Handle("GET /results/{resultId}/round/{roundId}/edit", handler.EditRound(resultService, playerService, courseService, validator))
	api.Handle("POST /results/{id}/round", handler.InsertRound(resultService, validator))
	api.Handle("POST /results/{id}/publish", handler.PublishRound(resultService, validator))
	api.Handle("PUT /results/{resultId}/round/{roundId}", handler.UpdateRound(resultService, validator))
	api.Handle("DELETE /results/{resultId}/round/{id}", handler.DeleteRound(resultService, validator))
	api.Handle("DELETE /results/{id}", handler.DeleteResult(resultService, validator))

	// Auth
	auth.Handle("POST /login", h.Login(userService, sessionService, validator))
	auth.Handle("GET /logout", protected(h.Logout(sessionService)))
	auth.Handle("POST /forgot-password", h.ForgotPassword(userService, resetPasswordService, emailService, validator))
	auth.Handle("POST /reset-password", h.ResetPassword(resetPasswordService, validator))

	// Healthz
	api.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Static files
	main.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./web/static"))))

	main.Handle("/admin/", moderatorOnly(http.StripPrefix("/admin", admin)))
	main.Handle("/auth/", http.StripPrefix("/auth", auth))
	main.Handle("/api/", moderatorOnly(http.StripPrefix("/api", api)))
}
