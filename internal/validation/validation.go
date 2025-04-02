package validation

import (
	"bytes"
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

type validationError struct {
	Field   string
	Tag     string
	Message string
}
type FieldErrors []validationError

func (fe FieldErrors) Error() string {
	buff := bytes.NewBufferString("Validation failed:\n")

	for i := range fe {
		buff.WriteString(fe[i].Message)
		buff.WriteString("\n")
	}

	return buff.String()
}

func New() *customValidator {
	return &customValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (cv *customValidator) ValidateIdParam(r *http.Request, name string) (string, error) {
	v := r.PathValue(name)

	if v == "" {
		return "", errors.New("could not find id param")
	}

	err := cv.validator.Var(v, "required,uuid4")

	if err != nil {
		return "", err
	}

	val, err := uuid.Parse(v)

	if err != nil {
		return "", err
	}

	return val.String(), nil
}

func (cv *customValidator) Validate(r *http.Request, i any) error {
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
		}

		var vErr FieldErrors

		for _, err := range err.(validator.ValidationErrors) {
			vErr = append(vErr, validationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Message: fmt.Sprintf("Field: %s - tag: %s", err.Field(), err.Tag()),
			})
		}

		return vErr
	}

	return nil
}
