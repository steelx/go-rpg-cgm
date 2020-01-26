package combat

import (
	"github.com/steelx/go-rpg-cgm/dice"
)

const (
	attack = "attack"
	item   = "item"
)

var HeroDef = ActorDef{
	Id: "hero",
	Stats: BaseStats{
		HpNow:    300,
		HpMax:    300,
		MpNow:    300,
		MpMax:    300,
		Strength: 10, Speed: 10, Intelligence: 10,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("4d50+100"),
		"MpMax":        dice.Create("2d50+100"),
		"Strength":     StatsGrowth.Fast,
		"Speed":        StatsGrowth.Fast,
		"Intelligence": StatsGrowth.Med,
	},
	Name:             "Chandragupta",
	Portrait:         "../resources/avatar_hero.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{3, 1, 2}, //ItemsDB ID
}

var MageDef = ActorDef{
	Id: "mage",
	Stats: BaseStats{
		HpNow:    200,
		HpMax:    200,
		MpNow:    280,
		MpMax:    280,
		Strength: 8, Speed: 10, Intelligence: 20,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("3d40+100"),
		"MpMax":        dice.Create("4d50+100"),
		"Strength":     StatsGrowth.Med,
		"Speed":        StatsGrowth.Med,
		"Intelligence": StatsGrowth.Fast,
	},
	Name:             "Mrignayani",
	Portrait:         "../resources/avatar_mage.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{1, 2, 3}, //ItemsDB ID
}

var ThiefDef = ActorDef{
	Id: "thief",
	Stats: BaseStats{
		HpNow:    280,
		HpMax:    280,
		MpNow:    150,
		MpMax:    150,
		Strength: 10, Speed: 15, Intelligence: 10,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("3d40+100"),
		"MpMax":        dice.Create("4d50+100"),
		"Strength":     StatsGrowth.Med,
		"Speed":        StatsGrowth.Med,
		"Intelligence": StatsGrowth.Med,
	},
	Name:             "Shashank",
	Portrait:         "../resources/avatar_thief.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{1, 2, 3}, //ItemsDB ID
}
