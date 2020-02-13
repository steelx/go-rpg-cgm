package combat

import (
	"github.com/steelx/go-rpg-cgm/dice"
	"github.com/steelx/go-rpg-cgm/world"
)

const (
	ActionAttack  = "Attack"
	ActionItem    = "Item"
	ActionMagic   = "Magic"
	ActionSpecial = "Special"
	ActionFlee    = "Flee"
)

var PartyMembersDefinitions = map[string]ActorDef{
	"hero":  HeroDef,
	"mage":  MageDef,
	"thief": ThiefDef,
}

var HeroDef = ActorDef{
	Id:       "hero",
	IsPlayer: true,
	Stats: world.BaseStats{
		HpNow:    300,
		HpMax:    300,
		MpNow:    200,
		MpMax:    200,
		Strength: 10, Speed: 10, Intelligence: 10,
		Attack:  10,
		Counter: 0, //make it 0 to stop counter attacks by this Actor
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
	Actions:          []string{ActionAttack, ActionSpecial, ActionItem, ActionFlee},
	Special:          []string{world.SpecialSlash},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}

var MageDef = ActorDef{
	Id:       "mage",
	IsPlayer: true,
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
	Actions:          []string{ActionAttack, ActionMagic, ActionItem, ActionFlee},
	Magic:            []string{world.SpellFire, world.SpellBurn, world.SpellBolt},
	ActiveEquipSlots: []int{0, 1, 2, 3}, //mage if no attack slot, Access goes to Attack slot(fix pending)
}

var ThiefDef = ActorDef{
	Id:       "thief",
	IsPlayer: true,
	Stats: world.BaseStats{
		HpNow:    200,
		HpMax:    200,
		MpNow:    150,
		MpMax:    150,
		Strength: 10, Speed: 11, Intelligence: 10,
		Attack: 8,
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
	Actions:          []string{ActionAttack, ActionSpecial, ActionItem, ActionFlee},
	Special:          []string{world.SpecialSteal},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}
