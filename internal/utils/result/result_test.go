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

func createScores() []armadan.Score {
	scores := make([]armadan.Score, 18)

	for i := range 9 {
		scores[i].Strokes = 4
		scores[i].Par = 4
		scores[i].Index = int64((i * 2) + 1)
	}

	for i := range 9 {
		scores[9+i].Strokes = 4
		scores[9+i].Par = 4
		scores[9+i].Index = int64((i + 1) * 2)
	}

	return scores
}

func TestGetRoundSummary(t *testing.T) {
	scores := createScores()

	result := GetRoundSummary(scores, 0)
	assert.EqualValues(t, 36, result.NetOut)
	assert.EqualValues(t, 36, result.NetIn)

	result = GetRoundSummary(scores, -3)
	assert.EqualValues(t, 37, result.NetOut)
	assert.EqualValues(t, 38, result.NetIn)

	result = GetRoundSummary(scores, 5)
	assert.EqualValues(t, 33, result.NetOut)
	assert.EqualValues(t, 34, result.NetIn)
}
