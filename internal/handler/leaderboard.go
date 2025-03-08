package handler

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/web/template/partials"
)

func GetLeaderboardSummary(rs armadan.ResultService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		summary, err := rs.LeaderboardSummary(r.Context(), id)
		if err != nil {
			return
		}

		partials.LeaderboardSummary(summary).Render(r.Context(), w)
	})
}
