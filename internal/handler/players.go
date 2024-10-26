package handler

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/shopspring/decimal"
)

// TODO: Error handling!

func ManagePlayersView(ps armadan.PlayerService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		players, err := ps.All(r.Context())

		if err != nil {
			return
		}

		views.ManagePlayers(players).Render(r.Context(), w)
	})
}

func NewPlayer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		partials.AddPlayer().Render(r.Context(), w)
	})
}

func EditPlayer(ps armadan.PlayerService, v armadan.Validator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)
		if err != nil {
			return
		}

		player, err := ps.Get(r.Context(), *id)
		if err != nil {
			return
		}

		partials.EditPlayer(*player).Render(r.Context(), w)
	})
}

func CancelEditPlayer(ps armadan.PlayerService, v armadan.Validator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)

		if err != nil {
			return
		}

		player, err := ps.Get(r.Context(), *id)

		if err != nil {
			return
		}

		partials.Player(*player).Render(r.Context(), w)
	})
}

func InsertPlayer(ps armadan.PlayerService, v armadan.Validator) http.Handler {
	type createPlayerData struct {
		FirstName string          `json:"firstName" validate:"required"`
		LastName  string          `json:"lastName" validate:"required"`
		Email     string          `json:"email" validate:"required,email"`
		HCP       decimal.Decimal `json:"hcp" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := createPlayerData{}
		if err := v.Validate(r, &data); err != nil {
			return
		}

		_, err := ps.Create(r.Context(), &armadan.Player{
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			Hcp:       data.HCP,
		})

		if err != nil {
			return
		}

		players, err := ps.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, partials.PlayerList(players)).
			WithSuccess("Spelare sparad").
			HTML()
	})
}

func UpdatePlayer(ps armadan.PlayerService, v armadan.Validator) http.Handler {
	type updatePlayerData struct {
		FirstName string          `json:"firstName" validate:"required"`
		LastName  string          `json:"lastName" validate:"required"`
		Email     string          `json:"email" validate:"required,email"`
		HCP       decimal.Decimal `json:"hcp" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)
		if err != nil {
			return
		}

		data := updatePlayerData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		_, err = ps.Update(r.Context(), &armadan.Player{
			ID:        *id,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			Hcp:       data.HCP,
		})

		if err != nil {
			return
		}

		players, err := ps.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, partials.PlayerList(players)).
			WithSuccess("Spelare uppdaterad").
			HTML()
	})
}

func DeletePlayer(ps armadan.PlayerService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)
		if err != nil {
			return
		}

		err = ps.Delete(r.Context(), *id)

		if err != nil {
			return
		}

		partials.SuccessToast("Spelare borttagen").Render(r.Context(), w)
	})
}
