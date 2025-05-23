package handler

import (
	"net/http"
	"strings"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/utils/session"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/rs/zerolog"
)

func Login(us armadan.UserService, ss armadan.SessionService, v armadan.Validator) http.Handler {
	type loginData struct {
		Email        string `json:"email" validate:"required,email"`
		Password     string `json:"password" validate:"required"`
		KeepLoggedIn bool   `json:"keepLoggedIn"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:Login").Logger()

		data := &loginData{}
		if err := v.Validate(r, data); err != nil {
			l.Error().AnErr("raw_err", err).Msg("validation failed")
			response.LoginValidationError(w, r)
			return
		}

		user, err := us.GetByEmail(r.Context(), data.Email)
		if err != nil {
			l.Warn().Str("email", data.Email).AnErr("raw_err", err).Msg("login failed :: invalid email")
			response.InvalidCredentialsError(w, r)
			return
		}

		hash, err := utils.DecodeHash(user.Hash)
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to decode password")
			response.GeneralLoginError(w, r)
			return
		}

		match, _ := hash.Compare(data.Password)
		if !match {
			response.InvalidCredentialsError(w, r)
			return
		}

		sess, err := ss.Create(r.Context(), user.ID, data.KeepLoggedIn)
		if err != nil {
			l.Error().
				AnErr("raw_err", err).
				Str("user_id", user.ID).
				Msg("failed to create session")
			response.GeneralLoginError(w, r)
			return
		}

		session.NewCookie(w, sess)
		w.Header().Add("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	})
}

func Logout(s armadan.SessionService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:Logout").Logger()
		sess, err := session.GetCookie(r)

		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to get session cookie")
			response.GeneralError(w, r)
			return
		}

		err = s.DeleteByToken(r.Context(), sess.Value)

		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to delete session")
			response.GeneralError(w, r)
			return
		}

		session.ClearCookie(w)
		w.Header().Add("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	})
}

func ForgotPassword(us armadan.UserService, rps armadan.ResetPasswordService, es armadan.EmailService, v armadan.Validator) http.Handler {
	type forgotPasswordData struct {
		Email string `json:"email" validate:"required,email"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:forgotPassword").Logger()

		data := &forgotPasswordData{}
		if err := v.Validate(r, data); err != nil {
			l.Error().AnErr("raw_err", err).Msg("validation failed")
			return
		}

		email := strings.ToLower(data.Email)

		user, err := us.GetByEmail(r.Context(), email)
		if err != nil {
			l.Error().AnErr("raw_err", err).Msgf("reset password :: email %s not found", data.Email)
			response.ResetPasswordEmailSent(w, r, email)
			return
		}

		rsToken, err := rps.Create(r.Context(), user.ID)
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to create reset password token")
			response.ResetPasswordEmailSent(w, r, email)
			return
		}

		go func(token string) {
			if err := es.SendResetPassword(email, token); err != nil {
				l.Error().AnErr("raw_err", err).Msg("failed to send reset password email")
			}
		}(rsToken.Token)

		response.ResetPasswordEmailSent(w, r, email)
	})
}

func ResetPassword(rps armadan.ResetPasswordService, v armadan.Validator) http.Handler {
	type resetPasswordData struct {
		ResetToken     string `json:"resetToken"`
		NewPassword    string `json:"newPassword" validate:"required,min=8"`
		RepeatPassword string `json:"repeatPassword" validate:"required,min=8"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:resetPassword").Logger()

		data := &resetPasswordData{}
		if err := v.Validate(r, data); err != nil {
			l.Error().AnErr("raw_err", err).Msg("validation failed")

			// TODO: make this nicer?
			if vErr, ok := err.(validation.FieldErrors); ok {
				for _, e := range vErr {
					if e.Tag == "min" {
						response.ResetPasswordMessage(w, r, "Lösenord måste vara minst 8 karaktärer", "error")
						return
					}
				}
			}

			response.ResetPasswordMessage(w, r, "Något gick fel, försök igen", "error")
			return
		}

		if data.NewPassword != data.RepeatPassword {
			response.ResetPasswordMessage(w, r, "Lösenord måste vara lika", "error")
			return
		}

		token, err := rps.Get(r.Context(), data.ResetToken)
		if err != nil {
			response.ResetPasswordMessage(w, r, "Något gick fel, försök igen", "error")
			return
		}

		if token.IsExpired() {
			l.Warn().Msg("reset password failed :: token expired")
			response.ResetPasswordMessage(w, r, "Denna länk har gått ut", "error")
			return
		}

		if err = rps.UpdateUserPassword(r.Context(), token, data.NewPassword); err != nil {
			l.Error().AnErr("raw_err", err).Msg("reset password failed :: could not set new password")
			response.ResetPasswordMessage(w, r, "Något gick fel, försök igen", "error")
			return
		}

		w.Header().Add("HX-Replace-Url", "/login")
		response.New(w, r, views.Login()).WithSuccess("Lösenord återställt").HTML()
	})
}
