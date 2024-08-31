package handlers

import (
	"net/http"
	"time"

	"github.com/caarvid/armadan/internal/constants"
	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ManageWeeksView(c echo.Context) error {
	weeks, err := h.db.GetWeeks(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	courses, err := h.db.GetCourses(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return views.ManageWeeks(weeks, courses).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CourseTees(c echo.Context) error {
	id, err := uuid.Parse(c.QueryParam("courseId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error")
	}

	tees, err := h.db.GetTeesByCourse(c.Request().Context(), id)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return partials.WeekTeeSelect(tees).Render(c.Request().Context(), c.Response().Writer)
}

type insertWeekData struct {
	Nr         int32       `json:"nr" validate:"required"`
	FinalsDate string      `json:"finalsDate" validate:"required_with=IsFinalWeek"`
	CourseID   uuid.UUID   `json:"courseId" validate:"required,uuid4"`
	TeeID      uuid.UUID   `json:"teeId" validate:"required,uuid4"`
	IsFinals   pgtype.Bool `json:"isFinalsWeek"`
}

func (h *Handler) InsertWeek(c echo.Context) error {
	week := insertWeekData{}

	if err := validation.ValidateRequest(c, &week); err != nil {
		return err
	}

	newWeek := &schema.CreateWeekParams{
		Nr:       week.Nr,
		CourseID: week.CourseID,
		TeeID:    week.TeeID,
		IsFinals: pgtype.Bool{Bool: false, Valid: true},
	}

	if week.IsFinals.Bool {
		finalsDate, _ := time.Parse(time.DateOnly, week.FinalsDate)
		newWeek.IsFinals = week.IsFinals
		newWeek.FinalsDate = pgtype.Timestamptz{
			Time:  finalsDate,
			Valid: true,
		}
	}

	_, err := h.db.CreateWeek(c.Request().Context(), newWeek)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	weeks, err := h.db.GetWeeks(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	h.cache.Delete(constants.ScheduleCache)

	return response.
		New(c, partials.WeekTable(weeks)).
		WithToast(response.Success, "Vecka sparad").
		HTML()
}

func (h *Handler) EditWeek(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	week, err := h.db.GetWeek(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	courses, err := h.db.GetCourses(c.Request().Context())

	if err != nil {
		return err
	}

	tees, err := h.db.GetTeesByCourse(c.Request().Context(), week.Course.Id)

	if err != nil {
		return err
	}

	return partials.EditWeekRow(week, courses, tees).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CancelEditWeek(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	week, err := h.db.GetWeek(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	return partials.WeekRow(schema.GetWeeksRow(week), utils.GetWeekDates(int(week.Nr))).Render(c.Request().Context(), c.Response().Writer)
}

type updateWeekData struct {
	ID       uuid.UUID `param:"id" validate:"required,uuid4"`
	Nr       int32     `json:"nr" validate:"required"`
	CourseID uuid.UUID `json:"courseId" validate:"required,uuid4"`
	TeeID    uuid.UUID `json:"teeId" validate:"required,uuid4"`
}

func (h *Handler) UpdateWeek(c echo.Context) error {
	data := updateWeekData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	_, err := h.db.UpdateWeek(c.Request().Context(), &schema.UpdateWeekParams{
		ID:       data.ID,
		Nr:       data.Nr,
		CourseID: data.CourseID,
		TeeID:    data.TeeID,
	})

	if err != nil {
		return err
	}

	weeks, err := h.db.GetWeeks(c.Request().Context())

	if err != nil {
		return err
	}

	h.cache.Delete(constants.ScheduleCache)

	return response.
		New(c, partials.WeekTable(weeks)).
		WithToast(response.Success, "Vecka uppdaterad").
		HTML()
}

func (h *Handler) DeleteWeek(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	err := h.db.DeleteWeek(c.Request().Context(), params.ID)

	if err != nil {
		return err
	}

	h.cache.Delete(constants.ScheduleCache)

	return partials.SuccessToast("Vecka borttagen").Render(c.Request().Context(), c.Response().Writer)
}
