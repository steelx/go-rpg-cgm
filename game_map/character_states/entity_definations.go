package character_states

import (
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

var gEntities map[string]game_map.EntityDefinition

func init() {

	walkCyclePng, err := globals.LoadPicture("../resources/walk_cycle.png")
	globals.PanicIfErr(err)

	sleepingPng, err := globals.LoadPicture("../resources/sleeping.png")
	globals.PanicIfErr(err)

	//gEntities
	gEntities = map[string]game_map.EntityDefinition{
		"hero": {
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 24,
			TileX:      20,
			TileY:      20,
		},
		"sleeper": {
			Texture: sleepingPng, Width: 32, Height: 32,
			StartFrame: 3,
			TileX:      14,
			TileY:      19,
		},
		"npc1": {
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 46,
			TileX:      24,
			TileY:      19,
		},
		"npc2": {
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 56,
			TileX:      19,
			TileY:      24,
		},
	}
}

func Hero(gMap *game_map.GameMap) *game_map.Character {
	var gameCharacter *game_map.Character
	gameCharacter = game_map.CharacterCreate("Ajinkya",
		map[string][]int{
			"up": {16, 17, 18, 19}, "right": {20, 21, 22, 23}, "down": {24, 25, 26, 27}, "left": {28, 29, 30, 31},
		},
		game_map.CharacterFacingDirection[2],
		gEntities["hero"],
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

func NPC1(gMap *game_map.GameMap) *game_map.Character {
	var NPC *game_map.Character
	NPC = game_map.CharacterCreate("Aghori Baba",
		nil,
		game_map.CharacterFacingDirection[2],
		gEntities["npc1"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCWaitStateCreate(NPC, gMap)
			},
		},
	)
	return NPC
}

func NPC2(gMap *game_map.GameMap) *game_map.Character {
	var NPC *game_map.Character
	NPC = game_map.CharacterCreate("Bhadrasaal",
		map[string][]int{
			"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
		},
		game_map.CharacterFacingDirection[2],
		gEntities["npc2"],
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
