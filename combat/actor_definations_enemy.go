package combat

import (
	"github.com/steelx/go-rpg-cgm/dice"
	"github.com/steelx/go-rpg-cgm/world"
)

var GoblinDef = ActorDef{
	Id:       "goblin",
	IsPlayer: false,
	Stats: world.BaseStats{
		HpNow:    50,
		HpMax:    300,
		MpNow:    300,
		MpMax:    300,
		Strength: 10, Speed: 10, Intelligence: 10,
		Attack: 10,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("4d50+100"),
		"MpMax":        dice.Create("2d50+100"),
		"Strength":     world.StatsGrowth.Fast,
		"Speed":        world.StatsGrowth.Fast,
		"Intelligence": world.StatsGrowth.Med,
	},
	Name:             "Goblin",
	Portrait:         "../resources/avatar_hero.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}
