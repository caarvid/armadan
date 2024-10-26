package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWeekDates(t *testing.T) {
	v20 := GetWeekDates(20)
	v30 := GetWeekDates(30)

	assert.Equal(t, v20.start, time.Date(2024, time.May, 14, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, v20.end, time.Date(2024, time.May, 19, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, v30.start, time.Date(2024, time.July, 23, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, v30.end, time.Date(2024, time.July, 28, 0, 0, 0, 0, time.UTC))
}
