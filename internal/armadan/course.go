package armadan

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CourseService interface {
	All(context.Context) ([]Course, error)
	Get(context.Context, uuid.UUID) (*Course, error)
	GetTees(context.Context, uuid.UUID) ([]Tee, error)
	Create(context.Context, *Course) (*Course, error)
	Update(context.Context, *Course) (*Course, error)
	Delete(context.Context, uuid.UUID) error
	DeleteTee(context.Context, uuid.UUID) error
}

type Hole struct {
	ID       uuid.UUID
	CourseID uuid.UUID `json:"course_id"`
	Nr       int32
	Par      int32
	Index    int32
}

type Tee struct {
	ID       uuid.UUID
	CourseID uuid.UUID `json:"course_id"`
	Name     string
	Slope    int32
	Cr       decimal.Decimal
}

type Course struct {
	ID    uuid.UUID
	Name  string
	Par   int32
	Holes []Hole
	Tees  []Tee
}
