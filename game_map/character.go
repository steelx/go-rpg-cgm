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
	Id                             string
	Anims                          map[string][]int
	Facing                         string
	Entity                         *Entity
	Controller                     *state_machine.StateMachine //[Name] -> [function that returns state]
	DefaultState, PrevDefaultState string                      //"wait"
	PathIndex                      int
	Path                           []string //e.g. ["up", "up", "up", "left", "right", "right",]
	TalkIndex                      int      //used during speech tracking
}

func CharacterCreate(
	def CharacterDefinition, controllerStates map[string]func() state_machine.State) *Character {
	player := &Character{
		Id:           def.Id,
		Facing:       def.FacingDirection,
		Entity:       CreateEntity(def.EntityDef),
		Controller:   state_machine.Create(controllerStates),
		DefaultState: def.DefaultState,
	}

	//AnimUp, AnimRight, AnimDown, AnimLeft []int
	player.Anims = make(map[string][]int, 0)

	if anim, ok := def.Animations["left"]; ok {
		player.Anims["left"] = anim
	}
	if anim, ok := def.Animations["right"]; ok {
		player.Anims["right"] = anim
	}
	if anim, ok := def.Animations["up"]; ok {
		player.Anims["up"] = anim
	}
	if anim, ok := def.Animations["down"]; ok {
		player.Anims["down"] = anim
	}

	return player
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

func (ch *Character) FollowPath(path []string) {
	ch.PathIndex = 0
	ch.Path = path
	//ch.PrevDefaultState = ch.DefaultState //this is causing problem
	ch.PrevDefaultState = "wait"
	ch.DefaultState = "follow_path"
	ch.Controller.Change("follow_path", nil)
}

func (ch *Character) GetCombatAnim(id string) []int {
	if _, ok := ch.Anims[id]; ok {
		return ch.Anims[id]
	} else {
		return []int{ch.Entity.StartFrame}
	}
}
