package game_map

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/combat"
)

type Actors struct {
	Party   []*combat.Actor
	Enemies []*combat.Actor
}

type CombatCharacters struct {
	Party   []*Character
	Enemies []*Character
}

type CombatDef struct {
	Background string
	Actors     Actors
	Characters CombatCharacters
}

const (
	enemies = "enemies"
	party   = "party"
)

//All positions are in the range 0 - 1. Each number is a percentage of screen width and
//height offset from the center of the screen.
var combatLayout = map[string][][]pixel.Vec{
	party: {
		{
			pixel.V(0.25, -0.056),
		},
		{
			pixel.V(0.23, 0.024),
			pixel.V(0.27, -0.136),
		},
		{
			pixel.V(0.23, 0.024),
			pixel.V(0.25, -0.056),
			pixel.V(0.27, -0.136),
		},
	},
	enemies: {
		{
			pixel.V(-0.25, -0.056),
		},
		{
			pixel.V(-0.23, 0.024),
			pixel.V(-0.27, -0.136),
		},
		{
			pixel.V(-0.21, -0.056),
			pixel.V(-0.23, 0.024),
			pixel.V(-0.27, -0.136),
		},
		{
			pixel.V(-0.18, -0.056),
			pixel.V(-0.23, 0.056),
			pixel.V(-0.25, -0.056),
			pixel.V(-0.27, -0.168),
		},
		{
			pixel.V(-0.28, 0.032),
			pixel.V(-0.3, -0.056),
			pixel.V(-0.32, -0.144),
			pixel.V(-0.2, 0.004),
			pixel.V(-0.24, -0.116),
		},
		{
			pixel.V(-0.28, 0.032),
			pixel.V(-0.3, -0.056),
			pixel.V(-0.32, -0.144),
			pixel.V(-0.16, 0.032),
			pixel.V(-0.205, -0.056),
			pixel.V(-0.225, -0.144),
		},
	},
}
