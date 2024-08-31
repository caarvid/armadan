package handlers

import (
	"github.com/caarvid/armadan/internal/schema"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
)

type Handler struct {
	db    *schema.Queries
	pool  *pgxpool.Pool
	cache *cache.Cache
}

type idParam struct {
	ID uuid.UUID `param:"id" validate:"required,uuid4"`
}

func Init(q *schema.Queries, p *pgxpool.Pool, c *cache.Cache) *Handler {
	return &Handler{
		db:    q,
		pool:  p,
		cache: c,
	}
}
