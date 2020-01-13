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

type Direction struct {
	X, Y float64
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

func CharacterCreate(
	name string, animations [][]int, facingDirection string, charDef EntityDefinition, controllerStates map[string]func() state_machine.State) *Character {
	player := &Character{
		Name:       name,
		Facing:     facingDirection,
		Entity:     CreateEntity(charDef),
		Controller: state_machine.Create(controllerStates),
	}
	if animations != nil && len(animations) == 4 {
		player.AnimUp = animations[0]
		player.AnimRight = animations[1]
		player.AnimDown = animations[2]
		player.AnimLeft = animations[3]
	}
	return player
}
