package dateutil

import (
	"math"
	"time"
)

type weekDates struct {
	Start time.Time
	End   time.Time
}

func getFirstOfJanuary(year int) time.Time {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
}

// Gets the start and end dates for a week (Tue - Sun)
// Set year = 0 to default to the current year
// Week number is clamped between 1 - 52
func GetWeekDates(year, nr int) weekDates {
	if year <= 0 {
		year = time.Now().Year()
	}

	nr = int(math.Max(1, math.Min(float64(nr), 52)))

	firstOfJan := getFirstOfJanuary(year)
	startWeekday := time.Tuesday
	startDate := firstOfJan.AddDate(0, 0, (nr-1)*7-int(firstOfJan.Weekday())).AddDate(0, 0, int(startWeekday))
	endDate := startDate.AddDate(0, 0, 5)

	return weekDates{
		Start: startDate,
		End:   endDate,
	}
}
