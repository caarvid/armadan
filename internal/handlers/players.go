package handlers

import (
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.step.sm/crypto/randutil"
)

func (h *Handler) ManagePlayersView(c echo.Context) error {
	players, err := h.db.GetPlayers(c.Request().Context())

	if err != nil {
		return err
	}

	return views.ManagePlayers(players).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) EditPlayer(c echo.Context) error {
	p := &idParam{}

	if err := validation.ValidateRequest(c, p); err != nil {
		return err
	}

	player, err := h.db.GetPlayer(c.Request().Context(), p.ID)

	if err != nil {
		return err
	}

	return partials.PlayerCardForm(player).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CancelEditPlayer(c echo.Context) error {
	p := &idParam{}

	if err := validation.ValidateRequest(c, p); err != nil {
		return err
	}

	player, err := h.db.GetPlayer(c.Request().Context(), p.ID)

	if err != nil {
		return err
	}

	return partials.Player(schema.GetPlayersRow(player)).Render(c.Request().Context(), c.Response().Writer)
}

type createPlayerData struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

func (h *Handler) InsertPlayer(c echo.Context) error {
	data := &createPlayerData{}

	if err := validation.ValidateRequest(c, data); err != nil {
		return err
	}

	pw, err := randutil.Alphanumeric(24)

	if err != nil {
		return err
	}

	user, err := h.db.CreateUser(c.Request().Context(), &schema.CreateUserParams{
		Email:    data.Email,
		Password: pw,
	})

	if err != nil {
		return err
	}

	_, err = h.db.CreatePlayer(c.Request().Context(), &schema.CreatePlayerParams{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		UserID:    user.ID,
	})

	if err != nil {
		return err
	}

	players, err := h.db.GetPlayers(c.Request().Context())

	if err != nil {
		return err
	}

	return partials.PlayerList(players).Render(c.Request().Context(), c.Response().Writer)
}

type updatePlayerData struct {
	ID        uuid.UUID `param:"id" validate:"required,uuid4"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
}

func (h *Handler) UpdatePlayer(c echo.Context) error {
	data := &updatePlayerData{}

	if err := validation.ValidateRequest(c, data); err != nil {
		return err
	}

	tx, err := h.pool.Begin(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	defer tx.Rollback(c.Request().Context())
	qtx := h.db.WithTx(tx)

	player, err := qtx.UpdatePlayer(c.Request().Context(), &schema.UpdatePlayerParams{
		ID:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	})

	if err != nil {
		return err
	}

	_, err = qtx.UpdateUserEmail(c.Request().Context(), &schema.UpdateUserEmailParams{
		Email: data.Email,
		ID:    player.UserID,
	})

	if err != nil {
		return err
	}

	err = tx.Commit(c.Request().Context())

	if err != nil {
		return err
	}

	players, err := h.db.GetPlayers(c.Request().Context())

	if err != nil {
		return err
	}

	return partials.PlayerList(players).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) DeletePlayer(c echo.Context) error {
	if err := validation.ValidateRequest(c, &idParam{}); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "")
}