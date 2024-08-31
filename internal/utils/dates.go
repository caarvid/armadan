package utils

import (
	"fmt"
	"time"
)

type WeekDates struct {
	start time.Time
	end   time.Time
}

func (wd *WeekDates) Format() string {
	return fmt.Sprintf("%s - %s", wd.start.Format("2/1"), wd.end.Format("2/1"))
}

func (wd *WeekDates) IsActive() bool {
	now := time.Now().YearDay()

	return now >= wd.start.YearDay() && now <= wd.end.YearDay()
}

func (wd *WeekDates) IsPreviousWeek() bool {
	return wd.end.YearDay() < time.Now().YearDay()
}

func getCurrentYear() int {
	return time.Now().Year()
}

func GetFirstOfJanuary() time.Time {
	return time.Date(getCurrentYear(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func GetWeekDates(nr int) *WeekDates {
	if nr <= 0 {
		return &WeekDates{}
	}

	firstOfJan := GetFirstOfJanuary()
	startWeekday := time.Tuesday
	startDate := firstOfJan.AddDate(0, 0, (nr-1)*7-int(firstOfJan.Weekday())).AddDate(0, 0, int(startWeekday))
	endDate := startDate.AddDate(0, 0, 5)

	return &WeekDates{
		start: startDate,
		end:   endDate,
	}
}
