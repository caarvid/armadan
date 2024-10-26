package middleware

import (
	"fmt"
	"net/http"
)

func HxPushPath(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HX-Push-URL", fmt.Sprintf("%s%s", prefix, r.URL.Path))

		next.ServeHTTP(w, r)
	})
}
