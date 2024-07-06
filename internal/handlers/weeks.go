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

func (h *Handler) InsertWeek(c echo.Context) error {
	week := schema.CreateWeekParams{}

	if err := validation.ValidateRequest(c, &week); err != nil {
		return err
	}

	_, err := h.db.CreateWeek(c.Request().Context(), &week)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	weeks, err := h.db.GetWeeks(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

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
	Nr       int32     `json:"nr"`
	CourseID uuid.UUID `json:"courseId"`
	TeeID    uuid.UUID `json:"teeId"`
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

	return partials.SuccessToast("Vecka borttagen").Render(c.Request().Context(), c.Response().Writer)
}
