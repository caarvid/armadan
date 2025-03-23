package validation

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdParam(t *testing.T) {
	v := New()

	tests := map[string]struct {
		expectsError bool
		path         string
		value        string
	}{
		"valid id": {
			expectsError: false,
			path:         "id",
			value:        uuid.NewString(),
		},
		"invalid id": {
			expectsError: true,
			path:         "id",
			value:        "invalid_id",
		},
		"missing id": {
			expectsError: true,
			path:         "not_id",
			value:        "invalid_id",
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.SetPathValue(test.path, test.value)
			id, err := v.ValidateIdParam(req, "id")

			if test.expectsError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.value, id)
				assert.NoError(t, err)
			}

		})
	}
}

type data struct {
	Name string `json:"name" validate:"required"`
}

func TestValidate(t *testing.T) {
	v := New()

	tests := map[string]struct {
		expectsError bool
		body         io.Reader
		expected     string
	}{
		"valid json": {
			expectsError: false,
			body:         strings.NewReader("{\"name\":\"test\"}"),
			expected:     "test",
		},
		"invalid json": {
			expectsError: true,
			body:         strings.NewReader("{"),
			expected:     "",
		},
		"empty body": {
			expectsError: true,
			body:         nil,
			expected:     "",
		},
		"missing required field": {
			expectsError: true,
			body:         strings.NewReader("{\"email\":\"test\"}"),
			expected:     "",
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d := data{}
			req := httptest.NewRequest(http.MethodGet, "/", test.body)
			err := v.Validate(req, &d)

			if test.expectsError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expected, d.Name)
				assert.NoError(t, err)
			}
		})
	}
}
