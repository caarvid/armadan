package session

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
)

const sessionCookie = "armadan_sid"

func GetCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(sessionCookie)
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func NewCookie(w http.ResponseWriter, sess *armadan.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sess.Token,
		Expires:  sess.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
