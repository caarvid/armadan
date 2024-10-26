package handler

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/rs/zerolog"
)

func HomeView(ps armadan.PostService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())
		posts, err := ps.All(r.Context())

		if err != nil {
			l.Error().Str("location", "handlers:HomeView").Err(err).Msg("failed to get posts")

			views.Index([]armadan.Post{}).Render(r.Context(), w)
			return
		}

		views.Index(posts).Render(r.Context(), w)
	})
}

func ScheduleView(ws armadan.WeekService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())
		weeks, err := ws.All(r.Context())

		if err != nil {
			l.Error().Str("location", "handlers:ScheduleView").AnErr("raw_err", err).Msg("failed to get weeks")

			views.Schedule([]armadan.Week{}).Render(r.Context(), w)
			return
		}

		views.Schedule(weeks).Render(r.Context(), w)
	})
}

func LeaderboardView(rs armadan.ResultService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())
		players, err := rs.Leaderboard(r.Context())

		if err != nil {
			l.Error().Str("location", "handlers:LeaderboardView").AnErr("raw_err", err).Msg("failed to get leaderboard")

			views.Leaderboard([]armadan.Leader{}).Render(r.Context(), w)
			return
		}

		views.Leaderboard(players).Render(r.Context(), w)
	})
}

func LoginView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Login().Render(r.Context(), w)
	})
}

func ForgotPasswordView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.ForgotPassword().Render(r.Context(), w)
	})
}

func ResetPasswordView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.ResetPassword().Render(r.Context(), w)
	})
}

func AdminView() http.Handler {
	return http.RedirectHandler("/admin/posts", http.StatusPermanentRedirect)
}
