package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/patrickmn/go-cache"
)

const COURSES_CACHE_KEY = "courses:all"

func toCourse(data any) *armadan.Course {
	switch c := data.(type) {
	case schema.CourseDetail:
		var holes []armadan.Hole
		var tees []armadan.Tee

		if c.Holes.Valid {
			json.Unmarshal([]byte(c.Holes.String), &holes)
		}

		if c.Tees.Valid {
			json.Unmarshal([]byte(c.Tees.String), &tees)
		}

		return &armadan.Course{
			ID:    c.ID,
			Par:   c.Par,
			Name:  c.Name,
			Holes: holes,
			Tees:  tees,
		}
	}

	return &armadan.Course{}
}

type courses struct {
	dbReader schema.Querier
	dbWriter schema.Querier
	pool     *sql.DB
	cache    *cache.Cache
}

func NewCourseService(reader, writer schema.Querier, pool *sql.DB, cache *cache.Cache) *courses {
	return &courses{
		dbReader: reader,
		dbWriter: writer,
		pool:     pool,
		cache:    cache,
	}
}

func (cs *courses) All(ctx context.Context) ([]armadan.Course, error) {
	if cachedCourses, found := cs.cache.Get(COURSES_CACHE_KEY); found {
		return cachedCourses.([]armadan.Course), nil
	}

	courses, err := cs.dbReader.GetCourses(ctx)

	if err != nil {
		return nil, err
	}

	mappedCourses := armadan.MapEntities(courses, toCourse)

	cs.cache.Set(COURSES_CACHE_KEY, mappedCourses, cache.NoExpiration)

	return mappedCourses, nil
}

func (cs *courses) Get(ctx context.Context, id string) (*armadan.Course, error) {
	course, err := cs.dbReader.GetCourse(ctx, id)

	if err != nil {
		return nil, err
	}

	return toCourse(course), nil
}

func (cs *courses) GetTees(ctx context.Context, id string) ([]armadan.Tee, error) {
	tees, err := cs.dbReader.GetTeesByCourse(ctx, id)

	if err != nil {
		return nil, err
	}

	return armadan.MapEntities(tees, func(a any) *armadan.Tee {
		switch t := a.(type) {
		case schema.Tee:
			return &armadan.Tee{
				ID:       t.ID,
				Name:     t.Name,
				Slope:    t.Slope,
				Cr:       t.Cr,
				CourseID: t.CourseID,
			}
		}

		return &armadan.Tee{}
	}), nil
}

func (cs *courses) Create(ctx context.Context, data *armadan.Course) (*armadan.Course, error) {
	tx, err := cs.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	var par int64
	for _, h := range data.Holes {
		par += h.Par
	}

	course, err := qtx.CreateCourse(ctx, &schema.CreateCourseParams{
		ID:   armadan.GetId(),
		Name: data.Name,
		Par:  par,
	})

	if err != nil {
		return nil, err
	}

	for _, newHole := range data.Holes {
		_, err = qtx.CreateHoles(ctx, &schema.CreateHolesParams{
			ID:          armadan.GetId(),
			Nr:          newHole.Nr,
			Par:         newHole.Par,
			StrokeIndex: newHole.Index,
			CourseID:    course.ID,
		})

		if err != nil {
			return nil, err
		}
	}

	for _, newTee := range data.Tees {
		_, err = qtx.CreateTees(ctx, &schema.CreateTeesParams{
			ID:       armadan.GetId(),
			Name:     newTee.Name,
			Slope:    newTee.Slope,
			Cr:       newTee.Cr,
			CourseID: course.ID,
		})

		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	return nil, nil
}

func (cs *courses) Update(ctx context.Context, data *armadan.Course) (*armadan.Course, error) {
	tx, err := cs.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	qtx := schema.New(tx)

	var par int64

	for _, h := range data.Holes {
		par += h.Par
	}

	course, err := qtx.UpdateCourse(ctx, &schema.UpdateCourseParams{
		ID:   data.ID,
		Name: data.Name,
		Par:  par,
	})

	if err != nil {
		return nil, err
	}

	for _, h := range data.Holes {
		err = qtx.UpdateHoles(ctx, &schema.UpdateHolesParams{
			Nr:          h.Nr,
			ID:          h.ID,
			StrokeIndex: h.Index,
			Par:         h.Par,
		})

		if err != nil {
			return nil, err
		}
	}

	for _, t := range data.Tees {
		if len(t.ID) == 0 {
			_, err = qtx.CreateTees(ctx, &schema.CreateTeesParams{
				ID:       armadan.GetId(),
				Name:     t.Name,
				Slope:    t.Slope,
				Cr:       t.Cr,
				CourseID: course.ID,
			})

			if err != nil {
				return nil, err
			}
		} else {
			err = qtx.UpdateTees(ctx, &schema.UpdateTeesParams{
				ID:    t.ID,
				Name:  t.Name,
				Slope: t.Slope,
				Cr:    t.Cr,
			})

			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	// TODO: Do I really need to return the course here?
	return nil, nil
}

func (cs *courses) Delete(ctx context.Context, id string) error {
	if err := cs.dbWriter.DeleteCourse(ctx, id); err != nil {
		return err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	return nil
}

func (cs *courses) DeleteTee(ctx context.Context, id string) error {
	if err := cs.dbWriter.DeleteTee(ctx, id); err != nil {
		return err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	return nil
}
