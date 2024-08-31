package hcp

import "math"

func getHcpGroup(hcp float64) float64 {
	switch {
	case hcp <= 4.4:
		return 0.1
	case hcp <= 11.4:
		return 0.2
	case hcp <= 18.4:
		return 0.3
	case hcp <= 26.4:
		return 0.4
	default:
		return 0.5
	}
}

func GetNewHcp(hcp float64, par, total int32) float64 {
	hcpGroup := getHcpGroup(hcp)

	if total < par {
		return math.Round((hcp-float64(par-total)*hcpGroup)*10) / 10
	}

	if total > (par + int32(hcpGroup*10)) {
		return math.Round((hcp+0.1)*10) / 10
	}

	return hcp
}

func GetStrokes(hcp, cr float64, slope, par int) int {
	return int(math.Min(math.Round(hcp*float64(slope/113)+(cr-float64(par))), 18.0))
}
