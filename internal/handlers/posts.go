package handlers

import (
	"net/http"

	"github.com/caarvid/armadan/internal/constants"
	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/utils/markdown"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/internal/validation"
	"github.com/caarvid/armadan/web/template/partials"
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

	return partials.EditPost(post).Render(c.Request().Context(), c.Response().Writer)
}

type previewPostData struct {
	Body string `json:"body" validate:"required"`
}

func (h *Handler) PreviewPost(c echo.Context) error {
	data := previewPostData{}

	if err := validation.ValidateRequest(c, &data); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, string(markdown.MdToHtml([]byte(data.Body))))
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

	_, err := h.db.UpdatePost(c.Request().Context(), &schema.UpdatePostParams{
		ID:     data.ID,
		Title:  data.Title,
		Body:   data.Body,
		Author: data.Author,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	posts, err := h.db.GetPosts(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	h.cache.Delete(constants.HomeCache)

	return response.
		New(c, partials.PostList(posts)).
		WithToast(response.Success, "Inlägg uppdaterat").
		HTML()
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

	posts, err := h.db.GetPosts(c.Request().Context())

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error")
	}

	h.cache.Delete(constants.HomeCache)

	return response.
		New(c, views.ManagePosts(posts)).
		WithToast(response.Success, "Inlägg sparat").
		HTML()
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

	h.cache.Delete(constants.HomeCache)

	return partials.SuccessToast("Inlägg borttaget").Render(c.Request().Context(), c.Response().Writer)
}
