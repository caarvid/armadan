package handlers

import (
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ManageUsersView(c echo.Context) error {
	users, err := h.db.GetUsers(c.Request().Context())

	if err != nil {
		return err
	}

	return views.ManageUsers(users).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) EditUser(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	user, err := h.db.GetUserById(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	return partials.EditUserRow(user).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CancelEditUser(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	user, err := h.db.GetUserById(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	return partials.UserRow(schema.GetUsersRow(user)).Render(c.Request().Context(), c.Response().Writer)
}

type updateUserData struct {
	ID   uuid.UUID            `param:"id" validate:"required,uuid4"`
	Role schema.UsersRoleEnum `json:"role" validate:"required"`
}

func (h *Handler) UpdateUser(c echo.Context) error {
	data := updateUserData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	user, err := h.db.UpdateUserRole(c.Request().Context(), &schema.UpdateUserRoleParams{
		ID:   data.ID,
		Role: data.Role,
	})

	if err != nil {
		return err
	}

	return response.
		New(c, partials.UserRow(schema.GetUsersRow{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		})).
		WithToast(response.Success, "Användare uppdaterad").
		HTML()
}

func (h *Handler) InsertAdminUser(c echo.Context) error {
	data := schema.CreateUserParams{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	hash, err := utils.GenerateHash(data.Password, nil)

	if err != nil {
		return err
	}

	user, err := h.db.CreateUser(c.Request().Context(), &schema.CreateUserParams{
		Email:    data.Email,
		Password: hash.Encode(),
	})

	if err != nil {
		return err
	}

	h.db.UpdateUserRole(c.Request().Context(), &schema.UpdateUserRoleParams{
		Role: schema.UsersRoleEnumAdmin,
		ID:   user.ID,
	})

	return c.NoContent(http.StatusNoContent)
}
