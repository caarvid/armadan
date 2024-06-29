package errors

import (
	"net/http"

	"github.com/caarvid/armadan/web/template/partials"
	"github.com/labstack/echo/v4"
)

func NewInvalidCredentialsError(c echo.Context) error {
	c.Response().Header().Add("HX-Retarget", "#login-error")
	c.Response().Header().Add("HX-Reselect", "#login-error")
	c.Response().WriteHeader(http.StatusUnprocessableEntity)
	return partials.WrongCredentials().Render(c.Request().Context(), c.Response().Writer)
}
