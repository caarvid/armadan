package utils_test

import (
	"testing"

	"github.com/caarvid/armadan/internal/utils"
)

func TestGetWeekDates(t *testing.T) {
	t.Log(utils.GetWeekDates(10).Format())
}
