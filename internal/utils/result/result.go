package result

import (
	"math"
	"slices"
)

func getPricePool(players int) int {
	return 50 * players
}

func getNrOfWinners(players int) int {
	switch {
	case players <= 10:
		return 1
	case players <= 20:
		return 2
	case players <= 30:
		return 3
	default:
		return 4
	}
}

type Player struct {
	Id    string
	Score int
}

type Winner struct {
	Id     string
	Points int
}

func groupByPlacement(players []Player) ([]Player, []Player, []Player, []Player) {
	m := make(map[int][]Player)

	for _, p := range players {
		m[p.Score] = append(m[p.Score], p)
	}

	res := make([][]Player, 0, len(m))

	for _, v := range m {
		res = append(res, v)
	}

	if len(res) < 4 {
		l := 4 - len(res)

		for i := 0; i < l; i++ {
			res = append(res, make([]Player, 0))
		}
	}

	slices.SortFunc(res, func(a, b []Player) int {
		if len(a) == 0 || len(b) == 0 {
			return len(b) - len(a)
		}

		return a[0].Score - b[0].Score
	})

	return res[0], res[1], res[2], res[3]
}

func createWinners(players []Player, points float64) []Winner {
	winners := make([]Winner, len(players))

	for i, p := range players {
		winners[i].Id = p.Id
		winners[i].Points = int(points)
	}

	return winners
}

func GetWinners(players []Player) []Winner {
	nrOfWinners := getNrOfWinners(len(players))
	pool := float64(getPricePool(len(players)))
	first, second, third, fourth := groupByPlacement(players)

	if nrOfWinners == 1 || len(first) >= nrOfWinners {
		return createWinners(first, math.Round(pool/float64(len(first))))
	} else if nrOfWinners == 2 {
		return slices.Concat(
			createWinners(first, math.Round(pool*0.6)),
			createWinners(second, math.Round(pool*0.4/float64(len(second)))),
		)
	} else if nrOfWinners == 3 {
		if len(first) == 2 {
			return slices.Concat(
				createWinners(first, math.Round(pool*0.8/float64(len(first)))),
				createWinners(second, math.Round(pool*0.2/float64(len(second)))),
			)
		} else if len(first) == 1 && len(second) >= 2 {
			return slices.Concat(
				createWinners(first, math.Round(pool*0.5/float64(len(first)))),
				createWinners(second, math.Round(pool*0.5/float64(len(second)))),
			)
		} else {
			return slices.Concat(
				createWinners(first, math.Round(pool*0.5)),
				createWinners(second, math.Round(pool*0.3)),
				createWinners(third, math.Round(pool*0.2/float64(len(third)))),
			)
		}
	} else {
		if len(first) == 3 {
			return slices.Concat(
				createWinners(first, math.Round(pool*0.9/float64(len(first)))),
				createWinners(second, math.Round(pool*0.1/float64(len(second)))),
			)
		} else if len(first) == 2 {
			if len(second) >= 2 {
				return slices.Concat(
					createWinners(first, math.Round(pool*0.75/float64(len(first)))),
					createWinners(second, math.Round(pool*0.25/float64(len(second)))),
				)
			} else if len(second) == 1 {
				return slices.Concat(
					createWinners(first, math.Round(pool*0.75/float64(len(first)))),
					createWinners(second, math.Round(pool*0.15)),
					createWinners(third, math.Round(pool*0.1/float64(len(third)))),
				)
			}
		} else if len(first) == 1 {
			if len(second) >= 3 {
				return slices.Concat(
					createWinners(first, math.Round(pool*0.5)),
					createWinners(second, math.Round(pool*0.5/float64(len(second)))),
				)
			} else if len(second) == 2 {
				return slices.Concat(
					createWinners(first, math.Round(pool*0.5)),
					createWinners(second, math.Round(pool*0.4/float64(len(second)))),
					createWinners(third, math.Round(pool*0.1/float64(len(third)))),
				)
			} else if len(second) == 1 && len(third) >= 2 {
				return slices.Concat(
					createWinners(first, math.Round(pool*0.5)),
					createWinners(second, math.Round(pool*0.25)),
					createWinners(third, math.Round(pool*0.25/float64(len(third)))),
				)
			}
		} else {
			return slices.Concat(
				createWinners(first, math.Round(pool*0.5)),
				createWinners(second, math.Round(pool*0.25)),
				createWinners(third, math.Round(pool*0.15)),
				createWinners(fourth, math.Round(pool*0.1/float64(len(fourth)))),
			)
		}
	}

	return nil
}
