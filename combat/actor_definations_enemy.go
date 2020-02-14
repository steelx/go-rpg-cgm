package combat

import (
	"github.com/steelx/go-rpg-cgm/world"
)

var GoblinDef = ActorDef{
	Id:       "goblin",
	IsPlayer: false,
	Stats: world.BaseStats{
		HpNow:    90,
		HpMax:    90,
		MpNow:    0,
		MpMax:    0,
		Strength: 15, Speed: 8, Intelligence: 2,
	},
	Name:     "Goblin",
	Portrait: "../resources/avatar_hero.png",
	Actions:  []string{ActionAttack},
	Drop: Drop{
		XP:     150,
		Gold:   [2]int{5, 15},
		Always: nil,
		Chance: []DropChanceItem{
			{Oddment: 1, ItemId: -1},
			{Oddment: 3, ItemId: 11},
		},
	},
	StealItem: 14,
}
