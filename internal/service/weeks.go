package service

import (
	"context"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickmn/go-cache"
)

const WEEKS_CACHE_KEY = "weeks:all"

func toWeek(ws any) *armadan.Week {
	switch w := ws.(type) {
	case schema.WeekDetail:
		return &armadan.Week{
			ID:         w.ID,
			Nr:         w.Nr,
			FinalsDate: w.FinalsDate.Time,
			IsFinals:   w.IsFinals.Bool,
			CourseID:   w.CourseID,
			CourseName: w.CourseName,
			TeeID:      w.TeeID,
			TeeName:    w.TeeName,
			Dates:      armadan.GetWeekDates(int(w.Nr)),
		}
	case schema.Week:
		return &armadan.Week{
			ID:         w.ID,
			Nr:         w.Nr,
			FinalsDate: w.FinalsDate.Time,
			IsFinals:   w.IsFinals.Bool,
			CourseID:   w.CourseID,
			TeeID:      w.TeeID,
			Dates:      armadan.GetWeekDates(int(w.Nr)),
		}
	}

	return &armadan.Week{}
}

type weeks struct {
	db    schema.Querier
	cache *cache.Cache
}

func NewWeekService(db schema.Querier, cache *cache.Cache) *weeks {
	return &weeks{
		db:    db,
		cache: cache,
	}
}

func (s *weeks) All(ctx context.Context) ([]armadan.Week, error) {
	if cachedWeeks, found := s.cache.Get(WEEKS_CACHE_KEY); found {
		return cachedWeeks.([]armadan.Week), nil
	}

	weeks, err := s.db.GetWeeks(ctx)

	if err != nil {
		return nil, err
	}

	mappedWeeks := armadan.MapEntities(weeks, toWeek)

	s.cache.Set(WEEKS_CACHE_KEY, mappedWeeks, cache.NoExpiration)

	return mappedWeeks, nil
}

func (s *weeks) Get(ctx context.Context, id uuid.UUID) (*armadan.Week, error) {
	week, err := s.db.GetWeek(ctx, id)

	if err != nil {
		return nil, err
	}

	return toWeek(week), nil
}

func (s *weeks) Create(ctx context.Context, data *armadan.Week) (*armadan.Week, error) {
	week, err := s.db.CreateWeek(ctx, &schema.CreateWeekParams{
		Nr:         data.Nr,
		IsFinals:   pgtype.Bool{Bool: data.IsFinals, Valid: true},
		FinalsDate: pgtype.Timestamptz{Time: data.FinalsDate, Valid: true},
		CourseID:   data.CourseID,
		TeeID:      data.TeeID,
	})

	if err != nil {
		return nil, err
	}

	s.cache.Delete(WEEKS_CACHE_KEY)

	return toWeek(week), nil
}

func (s *weeks) Update(ctx context.Context, data *armadan.Week) (*armadan.Week, error) {
	week, err := s.db.UpdateWeek(ctx, &schema.UpdateWeekParams{
		ID:       data.ID,
		Nr:       data.Nr,
		CourseID: data.CourseID,
		TeeID:    data.TeeID,
	})

	if err != nil {
		return nil, err
	}

	s.cache.Delete(WEEKS_CACHE_KEY)

	return toWeek(week), nil
}

func (s *weeks) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.db.DeletePost(ctx, id); err != nil {
		return err
	}

	s.cache.Delete(WEEKS_CACHE_KEY)

	return nil
}
