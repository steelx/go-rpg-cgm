package combat

import (
	"github.com/steelx/go-rpg-cgm/world"
)

var GoblinDef = ActorDef{
	Id:       "goblin",
	IsPlayer: false,
	Stats: world.BaseStats{
		HpNow:    150,
		HpMax:    150,
		MpNow:    300,
		MpMax:    300,
		Strength: 8, Speed: 8, Intelligence: 5,
		Attack: 10,
	},
	Name:     "Goblin",
	Portrait: "../resources/avatar_hero.png",
	Actions:  []string{attack},
}
