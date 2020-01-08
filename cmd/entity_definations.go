package main

import (
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

var (
	Hero *game_map.Character
	NPC1 *game_map.Character
	NPC2 *game_map.Character
)

func init() {
	pic, err := globals.LoadPicture("../resources/walk_cycle.png")
	globals.PanicIfErr(err)

	Hero = game_map.CharacterCreate("Ajinkya",
		[][]int{{16, 17, 18, 19}, {20, 21, 22, 23}, {24, 25, 26, 27}, {28, 29, 30, 31}},
		game_map.CharacterFacingDirection[2],
		game_map.CharacterDefinition{
			Texture: pic, Width: 16, Height: 24,
			StartFrame: 24,
			TileX:      2,
			TileY:      4,
			Map:        CastleRoomMap,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return character_states.WaitStateCreate(Hero, CastleRoomMap)
			},
			"move": func() state_machine.State {
				return character_states.MoveStateCreate(Hero, CastleRoomMap)
			},
		},
	)

	//Character:Create
	//Hero = &game_map.Character{
	//	Name:      "Ajinkya",
	//	AnimUp:    []int{16, 17, 18, 19},
	//	AnimRight: []int{20, 21, 22, 23},
	//	AnimDown:  []int{24, 25, 26, 27},
	//	AnimLeft:  []int{28, 29, 30, 31},
	//	Facing:    game_map.CharacterFacingDirection[2],
	//	Entity: game_map.CreateEntity(game_map.CharacterDefinition{
	//		Texture: pic, Width: 16, Height: 24,
	//		StartFrame: 24,
	//		TileX:      2,
	//		TileY:      4,
	//		Map:        CastleRoomMap,
	//	}),
	//	Controller: state_machine.Create(
	//		map[string]func() state_machine.State{
	//			"wait": func() state_machine.State {
	//				return character_states.WaitStateCreate(Hero, CastleRoomMap)
	//			},
	//			"move": func() state_machine.State {
	//				return character_states.MoveStateCreate(Hero, CastleRoomMap)
	//			},
	//		},
	//	),
	//}

	NPC1 = &game_map.Character{
		Name:   "Aghori Baba",
		Facing: game_map.CharacterFacingDirection[2],
		Entity: game_map.CreateEntity(game_map.CharacterDefinition{
			Texture: pic, Width: 16, Height: 24,
			StartFrame: 46,
			TileX:      9,
			TileY:      4,
		}),
		Controller: state_machine.Create(
			map[string]func() state_machine.State{
				"wait": func() state_machine.State {
					return character_states.NPCWaitStateCreate(NPC1, CastleRoomMap)
				},
			},
		),
	}

	NPC2 = &game_map.Character{
		Name:      "Bhadrasaal",
		AnimUp:    []int{48, 49, 50, 51},
		AnimRight: []int{52, 53, 54, 55},
		AnimDown:  []int{56, 57, 58, 59},
		AnimLeft:  []int{60, 61, 62, 63},
		Facing:    game_map.CharacterFacingDirection[2],
		Entity: game_map.CreateEntity(game_map.CharacterDefinition{
			Texture: pic, Width: 16, Height: 24,
			StartFrame: 56,
			TileX:      3,
			TileY:      8,
		}),
		Controller: state_machine.Create(
			map[string]func() state_machine.State{
				"wait": func() state_machine.State {
					return character_states.NPCStrollWaitStateCreate(NPC2, CastleRoomMap)
				},
				"move": func() state_machine.State {
					return character_states.MoveStateCreate(NPC2, CastleRoomMap)
				},
			},
		),
	}

	//Init Characters
	Hero.Controller.Change("wait", globals.Direction{0, 0})
	NPC1.Controller.Change("wait", globals.Direction{0, 0})
	NPC2.Controller.Change("wait", globals.Direction{0, 0})
}
