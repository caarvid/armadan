package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

func (h *Handler) GetLeaderboardSummary(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("leaderboard/%s", params.ID)

	if data, found := h.cache.Get(cacheKey); found {
		switch v := data.(type) {
		case *bytes.Buffer:
			return c.HTML(http.StatusOK, v.String())
		}
	}

	summary, err := h.db.GetLeaderboardSummary(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	return response.
		New(c, partials.LeaderboardSummary(summary)).
		Cache(h.cache, cacheKey, cache.DefaultExpiration).
		HTML()
}
