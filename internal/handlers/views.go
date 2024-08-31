package handlers

import (
	"bytes"
	"net/http"

	"github.com/caarvid/armadan/internal/constants"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

func (h *Handler) HomeView(c echo.Context) error {
	if data, found := h.cache.Get(constants.HomeCache); found {
		switch v := data.(type) {
		case *bytes.Buffer:
			return c.HTML(http.StatusOK, v.String())
		}
	}

	posts, err := h.db.GetPosts(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return response.New(c, views.Index(posts)).
		Cache(h.cache, constants.HomeCache, cache.NoExpiration).
		HTML()
}

func (h *Handler) ScheduleView(c echo.Context) error {
	if data, found := h.cache.Get(constants.ScheduleCache); found {
		switch v := data.(type) {
		case *bytes.Buffer:
			return c.HTML(http.StatusOK, v.String())
		}
	}

	weeks, err := h.db.GetWeeks(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return response.New(c, views.Schedule(weeks)).
		Cache(h.cache, constants.ScheduleCache, cache.NoExpiration).
		HTML()
}

func (h *Handler) LeaderboardView(c echo.Context) error {
	players, err := h.db.GetLeaderboard(c.Request().Context())

	if err != nil {
		return err
	}

	return views.Leaderboard(players).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) LoginView(c echo.Context) error {
	return views.Login().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) ForgotPasswordView(c echo.Context) error {
	return views.ForgotPassword().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) ResetPasswordView(c echo.Context) error {
	return views.ResetPassword().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) AdminView(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/admin/posts")
}
