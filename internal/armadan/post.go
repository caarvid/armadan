package armadan

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PostService interface {
	All(context.Context) ([]Post, error)
	Get(context.Context, uuid.UUID) (*Post, error)
	Create(context.Context, *Post) (*Post, error)
	Update(context.Context, *Post) (*Post, error)
	Delete(context.Context, uuid.UUID) error
}

type Post struct {
	ID        uuid.UUID
	Title     string
	Body      string
	Author    string
	CreatedAt time.Time
}
