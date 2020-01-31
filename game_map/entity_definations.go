package game_map

import (
	"github.com/steelx/go-rpg-cgm/state_machine"
)

var Entities = make(map[string]EntityDefinition)
var Characters = make(map[string]func(gMap *GameMap) *Character)
var CharacterDefinitions = make(map[string]CharacterDefinition)

func init() {
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

type CharacterDefinition struct {
	Id                         string
	Animations                 map[string][]int
	FacingDirection            string
	EntityDef, CombatEntityDef EntityDefinition
	CombatStates               map[string]func(args ...interface{}) state_machine.State
	DefaultCombatState,
	DefaultState string
}

func hero(gMap *GameMap) *Character {
	charDef := CharacterDefinitions["hero"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["thief"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["mage"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["sleeper"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["npc1"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["npc2"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["guard"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["prisoner"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
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
	charDef := CharacterDefinitions["chest"]
	var gameCharacter *Character
	gameCharacter = CharacterCreate(
		charDef,
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return NPCStandStateCreate(gameCharacter, gMap)
			},
		},
	)
	gameCharacter.Controller.Change("wait", Direction{0, 0})
	return gameCharacter
}
