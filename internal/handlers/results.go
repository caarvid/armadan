package handlers

import (
	"fmt"
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils/hcp"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func (h *Handler) ManageResultsView(c echo.Context) error {
	weeks, err := h.db.GetManageResultView(c.Request().Context())

	if err != nil {
		return err
	}

	return views.ManageResults(weeks).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) EditResultView(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	result, err := h.db.GetResultById(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	rounds, err := h.db.GetRoundsByResultId(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Push-URL", fmt.Sprintf("/admin/results/%s", params.ID))

	return views.EditResult(result, rounds).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) NewRound(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	players, err := h.db.GetPlayers(c.Request().Context())

	if err != nil {
		return err
	}

	return partials.NewRoundModal(params.ID.String(), players).Render(c.Request().Context(), c.Response().Writer)
}

// TODO: Improve performance?
func (h *Handler) NewRoundForm(c echo.Context) error {
	params := idParam{}
	playerId, err := uuid.Parse(c.QueryParam("playerId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error")
	}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	result, err := h.db.GetResultById(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	course, err := h.db.GetCourse(c.Request().Context(), result.CourseID)

	if err != nil {
		return err
	}

	player, err := h.db.GetPlayer(c.Request().Context(), playerId)

	if err != nil {
		return err
	}

	return partials.NewRoundForm(params.ID.String(), course, player).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) AddNewResult(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	_, err := h.db.CreateResult(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	return h.EditResultView(c)
}

type roundResults struct {
	NetIn    int32
	NetOut   int32
	GrossIn  int32
	GrossOut int32
}

type score struct {
	HoleID  uuid.UUID `json:"holeId" validate:"required,uuid4"`
	Strokes int32     `json:"strokes" validate:"required,gte=1,lte=10"`
	Index   int32     `json:"index"`
	Par     int32     `json:"par"`
}

type newRoundData struct {
	ID       uuid.UUID       `param:"id" validate:"required,uuid4"`
	CourseID uuid.UUID       `json:"courseId" validate:"required,uuid4"`
	PlayerID uuid.UUID       `json:"playerId" validate:"required,uuid4"`
	HCP      decimal.Decimal `json:"hcp"`
	Scores   []score         `json:"scores" validate:"required,len=18,dive"`
}

func getRoundResults(scores []score, strokes int) roundResults {
	results := roundResults{}

	for i, s := range scores {
		if i < 9 {
			results.GrossOut += s.Strokes
			results.NetOut += s.Strokes

			if s.Index <= int32(strokes) {
				results.NetOut -= 1
			}
		} else {
			results.GrossIn += s.Strokes
			results.NetIn += s.Strokes

			if s.Index <= int32(strokes) {
				results.NetIn -= 1
			}
		}
	}

	return results
}

func (h *Handler) InsertRound(c echo.Context) error {
	data := newRoundData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	result, err := h.db.GetResultById(c.Request().Context(), data.ID)

	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	tx, err := h.pool.Begin(ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	defer tx.Rollback(ctx)
	qtx := h.db.WithTx(tx)

	var par int32
	for _, h := range data.Scores {
		par += h.Par
	}

	strokes := hcp.GetStrokes(data.HCP.InexactFloat64(), result.Cr.InexactFloat64(), int(result.Slope.Int32), int(par))
	results := getRoundResults(data.Scores, strokes)
	newHcp := hcp.GetNewHcp(data.HCP.InexactFloat64(), par, results.NetIn+results.NetOut)

	newRound, err := qtx.CreateRound(ctx, &schema.CreateRoundParams{
		PlayerID: data.PlayerID,
		ResultID: data.ID,
		NewHcp:   decimal.NewFromFloat32(float32(newHcp)),
		OldHcp:   data.HCP,
		NetIn:    results.NetIn,
		NetOut:   results.NetOut,
		GrossIn:  results.GrossIn,
		GrossOut: results.GrossOut,
	})

	var scores []*schema.CreateScoresParams

	for _, s := range data.Scores {
		scores = append(scores, &schema.CreateScoresParams{
			Strokes: s.Strokes,
			RoundID: newRound.ID,
			HoleID:  s.HoleID,
		})
	}

	_, err = qtx.CreateScores(ctx, scores)

	if err != nil {
		return err
	}

	tx.Commit(ctx)

	rounds, err := h.db.GetRoundsByResultId(c.Request().Context(), data.ID)

	if err != nil {
		return err
	}

	return response.
		New(c, partials.RoundTable(rounds)).
		WithToast(response.Success, "Runda tillagd").
		HTML()
}

func (h *Handler) DeleteResult(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	err := h.db.DeleteResult(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	weeks, err := h.db.GetManageResultView(c.Request().Context())

	if err != nil {
		return err
	}

	return response.
		New(c, views.ManageResults(weeks)).
		WithToast(response.Success, "Resultat borttaget").
		HTML()
}

func (h *Handler) DeleteRound(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	err := h.db.DeleteRound(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	return partials.SuccessToast("Runda borttagen").Render(c.Request().Context(), c.Response().Writer)
}
