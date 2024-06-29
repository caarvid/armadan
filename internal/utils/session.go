package utils

import (
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/labstack/echo/v4"
)

const cookieName = "sId"

func GetSessionCookie(c echo.Context) (*http.Cookie, error) {
	return c.Request().Cookie(cookieName)
}

func ClearSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
}

func NewSessionCookie(sess schema.UserSession) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    sess.Token,
		Expires:  sess.ExpiresAt.Time,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
