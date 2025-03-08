package service

import (
	"context"
	"database/sql"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/patrickmn/go-cache"
)

const WEEKS_CACHE_KEY = "weeks:all"

func toWeek(ws any) *armadan.Week {
	switch w := ws.(type) {
	case schema.WeekDetail:
		return &armadan.Week{
			ID:         w.ID,
			Nr:         w.Nr,
			FinalsDate: armadan.ParseTime(w.FinalsDate.String),
			IsFinals:   w.IsFinals == 1,
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
			FinalsDate: armadan.ParseTime(w.FinalsDate.String),
			IsFinals:   w.IsFinals == 1,
			CourseID:   w.CourseID,
			TeeID:      w.TeeID,
			Dates:      armadan.GetWeekDates(int(w.Nr)),
		}
	}

	return &armadan.Week{}
}

type weeks struct {
	dbReader schema.Querier
	dbWriter schema.Querier
	cache    *cache.Cache
}

func NewWeekService(reader schema.Querier, writer schema.Querier, cache *cache.Cache) *weeks {
	return &weeks{
		dbReader: reader,
		dbWriter: writer,
		cache:    cache,
	}
}

func (s *weeks) All(ctx context.Context) ([]armadan.Week, error) {
	if cachedWeeks, found := s.cache.Get(WEEKS_CACHE_KEY); found {
		return cachedWeeks.([]armadan.Week), nil
	}

	weeks, err := s.dbReader.GetWeeks(ctx)

	if err != nil {
		return nil, err
	}

	mappedWeeks := armadan.MapEntities(weeks, toWeek)

	s.cache.Set(WEEKS_CACHE_KEY, mappedWeeks, cache.NoExpiration)

	return mappedWeeks, nil
}

func (s *weeks) Get(ctx context.Context, id string) (*armadan.Week, error) {
	week, err := s.dbReader.GetWeek(ctx, id)

	if err != nil {
		return nil, err
	}

	return toWeek(week), nil
}

func (s *weeks) Create(ctx context.Context, data *armadan.Week) (*armadan.Week, error) {
	week, err := s.dbWriter.CreateWeek(ctx, &schema.CreateWeekParams{
		ID:         armadan.GetId(),
		Nr:         data.Nr,
		IsFinals:   armadan.ToSqlBool(data.IsFinals),
		FinalsDate: sql.NullString{String: data.FinalsDate.Format(armadan.DEFAULT_TIME_FORMAT), Valid: true},
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
	week, err := s.dbWriter.UpdateWeek(ctx, &schema.UpdateWeekParams{
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

func (s *weeks) Delete(ctx context.Context, id string) error {
	if err := s.dbWriter.DeletePost(ctx, id); err != nil {
		return err
	}

	s.cache.Delete(WEEKS_CACHE_KEY)

	return nil
}
