package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

func DefaultContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetRequest(c.Request().WithContext(
			context.WithValue(
				c.Request().Context(),
				"isLoggedIn",
				false,
			),
		))

		return next(c)
	}
}
