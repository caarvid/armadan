package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

type customContext struct {
	echo.Context
}

func (ctx customContext) Get(key string) interface{} {
	val := ctx.Context.Get(key)

	if val != nil {
		return val
	}

	return ctx.Request().Context().Value(key)
}

func (ctx customContext) Set(key string, val interface{}) {
	ctx.SetRequest(ctx.Request().WithContext(
		context.WithValue(
			ctx.Request().Context(),
			key,
			val,
		),
	))
}

func DefaultContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &customContext{c}

		cc.Set("isLoggedIn", false)

		return next(cc)
	}
}
