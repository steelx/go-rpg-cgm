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
		MpNow:    50,
		MpMax:    50,
		Strength: 8, Speed: 8, Intelligence: 5,
		Attack: 10,
	},
	Name:     "Goblin",
	Portrait: "../resources/avatar_hero.png",
	Actions:  []string{ActionAttack},
	Drop: Drop{
		XP:     5,
		Gold:   [2]int{0, 5},
		Always: nil,
		Chance: []DropChanceItem{
			{Oddment: 95, ItemId: -1},
			{Oddment: 3, ItemId: 11},
		},
	},
}
