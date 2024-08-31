package handlers

import (
	// "time"

	"fmt"

	"github.com/caarvid/armadan/internal/errors"
	"github.com/caarvid/armadan/internal/schema"

	// "github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/labstack/echo/v4"
	// "go.step.sm/crypto/randutil"
)

type loginData struct {
	Email        string `json:"email" validate:"required"`
	Password     string `json:"password" validate:"required"`
	KeepLoggedIn bool   `json:"keepLoggedIn"`
}

type forgotPasswordData struct {
	Email string `json:"email" validate:"required"`
}

type resetPasswordData struct {
	NewPassword    string `json:"newPassword" validate:"required"`
	RepeatPassword string `json:"repeatPassword" validate:"required"`
}

func (h *Handler) ResetPassword(c echo.Context) error {
	return nil
}

func (h *Handler) ForgotPassword(c echo.Context) error {
	data := &forgotPasswordData{}

	if err := validation.ValidateRequest(c, data); err != nil {
		return err
	}

	// go func() {
	// 	user, err := h.db.GetUserByEmail(c.Request().Context(), data.Email)
	//
	// 	if err != nil {
	// 		return
	// 	}
	//
	// 	token, err := randutil.Alphanumeric(64)
	//
	// 	if err != nil {
	// 		return
	// 	}
	//
	// 	newToken := schema.CreateTokenParams{
	// 		UserID:    user.ID,
	// 		Hash:      token,
	// 	}
	//
	// 	_, err := h.db.CreateToken(c.Request().Context())
	//
	// 	if err != nil {
	// 		return
	// 	}
	// }()

	return nil
}

func (h *Handler) Login(c echo.Context) error {
	data := &loginData{}

	if err := validation.ValidateRequest(c, data); err != nil {
		return err
	}

	user, err := h.db.GetUserByEmail(c.Request().Context(), data.Email)

	if err != nil {
		return errors.NewInvalidCredentialsError(c)
	}

	hash, err := utils.DecodeHash(user.Password)

	if err != nil {
		fmt.Println("*** HERE 1")
		return err
	}

	match, err := hash.Compare(data.Password)

	if err != nil {
		fmt.Println("*** HERE 2")
		return err
	}

	if !match {
		return errors.NewInvalidCredentialsError(c)
	}

	sess, err := h.db.CreateSession(c.Request().Context(), schema.NewSession(user.ID, data.KeepLoggedIn))

	if err != nil {
		fmt.Println("*** HERE 3")
		return err
	}

	c.SetCookie(utils.NewSessionCookie(sess))
	c.Response().Header().Add("HX-Redirect", "/")

	return c.NoContent(200)
}

func (h *Handler) Logout(c echo.Context) error {
	session, err := utils.GetSessionCookie(c)

	if err != nil {
		return err
	}

	err = h.db.DeleteSession(c.Request().Context(), session.Value)

	if err != nil {
		return err
	}

	c.SetCookie(utils.ClearSessionCookie())
	c.Response().Header().Add("HX-Redirect", "/login")

	return c.NoContent(200)
}
