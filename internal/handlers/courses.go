package handlers

import (
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ManageCoursesView(c echo.Context) error {
	courses, err := h.db.GetCourses(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "error")
	}

	return views.ManageCourses(courses).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CreateCourseView(c echo.Context) error {
	return views.CreateCourse().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) GetEmptyTeeForm(c echo.Context) error {
	return views.EmptyTeeForm().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) RemoveEmptyTeeForm(c echo.Context) error {
	return c.HTML(http.StatusOK, "")
}

func (h *Handler) RemoveTee(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	err := h.db.DeleteTee(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return partials.SuccessToast("Tee borttagen").Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CancelEditCourse(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	course, err := h.db.GetCourse(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return views.CourseCard(schema.GetCoursesRow(course)).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) EditCourse(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	course, err := h.db.GetCourse(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return views.EditCourse(course).Render(c.Request().Context(), c.Response().Writer)
}

type newHole struct {
	Nr    int32 `json:"nr" validate:"required"`
	Par   int32 `json:"par" validate:"requried"`
	Index int32 `json:"index" validate:"required"`
}

type newTee struct {
	Name  string         `json:"name" validate:"required"`
	Slope int32          `json:"slope" validate:"requried"`
	CR    pgtype.Numeric `json:"cr" validate:"required"`
}

type createCourseData struct {
	Name  string    `json:"name" validate:"required"`
	Holes []newHole `json:"holes" validate:"required"`
	Tees  []newTee  `json:"tees"`
}

func (h *Handler) InsertCourse(c echo.Context) error {
	data := createCourseData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	tx, err := h.pool.Begin(ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	defer tx.Rollback(ctx)
	qtx := h.db.WithTx(tx)

	var sum int32

	for _, newHole := range data.Holes {
		sum = sum + newHole.Par
	}

	course, err := qtx.CreateCourse(c.Request().Context(), &schema.CreateCourseParams{
		Name: data.Name,
		Par:  sum,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var holes []*schema.CreateHolesParams

	for _, newHole := range data.Holes {
		holes = append(holes, &schema.CreateHolesParams{
			Nr:       newHole.Nr,
			Index:    newHole.Index,
			Par:      newHole.Par,
			CourseID: course.ID,
		})
	}

	_, err = qtx.CreateHoles(ctx, holes)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if len(data.Tees) > 0 {
		var tees []*schema.CreateTeesParams

		for _, newTee := range data.Tees {
			tees = append(tees, &schema.CreateTeesParams{
				Name:     newTee.Name,
				Slope:    newTee.Slope,
				Cr:       newTee.CR,
				CourseID: course.ID,
			})
		}

		_, err = qtx.CreateTees(ctx, tees)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	tx.Commit(ctx)

	courses, err := h.db.GetCourses(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "error")
	}

	return response.
		New(c, views.ManageCourses(courses)).
		WithToast(response.Success, "Bana sparad").
		HTML()
}

type updatedHole struct {
	ID    uuid.UUID `json:"id" validate:"required,uuid4"`
	Nr    int32     `json:"nr" validate:"required"`
	Par   int32     `json:"par" validate:"requried"`
	Index int32     `json:"index" validate:"required"`
}

type updatedTee struct {
	ID    uuid.UUID      `json:"id"`
	Name  string         `json:"name" validate:"required"`
	Slope int32          `json:"slope" validate:"requried"`
	CR    pgtype.Numeric `json:"cr" validate:"required"`
}

type updateCourseData struct {
	ID    uuid.UUID     `param:"id" validate:"required,uuid4"`
	Name  string        `json:"name" validate:"required"`
	Holes []updatedHole `json:"holes" validate:"required"`
	Tees  []updatedTee  `json:"tees"`
}

func (h *Handler) UpdateCourse(c echo.Context) error {
	data := updateCourseData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	tx, err := h.pool.Begin(ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	defer tx.Rollback(ctx)
	qtx := h.db.WithTx(tx)

	var sum int32

	for _, newHole := range data.Holes {
		sum = sum + newHole.Par
	}

	course, err := qtx.UpdateCourse(c.Request().Context(), &schema.UpdateCourseParams{
		Name: data.Name,
		Par:  sum,
		ID:   data.ID,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var holes []*schema.UpdateHolesParams

	for _, newHole := range data.Holes {
		holes = append(holes, &schema.UpdateHolesParams{
			Nr:    newHole.Nr,
			ID:    newHole.ID,
			Index: newHole.Index,
			Par:   newHole.Par,
		})
	}

	t := qtx.UpdateHoles(ctx, holes)
	t.Exec(nil)

	if len(data.Tees) > 0 {
		var tees []*schema.CreateTeesParams
		var updatedTees []*schema.UpdateTeesParams
		emptyId := uuid.UUID{}

		for _, newTee := range data.Tees {
			if newTee.ID.String() == emptyId.String() {
				tees = append(tees, &schema.CreateTeesParams{
					Name:     newTee.Name,
					Slope:    newTee.Slope,
					Cr:       newTee.CR,
					CourseID: course.ID,
				})
			} else {
				updatedTees = append(updatedTees, &schema.UpdateTeesParams{
					ID:    newTee.ID,
					Name:  newTee.Name,
					Slope: newTee.Slope,
					Cr:    newTee.CR,
				})
			}
		}

		if len(tees) > 0 {
			_, err = qtx.CreateTees(ctx, tees)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}

		t := qtx.UpdateTees(ctx, updatedTees)
		t.Exec(nil)
	}

	tx.Commit(ctx)

	courses, err := h.db.GetCourses(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "error")
	}

	return response.
		New(c, partials.CourseList(courses)).
		WithToast(response.Success, "Bana uppdaterad").
		HTML()
}

func (h *Handler) DeleteCourse(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	err := h.db.DeleteCourse(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return partials.SuccessToast("Bana borttagen").Render(c.Request().Context(), c.Response().Writer)
}
