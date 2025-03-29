package armadan

import (
	"context"
	"fmt"
	"time"
)

type WeekService interface {
	All(context.Context) ([]Week, error)
	Get(context.Context, string) (*Week, error)
	Create(context.Context, *Week) (*Week, error)
	Update(context.Context, *Week) (*Week, error)
	Delete(context.Context, string) error
}

type Week struct {
	ID         string
	Nr         int64
	FinalsDate time.Time
	IsFinals   bool
	CourseID   string
	CourseName string
	TeeID      string
	TeeName    string
	StartDate  time.Time
	EndDate    time.Time
}

func (w *Week) FormattedDate() string {
	return fmt.Sprintf("%s - %s", w.StartDate.Format("2/1"), w.EndDate.Format("2/1"))
}

func (w *Week) IsCurrent() bool {
	now := time.Now().YearDay()

	return now >= w.StartDate.YearDay() && now <= w.EndDate.YearDay()
}

func (w *Week) IsPrevious() bool {
	return w.EndDate.YearDay() < time.Now().YearDay()
}
