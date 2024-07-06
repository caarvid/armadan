package utils

import (
	"slices"
	"strings"

	"github.com/caarvid/armadan/internal/dto"
)

func TeeNameList(tees *dto.TeeList) string {
	names := []string{}

	slices.SortFunc(*tees, func(a dto.Tee, b dto.Tee) int {
		return int(a.Slope - b.Slope)
	})

	for _, tee := range *tees {
		names = append(names, tee.Name)
	}

	return strings.Join(names, ", ")
}
