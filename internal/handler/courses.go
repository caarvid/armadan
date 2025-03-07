package handler

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TODO: error handling

func ManageCoursesView(cs armadan.CourseService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		courses, err := cs.All(r.Context())

		if err != nil {
			return
		}

		views.ManageCourses(courses).Render(r.Context(), w)
	})
}

func CreateCourseView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.CreateCourse().Render(r.Context(), w)
	})
}

func GetEmptyTeeForm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.EmptyTeeForm().Render(r.Context(), w)
	})
}

func RemoveEmptyTeeForm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func RemoveTee(cs armadan.CourseService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		if err := cs.DeleteTee(r.Context(), *id); err != nil {
			return
		}

		partials.SuccessToast("Tee borttagen").Render(r.Context(), w)
	})
}

func EditCourse(cs armadan.CourseService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		course, err := cs.Get(r.Context(), *id)
		if err != nil {
			return
		}

		views.EditCourse(*course).Render(r.Context(), w)
	})
}

func InsertCourse(cs armadan.CourseService, v armadan.Validator) http.Handler {
	type newHole struct {
		Nr    int32 `json:"nr" validate:"required"`
		Par   int32 `json:"par" validate:"requried"`
		Index int32 `json:"index" validate:"required"`
	}

	type newTee struct {
		Name  string          `json:"name" validate:"required"`
		Slope int32           `json:"slope" validate:"requried"`
		CR    decimal.Decimal `json:"cr" validate:"required"`
	}

	type createCourseData struct {
		Name  string    `json:"name" validate:"required"`
		Holes []newHole `json:"holes" validate:"required,dive"`
		Tees  []newTee  `json:"tees" validate:"dive"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := createCourseData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		var (
			holes []armadan.Hole
			tees  []armadan.Tee
		)

		for _, h := range data.Holes {
			holes = append(holes, armadan.Hole{
				Nr:    h.Nr,
				Par:   h.Par,
				Index: h.Index,
			})
		}

		for _, t := range data.Tees {
			tees = append(tees, armadan.Tee{
				Slope: t.Slope,
				Cr:    t.CR,
				Name:  t.Name,
			})
		}

		_, err := cs.Create(r.Context(), &armadan.Course{
			Name:  data.Name,
			Holes: holes,
			Tees:  tees,
		})

		if err != nil {
			return
		}

		courses, err := cs.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, views.ManageCourses(courses)).
			WithSuccess("Bana sparad").
			HTML()
	})
}

func UpdateCourse(cs armadan.CourseService, v armadan.Validator) http.Handler {
	type updatedHole struct {
		ID    uuid.UUID `json:"id" validate:"required,uuid4"`
		Nr    int32     `json:"nr" validate:"required"`
		Par   int32     `json:"par" validate:"requried"`
		Index int32     `json:"index" validate:"required"`
	}

	type updatedTee struct {
		ID    uuid.UUID       `json:"id"`
		Name  string          `json:"name" validate:"required"`
		Slope int32           `json:"slope" validate:"requried"`
		CR    decimal.Decimal `json:"cr" validate:"required"`
	}

	type updateCourseData struct {
		Name  string        `json:"name" validate:"required"`
		Holes []updatedHole `json:"holes" validate:"required,dive"`
		Tees  []updatedTee  `json:"tees" validate:"dive"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		data := updateCourseData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		var (
			holes []armadan.Hole
			tees  []armadan.Tee
		)

		for _, h := range data.Holes {
			holes = append(holes, armadan.Hole{
				ID:    h.ID,
				Nr:    h.Nr,
				Par:   h.Par,
				Index: h.Index,
			})
		}

		for _, t := range data.Tees {
			tees = append(tees, armadan.Tee{
				ID:    t.ID,
				Slope: t.Slope,
				Cr:    t.CR,
				Name:  t.Name,
			})
		}

		_, err = cs.Update(r.Context(), &armadan.Course{
			ID:    *id,
			Name:  data.Name,
			Holes: holes,
			Tees:  tees,
		})

		if err != nil {
			return
		}

		courses, err := cs.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, views.ManageCourses(courses)).
			WithSuccess("Bana uppdaterad").
			HTML()
	})
}

func DeleteCourse(cs armadan.CourseService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		if err = cs.Delete(r.Context(), *id); err != nil {
			return
		}

		partials.SuccessToast("Bana borttagen").Render(r.Context(), w)
	})
}
