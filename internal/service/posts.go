package service

import (
	"context"
	"fmt"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/patrickmn/go-cache"
)

const POSTS_CACHE_KEY = "posts:all"

type posts struct {
	dbReader schema.Querier
	dbWriter schema.Querier
	cache    *cache.Cache
}

func toPost(post any) *armadan.Post {
	switch p := post.(type) {
	case schema.Post:
		return &armadan.Post{
			ID:        p.ID,
			Title:     p.Title,
			Body:      p.Body,
			Author:    p.Author,
			CreatedAt: armadan.ParseTime(p.CreatedAt),
		}
	}

	return &armadan.Post{}
}

func NewPostService(reader schema.Querier, writer schema.Querier, cache *cache.Cache) *posts {
	return &posts{
		dbReader: reader,
		dbWriter: writer,
		cache:    cache,
	}
}

func (s *posts) All(ctx context.Context) ([]armadan.Post, error) {
	fmt.Println("reading all posts")
	if cachedPosts, found := s.cache.Get(POSTS_CACHE_KEY); found {
		return cachedPosts.([]armadan.Post), nil
	}

	posts, err := s.dbReader.GetPosts(ctx)

	if err != nil {
		return nil, err
	}

	mappedPosts := armadan.MapEntities(posts, toPost)

	s.cache.Set(POSTS_CACHE_KEY, mappedPosts, cache.NoExpiration)

	return mappedPosts, nil
}

func (s *posts) Get(ctx context.Context, id string) (*armadan.Post, error) {
	post, err := s.dbReader.GetPost(ctx, id)

	if err != nil {
		return nil, err
	}

	return toPost(post), err
}

func (s *posts) Create(ctx context.Context, data *armadan.Post) (*armadan.Post, error) {
	post, err := s.dbWriter.CreatePost(ctx, &schema.CreatePostParams{
		ID:     armadan.GetId(),
		Title:  data.Title,
		Body:   data.Body,
		Author: data.Author,
	})

	if err != nil {
		return nil, err
	}

	s.cache.Delete(POSTS_CACHE_KEY)

	return toPost(post), nil
}

func (s *posts) Update(ctx context.Context, data *armadan.Post) (*armadan.Post, error) {
	post, err := s.dbWriter.UpdatePost(ctx, &schema.UpdatePostParams{
		ID:     data.ID,
		Title:  data.Title,
		Body:   data.Body,
		Author: data.Author,
	})

	if err != nil {
		return nil, err
	}

	s.cache.Delete(POSTS_CACHE_KEY)

	return toPost(post), nil
}

func (s *posts) Delete(ctx context.Context, id string) error {
	err := s.dbWriter.DeletePost(ctx, id)
	if err != nil {
		return err
	}

	s.cache.Delete(POSTS_CACHE_KEY)

	return nil
}
