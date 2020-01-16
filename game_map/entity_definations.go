package game_map

import (
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

var Entities = make(map[string]EntityDefinition)
var Characters = make(map[string]func(gMap *GameMap) *Character)

func init() {

	walkCyclePng, err := globals.LoadPicture("../resources/walk_cycle.png")
	globals.PanicIfErr(err)

	sleepingPng, err := globals.LoadPicture("../resources/sleeping.png")
	globals.PanicIfErr(err)

	//Entities
	Entities = map[string]EntityDefinition{
		"hero": {
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 24,
			TileX:      20,
			TileY:      20,
		},
		"sleeper": {
			Texture: sleepingPng, Width: 32, Height: 32,
			StartFrame: 12,
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

	Characters["hero"] = chanakya
	Characters["sleeper"] = Sleeper
	Characters["npc1"] = NPC1
	Characters["npc2"] = NPC2
	Characters["guard"] = guard
}

func chanakya(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("Chanakya",
		map[string][]int{
			"up": {16, 17, 18, 19}, "right": {20, 21, 22, 23}, "down": {24, 25, 26, 27}, "left": {28, 29, 30, 31},
		},
		CharacterFacingDirection[2],
		Entities["hero"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return WaitStateCreate(gameCharacter, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", globals.Direction{0, 0})
	return gameCharacter
}

func Sleeper(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("Ajinkya",
		map[string][]int{
			"left": {13},
		},
		CharacterFacingDirection[3],
		Entities["hero"],
		map[string]func() state_machine.State{
			"sleep": func() state_machine.State {
				return SleepStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("sleep", globals.Direction{0, 0})
	return gameCharacter
}

func NPC1(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("Aghori Baba",
		nil,
		CharacterFacingDirection[2],
		Entities["npc1"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStandStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", globals.Direction{0, 0})
	return gameCharacter
}

func NPC2(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("Bhadrasaal",
		map[string][]int{
			"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
		},
		CharacterFacingDirection[2],
		Entities["npc2"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStrollWaitStateCreate(gameCharacter, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", globals.Direction{0, 0})
	return gameCharacter
}

func guard(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("guard",
		map[string][]int{
			"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
		},
		CharacterFacingDirection[2],
		Entities["npc2"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStandStateCreate(gameCharacter, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(gameCharacter, gMap)
			},
			"follow_path": func() state_machine.State {
				return FollowPathStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", globals.Direction{0, 0})
	return gameCharacter
}
