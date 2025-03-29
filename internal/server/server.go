package server

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/middleware"
)

func New(
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
) http.Handler {
	mux := http.NewServeMux()

	setupRoutes(
		mux,
		postService,
		weekService,
		userService,
		playerService,
		sessionService,
		courseService,
		resultService,
		resetPasswordService,
		emailService,
		validator,
	)

	var handler http.Handler = mux

	handler = middleware.RequestLogger(handler)

	return handler
}
