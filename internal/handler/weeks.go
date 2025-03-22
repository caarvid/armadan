package handler

import (
	"net/http"
	"time"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/rs/zerolog"
)

// TODO: Error handling!

func ManageWeeksView(ws armadan.WeekService, cs armadan.CourseService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context()).With().Str("location", "handlers:ManageWeeksView").Logger()

		weeks, err := ws.All(r.Context())
		if err != nil {
			l.Error().Err(err).Msg("could not get weeks")

			views.ManageWeeks([]armadan.Week{}, []armadan.Course{}).Render(r.Context(), w)
			return
		}

		courses, err := cs.All(r.Context())
		if err != nil {
			l.Error().Err(err).Msg("could not get courses")

			views.ManageWeeks([]armadan.Week{}, []armadan.Course{}).Render(r.Context(), w)
			return
		}

		views.ManageWeeks(weeks, courses).Render(r.Context(), w)
	})
}

func CourseTees(cs armadan.CourseService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("courseId")
		tees, err := cs.GetTees(r.Context(), id)

		if err != nil {
			return
		}

		partials.WeekTeeSelect(tees).Render(r.Context(), w)
	})
}

func InsertWeek(ws armadan.WeekService, v armadan.Validator) http.Handler {
	type insertWeekData struct {
		Nr         int64  `json:"nr" validate:"required"`
		FinalsDate string `json:"finalsDate" validate:"required_with=IsFinalWeek"`
		CourseID   string `json:"courseId" validate:"required,uuid4"`
		TeeID      string `json:"teeId" validate:"required,uuid4"`
		IsFinals   bool   `json:"isFinalsWeek"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		week := insertWeekData{}

		if err := v.Validate(r, &week); err != nil {
			return
		}

		newWeek := armadan.Week{
			Nr:       week.Nr,
			CourseID: week.CourseID,
			TeeID:    week.TeeID,
		}

		if week.IsFinals {
			finalsDate, _ := time.Parse(time.DateOnly, week.FinalsDate)
			newWeek.IsFinals = true
			newWeek.FinalsDate = finalsDate
		}

		if _, err := ws.Create(r.Context(), &newWeek); err != nil {
			return
		}

		weeks, err := ws.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, partials.WeekTable(weeks)).
			WithSuccess("Vecka sparad").
			HTML()
	})
}

func EditWeek(ws armadan.WeekService, cs armadan.CourseService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		var week *armadan.Week
		var courses []armadan.Course
		var tees []armadan.Tee

		if week, err = ws.Get(r.Context(), id); err != nil {
			return
		}

		if courses, err = cs.All(r.Context()); err != nil {
			return
		}

		if tees, err = cs.GetTees(r.Context(), week.CourseID); err != nil {
			return
		}

		partials.EditWeekRow(*week, courses, tees).Render(r.Context(), w)
	})
}

func CancelEditWeek(ws armadan.WeekService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		week, err := ws.Get(r.Context(), id)
		if err != nil {
			return
		}

		partials.WeekRow(*week).Render(r.Context(), w)
	})
}

func UpdateWeek(ws armadan.WeekService, v armadan.Validator) http.Handler {
	type updateWeekData struct {
		Nr       int64  `json:"nr" validate:"required"`
		CourseID string `json:"courseId" validate:"required,uuid4"`
		TeeID    string `json:"teeId" validate:"required,uuid4"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		data := updateWeekData{}
		if err := v.Validate(r, &data); err != nil {
			return
		}

		_, err = ws.Update(r.Context(), &armadan.Week{
			ID:       id,
			Nr:       data.Nr,
			CourseID: data.CourseID,
			TeeID:    data.TeeID,
		})

		if err != nil {
			return
		}

		weeks, err := ws.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, partials.WeekTable(weeks)).
			WithSuccess("Vecka uppdaterad").
			HTML()
	})
}

func DeleteWeek(ws armadan.WeekService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		if err = ws.Delete(r.Context(), id); err != nil {
			return
		}

		partials.SuccessToast("Vecka borttagen").Render(r.Context(), w)
	})
}
