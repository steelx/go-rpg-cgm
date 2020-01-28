package world

import "github.com/steelx/go-rpg-cgm/dice"

type StatsGrowthT struct {
	Fast func() int
	Med  func() int
	Slow func() int
}

var StatsGrowth = StatsGrowthT{
	Fast: dice.Create("3d2"),
	Med:  dice.Create("1d3"),
	Slow: dice.Create("1d2"),
}
