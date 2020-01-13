package character_states

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

var walkCyclePng pixel.Picture
var sleepingPng pixel.Picture

//var gEntities map[string]EntityDefinition

func init() {
	var err error
	walkCyclePng, err = globals.LoadPicture("../resources/walk_cycle.png")
	globals.PanicIfErr(err)

	sleepingPng, err = globals.LoadPicture("../resources/sleeping.png")
	globals.PanicIfErr(err)

	//gEntities
	//gEntities = map[string]EntityDefinition{
	//	"hero":
	//}
}

func Hero(startPos pixel.Vec, gMap *game_map.GameMap) *game_map.Character {
	var gameCharacter *game_map.Character
	gameCharacter = game_map.CharacterCreate("Ajinkya",
		[][]int{{16, 17, 18, 19}, {20, 21, 22, 23}, {24, 25, 26, 27}, {28, 29, 30, 31}},
		game_map.CharacterFacingDirection[2],
		game_map.EntityDefinition{
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 24,
			TileX:      startPos.X,
			TileY:      startPos.Y,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return WaitStateCreate(gameCharacter, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(gameCharacter, gMap)
			},
			"sleep": func() state_machine.State {
				return SleepStateCreate(gameCharacter, gMap)
			},
		},
	)
	return gameCharacter
}

func Sleeper(startPos pixel.Vec, gMap *game_map.GameMap) *game_map.Character {
	var sleeping *game_map.Character
	sleeping = game_map.CharacterCreate("Ajinkya",
		[][]int{{}, {}, {}, {13}}, //left side only
		game_map.CharacterFacingDirection[3],
		game_map.EntityDefinition{
			Texture: sleepingPng, Width: 32, Height: 32,
			StartFrame: 3,
			TileX:      startPos.X, //18
			TileY:      startPos.Y, //32
		},
		map[string]func() state_machine.State{
			"sleep": func() state_machine.State {
				return SleepStateCreate(sleeping, gMap)
			},
		},
	)
	return sleeping
}

func NPC1(startPos pixel.Vec, gMap *game_map.GameMap) *game_map.Character {
	var NPC *game_map.Character
	NPC = game_map.CharacterCreate("Aghori Baba",
		nil,
		game_map.CharacterFacingDirection[2],
		game_map.EntityDefinition{
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 46,
			TileX:      startPos.X,
			TileY:      startPos.Y,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCWaitStateCreate(NPC, gMap)
			},
		},
	)
	return NPC
}

func NPC2(startPos pixel.Vec, gMap *game_map.GameMap) *game_map.Character {
	var NPC *game_map.Character
	NPC = game_map.CharacterCreate("Bhadrasaal",
		[][]int{{48, 49, 50, 51}, {52, 53, 54, 55}, {56, 57, 58, 59}, {60, 61, 62, 63}},
		game_map.CharacterFacingDirection[2],
		game_map.EntityDefinition{
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 56,
			TileX:      startPos.X,
			TileY:      startPos.Y,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStrollWaitStateCreate(NPC, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(NPC, gMap)
			},
		},
	)
	return NPC
}
