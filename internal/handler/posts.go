package handler

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/markdown"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/rs/zerolog"

	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
)

// TODO: error handling

func ManagePostsView(ps armadan.PostService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())
		posts, err := ps.All(r.Context())

		if err != nil {
			l.Error().Str("location", "handlers:ManagePostView").AnErr("raw_err", err).Msg("could not get posts")

			views.ManagePosts([]armadan.Post{}).Render(r.Context(), w)
			return
		}

		views.ManagePosts(posts).Render(r.Context(), w)
	})
}

func CreatePostView() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.CreatePost().Render(r.Context(), w)
	})
}

func EditPost(ps armadan.PostService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		post, err := ps.Get(r.Context(), *id)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		partials.EditPost(*post).Render(r.Context(), w)
	})
}

func PreviewPost(v armadan.Validator) http.Handler {
	type previewPostData struct {
		Body string `json:"body" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := previewPostData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		w.Write(markdown.MdToHtml([]byte(data.Body)))
	})
}

func InsertPost(ps armadan.PostService, v armadan.Validator) http.Handler {
	type createPostData struct {
		Title  string `json:"title" validate:"required"`
		Body   string `json:"body" validate:"required"`
		Author string `json:"author" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := createPostData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		_, err := ps.Create(r.Context(), &armadan.Post{
			Title:  data.Title,
			Body:   data.Body,
			Author: data.Author,
		})

		if err != nil {
			return
		}

		posts, err := ps.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, views.ManagePosts(posts)).
			WithSuccess("Inlägg sparat").
			HTML()
	})
}

func UpdatePost(ps armadan.PostService, v armadan.Validator) http.Handler {
	type updatePostData struct {
		Title  string `json:"title" validate:"required"`
		Body   string `json:"body" validate:"required"`
		Author string `json:"author" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)
		if err != nil {
			return
		}

		data := updatePostData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		_, err = ps.Update(r.Context(), &armadan.Post{
			ID:     *id,
			Title:  data.Title,
			Body:   data.Body,
			Author: data.Author,
		})

		if err != nil {
			return
		}

		posts, err := ps.All(r.Context())
		if err != nil {
			return
		}

		response.
			New(w, r, partials.PostList(posts)).
			WithSuccess("Inlägg uppdaterat").
			HTML()
	})
}

func DeletePost(ps armadan.PostService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r)
		if err != nil {
			return
		}

		err = ps.Delete(r.Context(), *id)
		if err != nil {
			return
		}

		partials.SuccessToast("Inlägg borttaget").Render(r.Context(), w)
	})
}
