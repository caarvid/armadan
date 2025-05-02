package middleware

import (
	"context"
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/route"
	"github.com/caarvid/armadan/internal/utils/session"
	"github.com/rs/zerolog"
)

type middleware func(http.Handler) http.Handler

var roleMap = map[armadan.Role]int8{
	armadan.Role(armadan.UserRole):      1,
	armadan.Role(armadan.ModeratorRole): 2,
	armadan.Role(armadan.AdminRole):     3,
}

func Protected(
	s armadan.SessionService,
	role armadan.Role,
	loginHandler http.Handler,
) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := zerolog.Ctx(r.Context())
			sCookie, err := session.GetCookie(r)

			if err != nil {
				loginHandler.ServeHTTP(w, r)
				return
			}

			session, err := s.GetByToken(r.Context(), sCookie.Value)

			if err != nil || !session.IsValid() {
				loginHandler.ServeHTTP(w, r)
				return
			}

			l.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("user_id", session.UserID)
			})

			if roleMap[session.Role] < roleMap[role] {
				l.Info().
					Str("location", "middleware:authorize").
					Str("user_role", string(session.Role)).
					Str("required_role", string(role)).
					Msg("unauthorized access attempt")
				route.RedirectToHome(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "isLoggedIn", true)
			ctx = context.WithValue(ctx, "role", session.Role)
			ctx = context.WithValue(ctx, "user_id", session.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
