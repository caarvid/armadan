package armadan

import (
	"context"
)

type CourseService interface {
	All(context.Context) ([]Course, error)
	Get(context.Context, string) (*Course, error)
	GetTees(context.Context, string) ([]Tee, error)
	Create(context.Context, *Course) (*Course, error)
	Update(context.Context, *Course) (*Course, error)
	Delete(context.Context, string) error
	DeleteTee(context.Context, string) error
}

type Hole struct {
	ID       string
	CourseID string `json:"course_id"`
	Nr       int64
	Par      int64
	Index    int64
}

type Tee struct {
	ID       string
	CourseID string `json:"course_id"`
	Name     string
	Slope    int64
	Cr       float64
}

type Course struct {
	ID    string
	Name  string
	Par   int64
	Holes []Hole
	Tees  []Tee
}
