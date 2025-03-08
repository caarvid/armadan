package handler

import (
	"fmt"
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/hcp"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
)

// TODO: Error handling!

func ManageResultsView(rs armadan.ResultService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		details, err := rs.ManagementView(r.Context())

		if err != nil {
			return
		}

		views.ManageResults(details).Render(r.Context(), w)
	})
}

func EditResultView(rs armadan.ResultService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		result, err := rs.Get(r.Context(), id)
		if err != nil {
			return
		}

		rounds, err := rs.GetRounds(r.Context(), id)
		if err != nil {
			return
		}

		players, err := rs.GetRemainingPlayers(r.Context(), id)
		if err != nil {
			return
		}

		views.EditResult(result, rounds, players).Render(r.Context(), w)
	})
}

func NewRoundForm(
	rs armadan.ResultService,
	cs armadan.CourseService,
	ps armadan.PlayerService,
	v armadan.Validator,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		playerId := r.URL.Query().Get("playerId")

		var (
			result *armadan.Result
			course *armadan.Course
			player *armadan.Player
		)

		if result, err = rs.Get(r.Context(), id); err != nil {
			return
		}

		if course, err = cs.Get(r.Context(), result.CourseID); err != nil {
			return
		}

		if player, err = ps.Get(r.Context(), playerId); err != nil {
			return
		}

		strokes := hcp.GetStrokes(player.Hcp, result.Cr, int(result.Slope), int(course.Par))

		partials.RoundForm(id, strokes, course, player).Render(r.Context(), w)
	})
}

func AddNewResult(rs armadan.ResultService, ps armadan.PlayerService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		fmt.Println(id)

		newResult, err := rs.Create(r.Context(), id)
		if err != nil {
			return
		}

		result, err := rs.Get(r.Context(), newResult.ID)
		if err != nil {
			return
		}

		rounds, err := rs.GetRounds(r.Context(), newResult.ID)
		if err != nil {
			return
		}

		players, err := ps.All(r.Context())
		if err != nil {
			return
		}

		w.Header().Add("HX-Push-URL", fmt.Sprintf("/admin/results/%s", newResult.ID))
		views.EditResult(result, rounds, players).Render(r.Context(), w)
	})
}

func InsertRound(
	rs armadan.ResultService,
	v armadan.Validator,
) http.Handler {
	type score struct {
		HoleID  string `json:"holeId" validate:"required,uuid4"`
		Strokes int64  `json:"strokes" validate:"required,gte=1,lte=10"`
		Index   int64  `json:"index"`
		Par     int64  `json:"par"`
	}

	type newRoundData struct {
		PlayerID string  `json:"playerId" validate:"required,uuid4"`
		HCP      float64 `json:"hcp"`
		Scores   []score `json:"scores" validate:"required,len=18,dive"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		round := newRoundData{}

		if err := v.Validate(r, &round); err != nil {
			return
		}

		var scores []armadan.Score
		for _, s := range round.Scores {
			scores = append(scores, armadan.Score{
				HoleID:  s.HoleID,
				Strokes: s.Strokes,
				Index:   s.Index,
				Par:     s.Par,
			})
		}

		err = rs.CreateRound(r.Context(), &armadan.Round{
			Hcp:      round.HCP,
			PlayerID: round.PlayerID,
			ResultID: id,
		}, scores)

		if err != nil {
			return
		}

		result, err := rs.Get(r.Context(), id)
		if err != nil {
			return
		}

		rounds, err := rs.GetRounds(r.Context(), id)
		if err != nil {
			return
		}

		players, err := rs.GetRemainingPlayers(r.Context(), id)
		if err != nil {
			return
		}

		response.
			New(w, r, partials.RoundTable(rounds)).
			WithPartial(partials.NewRoundPanel(result, players)).
			WithSuccess("Runda tillagd").
			HTML()
	})
}

func DeleteResult(rs armadan.ResultService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		if err = rs.Delete(r.Context(), id); err != nil {
			return
		}

		results, err := rs.ManagementView(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, views.ManageResults(results)).
			WithSuccess("Resultat borttaget").
			HTML()
	})
}

func DeleteRound(rs armadan.ResultService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		resultId, err := v.ValidateIdParam(r, "resultId")
		if err != nil {
			return
		}

		if err = rs.DeleteRound(r.Context(), id); err != nil {
			return
		}

		players, err := rs.GetRemainingPlayers(r.Context(), resultId)
		if err != nil {
			return
		}

		response.
			New(w, r, partials.PlayerDropdown(resultId, players)).
			WithSuccess("Runda borttagen").
			HTML()
	})
}
