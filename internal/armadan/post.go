package armadan

import (
	"context"
	"time"
)

type PostService interface {
	All(context.Context) ([]Post, error)
	Get(context.Context, string) (*Post, error)
	Create(context.Context, *Post) (*Post, error)
	Update(context.Context, *Post) (*Post, error)
	Delete(context.Context, string) error
}

type Post struct {
	ID        string
	Title     string
	Body      string
	Author    string
	CreatedAt time.Time
}
