package validation

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ValidateRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator(v *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validator: v,
	}
}

type errorMessage map[string]string
type validationErrors []errorMessage

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
		}

		var vErr validationErrors

		for _, err := range err.(validator.ValidationErrors) {
			e := make(errorMessage)

			e["field"] = err.Field()

			switch err.Tag() {
			case "required":
				e["message"] = "required field"
			case "email":
				e["message"] = "invalid email address"
			default:
				e["message"] = fmt.Sprintf("must satisfy '%s' '%v' criteria", err.Tag(), err.Param())
			}

			vErr = append(vErr, e)
		}

		return echo.NewHTTPError(http.StatusBadRequest, map[string]validationErrors{
			"errors": vErr,
		})
	}

	return nil
}
