package middleware

import (
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/labstack/echo/v4"
)

func redirectToLogin(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/login")
}

func redirectToHome(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Add("HX-Redirect", "/")
		return c.NoContent(http.StatusForbidden)
	}

	return c.Redirect(http.StatusFound, "/")
}

func Authorize(db *schema.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sCookie, err := utils.GetSessionCookie(c)

			if err != nil {
				return redirectToLogin(c)
			}

			session, err := db.GetSessionByToken(c.Request().Context(), sCookie.Value)

			if err != nil || !session.IsValid() {
				c.SetCookie(utils.ClearSessionCookie())
				return redirectToLogin(c)
			}

			c.Set("isLoggedIn", true)
			c.Set("isModerator", utils.IsModerator(session.Role.UsersRoleEnum))
			c.Set("isAdmin", utils.IsAdmin(session.Role.UsersRoleEnum))

			return next(c)
		}
	}
}
