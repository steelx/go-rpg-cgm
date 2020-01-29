package game_map

import (
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/utilz"
)

var Entities = make(map[string]EntityDefinition)
var Characters = make(map[string]func(gMap *GameMap) *Character)

func init() {

	walkCyclePng, err := utilz.LoadPicture("../resources/walk_cycle.png")
	utilz.PanicIfErr(err)

	sleepingPng, err := utilz.LoadPicture("../resources/sleeping.png")
	utilz.PanicIfErr(err)

	chestPng, err := utilz.LoadPicture("../resources/chest.png")
	utilz.PanicIfErr(err)

	//Entities
	Entities = map[string]EntityDefinition{
		"hero": {
			Texture: walkCyclePng,
			Width:   16, Height: 24,
			StartFrame: 24,
			TileX:      20,
			TileY:      20,
		},
		"thief": {
			Texture: walkCyclePng,
			Width:   16, Height: 24,
			StartFrame: 104,
			TileX:      11,
			TileY:      3,
		},
		"mage": {
			Texture: walkCyclePng,
			Width:   16, Height: 24,
			StartFrame: 120,
			TileX:      11,
			TileY:      3,
		},
		"sleeper": {
			Texture: sleepingPng,
			Width:   32, Height: 32,
			StartFrame: 12,
			TileX:      14,
			TileY:      19,
		},
		"npc1": {
			Texture: walkCyclePng,
			Width:   16, Height: 24,
			StartFrame: 46,
			TileX:      24,
			TileY:      19,
		},
		"npc2": {
			Texture: walkCyclePng,
			Width:   16, Height: 24,
			StartFrame: 56,
			TileX:      19,
			TileY:      24,
		},
		"prisoner": {
			Texture: walkCyclePng,
			Width:   16, Height: 24,
			StartFrame: 88,
			TileX:      19,
			TileY:      19, //jail map cords
		},
		"chest": {
			Texture: chestPng,
			Width:   16, Height: 16,
			StartFrame: 0,
			TileX:      20,
			TileY:      20,
		},
	}

	Characters["hero"] = hero
	Characters["thief"] = thief
	Characters["mage"] = mage
	Characters["sleeper"] = Sleeper
	Characters["npc1"] = NPC1
	Characters["npc2"] = NPC2
	Characters["guard"] = guard
	Characters["prisoner"] = prisoner
	Characters["chest"] = chest
}

func hero(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("hero",
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
	gameCharacter.Controller.Change("wait", Direction{0, 0})
	return gameCharacter
}

func thief(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("thief",
		map[string][]int{
			"up": {96, 97, 98, 99}, "right": {100, 101, 102, 103}, "down": {104, 105, 106, 107}, "left": {108, 109, 110, 111},
		},
		CharacterFacingDirection[2],
		Entities["thief"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStandStateCreate(gameCharacter, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", Direction{0, 0})
	return gameCharacter
}

func mage(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("mage",
		map[string][]int{
			"up": {112, 113, 114, 115}, "right": {116, 117, 118, 119}, "down": {120, 121, 122, 123}, "left": {124, 125, 126, 127},
		},
		CharacterFacingDirection[2],
		Entities["mage"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStandStateCreate(gameCharacter, gMap)
			},
			"move": func() state_machine.State {
				return MoveStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", Direction{0, 0})
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
	gameCharacter.Controller.Change("sleep", Direction{0, 0})
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
	gameCharacter.Controller.Change("wait", Direction{0, 0})
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
	gameCharacter.Controller.Change("wait", Direction{0, 0})
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
	gameCharacter.Controller.Change("wait", Direction{0, 0})
	return gameCharacter
}

func prisoner(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate("prisoner",
		map[string][]int{
			"up": {80, 81, 82, 83}, "right": {84, 85, 86, 87}, "down": {88, 89, 90, 91}, "left": {92, 93, 94, 95},
		},
		CharacterFacingDirection[2],
		Entities["prisoner"],
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
	gameCharacter.Controller.Change("wait", Direction{0, 0})
	return gameCharacter
}

func chest(gMap *GameMap) *Character {
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		"chest",
		map[string][]int{
			"down": {0, 1},
		},
		CharacterFacingDirection[2],
		Entities["chest"],
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStandStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", Direction{0, 0})
	return gameCharacter
}
