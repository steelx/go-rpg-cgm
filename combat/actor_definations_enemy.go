package combat

import (
	"github.com/steelx/go-rpg-cgm/world"
)

var EnemyDefinitions = map[string]ActorDef{
	"goblin": GoblinDef,
	"dragon": DragonDef,
	"ogre":   OgreDef,
}

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
	Portrait: "../resources/avatar_hero.png", //temp we need this at Actor Create
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

var DragonDef = ActorDef{
	Id:       "dragon",
	IsPlayer: false,
	Stats: world.BaseStats{
		HpNow:    200,
		HpMax:    200,
		MpNow:    0,
		MpMax:    0,
		Strength: 35, Speed: 8, Intelligence: 20,
		Counter: 0.1,
	},
	Name:     "Green Dragon",
	Portrait: "../resources/avatar_hero.png", //temp we need this at Actor Create
	Actions:  []string{ActionAttack},
	Drop: Drop{
		XP:     350,
		Gold:   [2]int{250, 300},
		Always: nil,
		Chance: []DropChanceItem{
			{Oddment: 1, ItemId: -1},
			{Oddment: 3, ItemId: 10},
		},
	},
	StealItem: 11,
}

var OgreDef = ActorDef{
	Id:       "ogre",
	IsPlayer: false,
	Stats: world.BaseStats{
		HpNow:    150,
		HpMax:    150,
		MpNow:    0,
		MpMax:    0,
		Strength: 20, Speed: 8, Intelligence: 2,
		Counter: 0,
	},
	Name:     "Ogre",
	Portrait: "../resources/avatar_hero.png", //temp we need this at Actor Create
	Actions:  []string{ActionAttack},
	Drop: Drop{
		XP:     250,
		Gold:   [2]int{100, 200},
		Always: nil,
		Chance: []DropChanceItem{
			{Oddment: 1, ItemId: -1},
			{Oddment: 3, ItemId: 10},
		},
	},
	StealItem: 12,
}
