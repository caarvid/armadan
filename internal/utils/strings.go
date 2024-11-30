package utils

import (
	"slices"
	"strings"

	"github.com/caarvid/armadan/internal/armadan"
)

func TeeNameList(tees []armadan.Tee) string {
	names := []string{}

	slices.SortFunc(tees, func(a armadan.Tee, b armadan.Tee) int {
		return int(a.Slope - b.Slope)
	})

	for _, tee := range tees {
		names = append(names, tee.Name)
	}

	return strings.Join(names, ", ")
}
