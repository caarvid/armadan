package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type customValidator struct {
	validator *validator.Validate
}

type errorMessage map[string]string
type validationErrors []errorMessage

func New() *customValidator {
	return &customValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (cv *customValidator) ValidateIdParam(r *http.Request) (*uuid.UUID, error) {
	v := r.PathValue("id")

	if v == "" {
		return nil, errors.New("could not find id param")
	}

	err := cv.validator.Var(v, "required,uuid4")

	if err != nil {
		return nil, err
	}

	val, err := uuid.Parse(v)

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (cv *customValidator) Validate(r *http.Request, i interface{}) error {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, i); err != nil {
		return err
	}

	if err := cv.validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid data")
			// return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
		}

		var vErr validationErrors

		for _, err := range err.(validator.ValidationErrors) {
			e := make(errorMessage)

			e["field"] = err.Field()

			switch err.Tag() {
			case "required":
				e["message"] = "required field"
			case "email":
				e["message"] = "Ogiltig email"
			default:
				e["message"] = fmt.Sprintf("must satisfy '%s' '%v' criteria", err.Tag(), err.Param())
			}

			vErr = append(vErr, e)
		}

		return errors.New("validation failed")
		// return echo.NewHTTPError(http.StatusBadRequest, map[string]validationErrors{
		// 	"errors": vErr,
		// })
	}

	return nil
}
