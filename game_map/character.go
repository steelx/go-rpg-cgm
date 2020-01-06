package game_map

import (
	"github.com/steelx/go-rpg-cgm/state_machine"
)

type CharacterStateBase struct {
	Character  *Character
	Map        *GameMap
	Entity     *Entity
	Controller *state_machine.StateMachine
}

type Character struct {
	Name                                  string
	AnimUp, AnimRight, AnimDown, AnimLeft []int
	Facing                                string
	Entity                                *Entity
	Controller                            *state_machine.StateMachine //[Name] -> [function that returns state]
}

func (ch Character) GetFacedTileCoords() (x, y float64) {
	var xOff, yOff float64 = 0, 0
	if ch.Facing == CharacterFacingDirection[3] {
		xOff = -1 //"left"
	} else if ch.Facing == CharacterFacingDirection[1] {
		xOff = 1 //"right"
	} else if ch.Facing == CharacterFacingDirection[0] {
		yOff = -1 //"up"
	} else if ch.Facing == CharacterFacingDirection[2] {
		yOff = 1 //"down"
	}

	x = ch.Entity.TileX + xOff
	y = ch.Entity.TileY + yOff
	return
}

func (ch *Character) SetFacing(dir int) {
	ch.Facing = CharacterFacingDirection[dir]
}
