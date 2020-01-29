package combat

import (
	"github.com/steelx/go-rpg-cgm/dice"
	"github.com/steelx/go-rpg-cgm/world"
)

const (
	attack = "attack"
	item   = "item"
)

var PartyMembersDefinitions = map[string]ActorDef{
	"hero":  HeroDef,
	"mage":  MageDef,
	"thief": ThiefDef,
}

var HeroDef = ActorDef{
	Id: "hero",
	Stats: world.BaseStats{
		HpNow:    300,
		HpMax:    300,
		MpNow:    300,
		MpMax:    300,
		Strength: 10, Speed: 10, Intelligence: 10,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("4d50+100"),
		"MpMax":        dice.Create("2d50+100"),
		"Strength":     world.StatsGrowth.Fast,
		"Speed":        world.StatsGrowth.Fast,
		"Intelligence": world.StatsGrowth.Med,
	},
	Name:             "Chandragupta",
	Portrait:         "../resources/avatar_hero.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}

var MageDef = ActorDef{
	Id: "mage",
	Stats: world.BaseStats{
		HpNow:    200,
		HpMax:    200,
		MpNow:    280,
		MpMax:    280,
		Strength: 8, Speed: 10, Intelligence: 20,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("3d40+100"),
		"MpMax":        dice.Create("4d50+100"),
		"Strength":     world.StatsGrowth.Med,
		"Speed":        world.StatsGrowth.Med,
		"Intelligence": world.StatsGrowth.Fast,
	},
	Name:             "Mrignayani",
	Portrait:         "../resources/avatar_mage.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{1, 2, 3}, //mage dont get Attack slot
}

var ThiefDef = ActorDef{
	Id: "thief",
	Stats: world.BaseStats{
		HpNow:    280,
		HpMax:    280,
		MpNow:    150,
		MpMax:    150,
		Strength: 10, Speed: 15, Intelligence: 10,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("3d40+100"),
		"MpMax":        dice.Create("4d50+100"),
		"Strength":     world.StatsGrowth.Med,
		"Speed":        world.StatsGrowth.Med,
		"Intelligence": world.StatsGrowth.Med,
	},
	Name:             "Shashank",
	Portrait:         "../resources/avatar_thief.png",
	Actions:          []string{attack, item},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}
