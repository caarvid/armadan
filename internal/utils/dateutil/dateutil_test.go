package dateutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWeekDates(t *testing.T) {
	year := 2025

	// 8/4 - 13/4
	week15 := GetWeekDates(year, 15)
	expectedStart := time.Date(year, time.April, 8, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, week15.Start, expectedStart)
	assert.Equal(t, week15.End, expectedStart.AddDate(0, 0, 5))

	// automatically clamped to week 52 (23/12 - 28/12)
	week52 := GetWeekDates(year, 100)
	expectedStart = time.Date(year, time.December, 23, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, week52.Start, expectedStart)
	assert.Equal(t, week52.End, expectedStart.AddDate(0, 0, 5))
}
