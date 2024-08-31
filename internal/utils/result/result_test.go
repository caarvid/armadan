package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestPlayers(scores []int) []Player {
	p := make([]Player, len(scores))

	for i, s := range scores {
		p[i].Id = fmt.Sprint(i)
		p[i].Score = s
	}

	return p
}

func TestGetWinners(t *testing.T) {
	// 2 winners (2 max)
	winners := GetWinners(createTestPlayers([]int{70, 70, 72, 75, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80}))
	expected := []Winner{{Id: "0", Points: 375}, {Id: "1", Points: 375}}

	assert.Equal(t, expected, winners)

	// 1 winner, 2 second place (2 max)
	winners = GetWinners(createTestPlayers([]int{70, 71, 71, 75, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80}))
	expected = []Winner{{Id: "0", Points: 420}, {Id: "1", Points: 140}, {Id: "2", Points: 140}}

	assert.Equal(t, expected, winners)
}
