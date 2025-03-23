package result

import (
	"fmt"
	"testing"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/stretchr/testify/assert"
)

func createRounds(scores []int64) []armadan.Round {
	p := make([]armadan.Round, len(scores))

	for i, s := range scores {
		p[i].PlayerID = fmt.Sprint(i)
		p[i].NetTotal = s
	}

	return p
}

func TestGetWinners(t *testing.T) {
	// 2 winners (2 max)
	winners := GetWinners(createRounds([]int64{70, 70, 72, 75, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80}))
	expected := []armadan.Winner{{PlayerID: "0", Points: 375}, {PlayerID: "1", Points: 375}}

	assert.Equal(t, 2, len(winners))

	for i, w := range winners {
		assert.Equal(t, expected[i].PlayerID, w.PlayerID)
		assert.Equal(t, expected[i].Points, w.Points)
	}

	// 1 winner, 2 second place (2 max)
	winners = GetWinners(createRounds([]int64{70, 71, 71, 75, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80}))
	expected = []armadan.Winner{{PlayerID: "0", Points: 420}, {PlayerID: "1", Points: 140}, {PlayerID: "2", Points: 140}}

	assert.Equal(t, 3, len(winners))

	for i, w := range winners {
		assert.Equal(t, expected[i].PlayerID, w.PlayerID)
		assert.Equal(t, expected[i].Points, w.Points)
	}
}
