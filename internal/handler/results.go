package handler

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/hcp"
	"github.com/caarvid/armadan/internal/utils/response"
	resultUtils "github.com/caarvid/armadan/internal/utils/result"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/rs/zerolog"
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

func PublishResultView(rs armadan.ResultService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:PublishResultView").Logger()

		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("validation failed")
			return
		}

		result, err := rs.Get(r.Context(), id)
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to get result")
			return
		}

		rounds, err := rs.GetRounds(r.Context(), id)
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to get rounds")
			return
		}

		winners := resultUtils.GetWinners(rounds)
		roundSummaries := make([]armadan.RoundSummary, len(rounds))

		for i, r := range rounds {
			winnerIdx := slices.IndexFunc(winners, func(w armadan.Winner) bool {
				return r.PlayerID == w.PlayerID
			})

			roundSummaries[i].ID = r.ID
			roundSummaries[i].Total = r.NetTotal
			roundSummaries[i].PlayerName = fmt.Sprintf("%s %s", r.FirstName, r.LastName)
			roundSummaries[i].Points = 0

			if winnerIdx > -1 {
				roundSummaries[i].Points = winners[winnerIdx].Points
			}
		}

		slices.SortFunc(roundSummaries, func(a, b armadan.RoundSummary) int {
			return int(b.Points) - int(a.Points)
		})

		views.PublishResult(result, roundSummaries).Render(r.Context(), w)
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

func EditRound(rs armadan.ResultService, ps armadan.PlayerService, cs armadan.CourseService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultId, err := v.ValidateIdParam(r, "resultId")
		if err != nil {
			return
		}

		roundId, err := v.ValidateIdParam(r, "roundId")
		if err != nil {
			return
		}

		var (
			result *armadan.Result
			course *armadan.Course
			player *armadan.Player
			round  *armadan.Round
		)

		if result, err = rs.Get(r.Context(), resultId); err != nil {
			return
		}

		if course, err = cs.Get(r.Context(), result.CourseID); err != nil {
			return
		}

		if round, err = rs.GetRound(r.Context(), roundId); err != nil {
			return
		}

		if player, err = ps.Get(r.Context(), round.PlayerID); err != nil {
			return
		}

		strokes := hcp.GetStrokes(round.OldHcp, result.Cr, int(result.Slope), int(course.Par))

		partials.EditRoundModal(resultId, strokes, course, player, round).Render(r.Context(), w)
	})
}

func AddNewResult(rs armadan.ResultService, ps armadan.PlayerService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

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

func UpdateRound(rs armadan.ResultService, v armadan.Validator) http.Handler {
	type score struct {
		HoleID  string `json:"holeId" validate:"required,uuid4"`
		Strokes int64  `json:"strokes" validate:"required,gte=1,lte=10"`
		Index   int64  `json:"index"`
		Par     int64  `json:"par"`
	}

	type roundData struct {
		PlayerID string  `json:"playerId" validate:"required,uuid4"`
		HCP      float64 `json:"hcp"`
		Scores   []score `json:"scores" validate:"required,len=18,dive"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:UpdateRound").Logger()

		resultId, err := v.ValidateIdParam(r, "resultId")
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("invalid resultId param")
			return
		}

		roundId, err := v.ValidateIdParam(r, "roundId")
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("invalid roundId param")
			return
		}

		round := roundData{}
		if err := v.Validate(r, &round); err != nil {
			l.Error().AnErr("raw_err", err).Msg("data validation failed")
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

		err = rs.UpdateRound(r.Context(), &armadan.Round{
			ID:       roundId,
			Hcp:      round.HCP,
			PlayerID: round.PlayerID,
			ResultID: resultId,
		}, scores)
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to update round")
			return
		}

		rounds, err := rs.GetRounds(r.Context(), resultId)
		if err != nil {
			return
		}

		response.
			New(w, r, partials.RoundTable(rounds)).
			WithSuccess("Runda uppdaterad").
			HTML()
	})
}

func PublishRound(rs armadan.ResultService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:PublishRound").Logger()

		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("validation failed")
			return
		}

		if err = rs.Publish(r.Context(), id); err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to publish results")
			return
		}

		results, err := rs.ManagementView(r.Context())
		if err != nil {
			l.Error().AnErr("raw_err", err).Msg("failed to get results view")
			return
		}

		w.Header().Add("HX-Push-URL", "/admin/results")
		response.
			New(w, r, views.ManageResults(results)).
			WithSuccess("Resultat publicerat").
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
