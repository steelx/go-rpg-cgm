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
		HpNow:    40,
		HpMax:    40,
		MpNow:    8,
		MpMax:    8,
		Level:    0,
		Strength: 10, Speed: 12, Intelligence: 10,
		Counter: 0, //make it 0 to stop counter attacks by this Actor
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("2d25+25"),
		"MpMax":        dice.Create("1d5+2"),
		"Strength":     world.StatsGrowth.Fast,
		"Speed":        world.StatsGrowth.Fast,
		"Intelligence": world.StatsGrowth.Med,
	},
	ActionGrowth: map[int]map[string][]string{
		5: {
			ActionSpecial: []string{world.SpecialSlash},
		},
	},
	Name:             "Chandragupta",
	Portrait:         "../resources/avatar_hero.png",
	Actions:          []string{ActionAttack, ActionItem, ActionFlee}, //removed ActionSpecial -> unlocks later
	Special:          []string{world.SpecialSlash},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}

var MageDef = ActorDef{
	Id:       "mage",
	IsPlayer: true,
	Stats: world.BaseStats{
		HpNow:    30,
		HpMax:    30,
		MpNow:    10,
		MpMax:    10,
		Level:    0,
		Strength: 8, Speed: 11, Intelligence: 20,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("2d25+18"),
		"MpMax":        dice.Create("1d5+2"),
		"Strength":     world.StatsGrowth.Med,
		"Speed":        world.StatsGrowth.Med,
		"Intelligence": world.StatsGrowth.Fast,
	},
	ActionGrowth: map[int]map[string][]string{
		1: {
			ActionMagic: []string{world.SpellBolt},
		},
		2: {
			ActionMagic: []string{world.SpellFire, world.SpellIce},
		},
		4: {
			ActionMagic: []string{world.SpellBurn},
		},
	},
	Name:             "Mrignayani",
	Portrait:         "../resources/avatar_mage.png",
	Actions:          []string{ActionAttack, ActionItem, ActionFlee}, //ActionMagic
	Magic:            []string{world.SpellFire, world.SpellBurn, world.SpellBolt},
	ActiveEquipSlots: []int{0, 1, 2, 3}, //mage if no attack slot, Access goes to Attack slot(fix pending)
}

var ThiefDef = ActorDef{
	Id:       "thief",
	IsPlayer: true,
	Stats: world.BaseStats{
		HpNow:    35,
		HpMax:    35,
		MpNow:    7,
		MpMax:    7,
		Level:    0,
		Strength: 10, Speed: 13, Intelligence: 10,
	},
	StatGrowth: map[string]func() int{
		"HpMax":        dice.Create("2d25+20"),
		"MpMax":        dice.Create("1d10+5"),
		"Strength":     world.StatsGrowth.Med,
		"Speed":        world.StatsGrowth.Med,
		"Intelligence": world.StatsGrowth.Med,
	},
	ActionGrowth: map[int]map[string][]string{
		2: {
			world.SpecialSteal: []string{ActionSpecial},
		},
	},
	Name:             "Shashank",
	Portrait:         "../resources/avatar_thief.png",
	Actions:          []string{ActionAttack, ActionItem, ActionFlee},
	Special:          []string{world.SpecialSteal},
	ActiveEquipSlots: []int{0, 1, 2, 3},
}
