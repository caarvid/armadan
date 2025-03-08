package course

import (
	"github.com/caarvid/armadan/internal/armadan"
)

type parInfo struct {
	In  int64
	Out int64
}

func GetParInfo(c *armadan.Course) parInfo {
	var in, out int64

	for _, h := range c.Holes[:9] {
		out = out + h.Par
	}

	for _, h := range c.Holes[9:] {
		in = in + h.Par
	}

	return parInfo{
		In:  in,
		Out: out,
	}
}
