package handlers

import (
	"github.com/caarvid/armadan/web/template/views"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ManageResultsView(c echo.Context) error {
	players, err := h.db.GetPlayers(c.Request().Context())

	if err != nil {
		return err
	}

	return views.ManageResults(players).Render(c.Request().Context(), c.Response().Writer)
}
