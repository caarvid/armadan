package handlers

import (
	"net/http"

	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/views"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ManagePostsView(c echo.Context) error {
	posts, err := h.db.GetPosts(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	return views.ManagePosts(posts).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CreatePostView(c echo.Context) error {
	return views.CreatePost().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) EditPost(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	post, err := h.db.GetPost(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return views.EditPost(post).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) CancelEditPost(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	post, err := h.db.GetPost(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return views.AdminPostItem(post).Render(c.Request().Context(), c.Response().Writer)
}

type updatePostData struct {
	ID     uuid.UUID `param:"id" validate:"required,uuid4"`
	Title  string    `json:"title" validate:"required"`
	Body   string    `json:"body" validate:"required"`
	Author string    `json:"author" validate:"required"`
}

func (h *Handler) UpdatePost(c echo.Context) error {
	data := updatePostData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	post, err := h.db.UpdatePost(c.Request().Context(), &schema.UpdatePostParams{
		ID:     data.ID,
		Title:  data.Title,
		Body:   data.Body,
		Author: data.Author,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return views.AdminPostItem(post).Render(c.Request().Context(), c.Response().Writer)
}

type createPostData struct {
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
	Author string `json:"author" validate:"required"`
}

func (h *Handler) InsertPost(c echo.Context) error {
	data := createPostData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	_, err := h.db.CreatePost(c.Request().Context(), &schema.CreatePostParams{
		Title:  data.Title,
		Body:   data.Body,
		Author: data.Author,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return h.ManagePostsView(c)
}

func (h *Handler) DeletePost(c echo.Context) error {
	params := idParam{}

	if err := validation.ValidateRequest(c, &params); err != nil {
		return err
	}

	err := h.db.DeletePost(c.Request().Context(), params.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.HTML(http.StatusOK, "")
}
