package handlers

import (
	"net/http"

	"github.com/caarvid/armadan/web/template/views"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HomeView(c echo.Context) error {
	posts, err := h.db.GetPosts(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return views.Index(posts).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) ScheduleView(c echo.Context) error {
	weeks, err := h.db.GetWeeks(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return views.Schedule(weeks).Render(c.Request().Context(), c.Response().Writer)
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
