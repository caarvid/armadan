package middleware

import "github.com/labstack/echo/v4"

func HxPushPath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("HX-Push-URL", c.Path())

		return next(c)
	}
}
