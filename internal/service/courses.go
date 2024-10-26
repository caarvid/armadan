package service

import (
	"context"
	"encoding/json"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
)

const COURSES_CACHE_KEY = "courses:all"

func toCourse(data any) *armadan.Course {
	switch c := data.(type) {
	case schema.GetCourseRow:
		var holes []armadan.Hole
		var tees []armadan.Tee

		json.Unmarshal(c.Holes, &holes)
		json.Unmarshal(c.Tees, &tees)

		return &armadan.Course{
			ID:    c.ID,
			Par:   c.Par,
			Name:  c.Name,
			Holes: holes,
			Tees:  tees,
		}
	case schema.GetCoursesRow:
		var holes []armadan.Hole
		var tees []armadan.Tee

		json.Unmarshal(c.Holes, &holes)
		json.Unmarshal(c.Tees, &tees)

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
	db    schema.Querier
	pool  *pgxpool.Pool
	cache *cache.Cache
}

func NewCourseService(db schema.Querier, pool *pgxpool.Pool, cache *cache.Cache) *courses {
	return &courses{db: db, pool: pool, cache: cache}
}

func (cs *courses) All(ctx context.Context) ([]armadan.Course, error) {
	if cachedCourses, found := cs.cache.Get(COURSES_CACHE_KEY); found {
		return cachedCourses.([]armadan.Course), nil
	}

	courses, err := cs.db.GetCourses(ctx)

	if err != nil {
		return nil, err
	}

	mappedCourses := armadan.MapEntities(courses, toCourse)

	cs.cache.Set(COURSES_CACHE_KEY, mappedCourses, cache.NoExpiration)

	return mappedCourses, nil
}

func (cs *courses) Get(ctx context.Context, id uuid.UUID) (*armadan.Course, error) {
	course, err := cs.db.GetCourse(ctx, id)

	if err != nil {
		return nil, err
	}

	return toCourse(course), nil
}

func (cs *courses) GetTees(ctx context.Context, id uuid.UUID) ([]armadan.Tee, error) {
	tees, err := cs.db.GetTeesByCourse(ctx, id)

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
	tx, err := cs.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	qtx := schema.New(tx)

	var par int32

	for _, h := range data.Holes {
		par += h.Par
	}

	course, err := qtx.CreateCourse(ctx, &schema.CreateCourseParams{
		Name: data.Name,
		Par:  par,
	})

	if err != nil {
		return nil, err
	}

	var newHoles []*schema.CreateHolesParams

	for _, newHole := range data.Holes {
		newHoles = append(newHoles, &schema.CreateHolesParams{
			Nr:       newHole.Nr,
			Par:      newHole.Par,
			Index:    newHole.Index,
			CourseID: course.ID,
		})
	}

	if _, err := qtx.CreateHoles(ctx, newHoles); err != nil {
		return nil, err
	}

	if len(data.Tees) > 0 {
		var newTees []*schema.CreateTeesParams

		for _, newTee := range data.Tees {
			newTees = append(newTees, &schema.CreateTeesParams{
				Name:  newTee.Name,
				Slope: newTee.Slope,
				Cr: pgtype.Numeric{
					Int:              newTee.Cr.BigInt(),
					Exp:              newTee.Cr.Exponent(),
					NaN:              false,
					Valid:            true,
					InfinityModifier: pgtype.Finite,
				},
				CourseID: course.ID,
			})
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	// TODO: Do I really need to return the new course here?
	return nil, nil
}

func (cs *courses) Update(ctx context.Context, data *armadan.Course) (*armadan.Course, error) {
	tx, err := cs.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	qtx := schema.New(tx)

	var par int32

	for _, h := range data.Holes {
		par += h.Par
	}

	course, err := qtx.UpdateCourse(ctx, &schema.UpdateCourseParams{
		Name: data.Name,
		Par:  par,
		ID:   data.ID,
	})

	if err != nil {
		return nil, err
	}

	var holes []*schema.UpdateHolesParams

	for _, h := range data.Holes {
		holes = append(holes, &schema.UpdateHolesParams{
			Nr:    h.Nr,
			ID:    h.ID,
			Index: h.Index,
			Par:   h.Par,
		})
	}

	qtx.UpdateHoles(ctx, holes).Exec(nil)

	if len(data.Tees) > 0 {
		var tees []*schema.CreateTeesParams
		var updatedTees []*schema.UpdateTeesParams
		emptyId := uuid.UUID{}

		for _, t := range data.Tees {
			cr := pgtype.Numeric{
				Valid:            true,
				NaN:              false,
				InfinityModifier: pgtype.Finite,
				Int:              t.Cr.BigInt(),
				Exp:              t.Cr.Exponent(),
			}

			if t.ID.String() == emptyId.String() {
				tees = append(tees, &schema.CreateTeesParams{
					Name:     t.Name,
					Slope:    t.Slope,
					Cr:       cr,
					CourseID: course.ID,
				})
			} else {
				updatedTees = append(updatedTees, &schema.UpdateTeesParams{
					ID:    t.ID,
					Name:  t.Name,
					Slope: t.Slope,
					Cr:    cr,
				})
			}
		}

		if len(tees) > 0 {
			if _, err := qtx.CreateTees(ctx, tees); err != nil {
				return nil, err
			}
		}

		qtx.UpdateTees(ctx, updatedTees).Exec(nil)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	// TODO: Do I really need to return the course here?
	return nil, nil
}

func (cs *courses) Delete(ctx context.Context, id uuid.UUID) error {
	if err := cs.db.DeleteCourse(ctx, id); err != nil {
		return err
	}

	cs.db.DeleteCourse(ctx, id)
	cs.cache.Delete(COURSES_CACHE_KEY)

	return nil
}

func (cs *courses) DeleteTee(ctx context.Context, id uuid.UUID) error {
	if err := cs.db.DeleteTee(ctx, id); err != nil {
		return err
	}

	cs.cache.Delete(WEEKS_CACHE_KEY)
	cs.cache.Delete(COURSES_CACHE_KEY)

	return nil
}
