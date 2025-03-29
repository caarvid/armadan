package armadan

import (
	"context"
	"iter"
	"slices"
	"strings"
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

func (c *Course) teeIter() iter.Seq[Tee] {
	return func(yield func(Tee) bool) {
		for _, v := range c.Tees {
			if !yield(v) {
				return
			}
		}
	}
}

func (c *Course) TeeList() string {
	names := []string{}
	sortedTees := slices.SortedFunc(c.teeIter(), func(a Tee, b Tee) int {
		return int(a.Slope - b.Slope)
	})

	for _, tee := range sortedTees {
		names = append(names, tee.Name)
	}

	return strings.Join(names, ", ")
}

type parInfo struct {
	In  int64
	Out int64
}

func (c *Course) ParInfo() parInfo {
	var in, out int64

	for _, h := range c.Holes[:9] {
		out = out + h.Par
	}

	for _, h := range c.Holes[9:] {
		in = in + h.Par
	}

	return parInfo{
		In:  in,
		Out: out,
	}
}
