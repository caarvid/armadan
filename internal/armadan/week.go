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

type WeekDates struct {
	Start time.Time
	End   time.Time
}

func (wd *WeekDates) String() string {
	return fmt.Sprintf("%s - %s", wd.Start.Format("2/1"), wd.End.Format("2/1"))
}

func getFirstOfJanuary() time.Time {
	return time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func GetWeekDates(nr int) WeekDates {
	if nr <= 0 {
		return WeekDates{}
	}

	firstOfJan := getFirstOfJanuary()
	startWeekday := time.Tuesday
	startDate := firstOfJan.AddDate(0, 0, (nr-1)*7-int(firstOfJan.Weekday())).AddDate(0, 0, int(startWeekday))
	endDate := startDate.AddDate(0, 0, 5)

	return WeekDates{
		Start: startDate,
		End:   endDate,
	}
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
