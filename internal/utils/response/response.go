package response

import (
	"bytes"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

const (
	Success = iota
	Error
)

type ResponseBuilder struct {
	context echo.Context
	buf     *bytes.Buffer
}

func New(ec echo.Context, comp templ.Component) *ResponseBuilder {
	buf := &bytes.Buffer{}

	comp.Render(ec.Request().Context(), buf)

	return &ResponseBuilder{
		buf:     buf,
		context: ec,
	}
}

func (rb *ResponseBuilder) Cache(c *cache.Cache, key string, dur time.Duration) *ResponseBuilder {
	c.Set(key, rb.buf, dur)

	return rb
}

func (rb *ResponseBuilder) WithToast(t int, msg string) *ResponseBuilder {
	switch t {
	case Success:
		partials.SuccessToast(msg).Render(rb.context.Request().Context(), rb.buf)
	}

	return rb
}

func (rb *ResponseBuilder) HTML() error {
	return rb.context.HTML(http.StatusOK, rb.buf.String())
}
