package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func RequestLogger(next http.Handler) http.Handler {
	h := hlog.NewHandler(log.Logger)

	uaHandler := hlog.UserAgentHandler("user_agent")
	idHandler := hlog.RequestIDHandler("req_id", "")

	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			if strings.Index(r.URL.String(), "public") == -1 {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Stringer("path", r.URL).
					Int("status", status).
					Int("res_size_bytes", size).
					Dur("res_time_ms", duration).
					Msg("request")
			}
		},
	)

	return h(accessHandler(idHandler(uaHandler(next))))
}
