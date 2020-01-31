package game_map

import "github.com/steelx/go-rpg-cgm/utilz"

func init() {
	walkCyclePng, err := utilz.LoadPicture("../resources/walk_cycle.png")
	utilz.PanicIfErr(err)

	sleepingPng, err := utilz.LoadPicture("../resources/sleeping.png")
	utilz.PanicIfErr(err)

	chestPng, err := utilz.LoadPicture("../resources/chest.png")
	utilz.PanicIfErr(err)

	combatHeroPng, err := utilz.LoadPicture("../resources/combat_hero.png")
	utilz.PanicIfErr(err)
	combatMagePng, err := utilz.LoadPicture("../resources/combat_mage.png")
	utilz.PanicIfErr(err)
	combatThiefPng, err := utilz.LoadPicture("../resources/combat_thief.png")
	utilz.PanicIfErr(err)

	//Entities
	Entities = map[string]EntityDefinition{
		"empty": {
			Texture: nil,
		},
		"combat_hero": {
			Texture: combatHeroPng,
			Width:   64, Height: 64,
			StartFrame: 37,
		},
		"combat_mage": {
			Texture: combatMagePng,
			Width:   64, Height: 64,
			StartFrame: 37,
		},
		"combat_thief": {
			Texture: combatThiefPng,
			Width:   64, Height: 64,
			StartFrame: 37,
		},
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

	CharacterDefinitions = map[string]CharacterDefinition{
		"hero": {
			Id: "hero",
			Animations: map[string][]int{
				"standby": {15, 16, 17, 18},
				"up":      {16, 17, 18, 19},
				"right":   {20, 21, 22, 23},
				"down":    {24, 25, 26, 27},
				"left":    {28, 29, 30, 31},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["hero"],
			CombatEntityDef: Entities["combat_hero"],
		},
		"thief": {
			Id: "thief",
			Animations: map[string][]int{
				"standby": {15, 16, 17, 18},
				"up":      {96, 97, 98, 99}, "right": {100, 101, 102, 103}, "down": {104, 105, 106, 107}, "left": {108, 109, 110, 111},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["thief"],
			CombatEntityDef: Entities["combat_thief"],
		},
		"mage": {
			Id: "mage",
			Animations: map[string][]int{
				"standby": {15, 16, 17, 18},
				"up":      {112, 113, 114, 115}, "right": {116, 117, 118, 119}, "down": {120, 121, 122, 123}, "left": {124, 125, 126, 127},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["mage"],
			CombatEntityDef: Entities["combat_mage"],
		},
		"sleeper": {
			Id: "sleeper",
			Animations: map[string][]int{
				"left": {13},
			},
			FacingDirection: CharacterFacingDirection[3],
			EntityDef:       Entities["hero"],
			CombatEntityDef: Entities["empty"],
		},
		"npc1": {
			Id:              "npc1",
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["npc1"],
			CombatEntityDef: Entities["empty"],
		},
		"npc2": {
			Id: "npc2",
			Animations: map[string][]int{
				"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["npc2"],
			CombatEntityDef: Entities["empty"],
		},
		"guard": {
			Id: "guard",
			Animations: map[string][]int{
				"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["npc2"],
			CombatEntityDef: Entities["empty"],
		},
		"prisoner": {
			Id: "prisoner",
			Animations: map[string][]int{
				"up": {80, 81, 82, 83}, "right": {84, 85, 86, 87}, "down": {88, 89, 90, 91}, "left": {92, 93, 94, 95},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["prisoner"],
			CombatEntityDef: Entities["empty"],
		},
		"chest": {
			Id: "chest",
			Animations: map[string][]int{
				"down": {0, 1},
			},
			FacingDirection: CharacterFacingDirection[2],
			EntityDef:       Entities["chest"],
			CombatEntityDef: Entities["empty"],
		},
	}
}
