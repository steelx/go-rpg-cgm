package game_map

const (
	walkCyclePng   = "../resources/walk_cycle.png"
	sleepingPng    = "../resources/sleeping.png"
	chestPng       = "../resources/chest.png"
	combatHeroPng  = "../resources/combat_hero.png"
	combatMagePng  = "../resources/combat_mage.png"
	combatThiefPng = "../resources/combat_thief.png"
	goblinPng      = "../resources/goblin.png"
	combatSlashPng = "../resources/combat_slash.png"
	combatClawPng  = "../resources/combat_claw.png"
)

//Entities
var Entities = map[string]EntityDefinition{
	"empty": {
		Texture: "",
	},
	"combat_hero": {
		Texture: combatHeroPng,
		Width:   64, Height: 64,
		StartFrame: 10,
	},
	"combat_mage": {
		Texture: combatMagePng,
		Width:   64, Height: 64,
		StartFrame: 10,
	},
	"combat_thief": {
		Texture: combatThiefPng,
		Width:   64, Height: 64,
		StartFrame: 10,
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
	"goblin": {
		Texture: goblinPng,
		Width:   32, Height: 32,
		StartFrame: 0,
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
	"slash": {
		Texture: combatSlashPng,
		Width:   64, Height: 64,
		StartFrame: 2,
		Frames:     []int{2, 1, 0},
	},
	"claw": {
		Texture: combatClawPng,
		Width:   64, Height: 64,
		StartFrame: 0,
		Frames:     []int{0, 1, 2},
	},
}

var CharacterDefinitions map[string]CharacterDefinition = map[string]CharacterDefinition{
	"hero": {
		Id: "hero",
		Animations: map[string][]int{
			"up":    {16, 17, 18, 19},
			"right": {20, 21, 22, 23},
			"down":  {24, 25, 26, 27},
			"left":  {28, 29, 30, 31},

			csStandby: {41, 42, 43, 44, 45, 46, 47, 48, 49, 50},
			csMove:    {81, 82, 83, 88, 85, 86},
			csRetreat: {21, 22, 23, 24, 25, 26},
			csProne:   {5, 6},
			csAttack: {
				61, 62, 63, 64, 65, 66, 67, 68, 69, 70,
				//71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
			},
			csVictory: {71, 72, 73, 74, 75, 76, 77, 78, 79, 80},
			csUse:     {10, 11, 12, 13, 14},
			csHurt:    {113, 112, 111, 110, 109},
			csDie:     {114, 113, 112, 111, 110, 109, 108, 107, 106, 105, 104, 103, 102, 101},
			csDeath:   {35, 36, 37, 38},
		},
		FacingDirection:    CharacterFacingDirection[2],
		EntityDef:          Entities["hero"],
		CombatEntityDef:    Entities["combat_hero"],
		DefaultState:       "wait",
		DefaultCombatState: csNpcStand,
	},
	"thief": {
		Id: "thief",
		Animations: map[string][]int{
			"up":    {96, 97, 98, 99},
			"right": {100, 101, 102, 103},
			"down":  {104, 105, 106, 107},
			"left":  {108, 109, 110, 111},

			csStandby: {25, 26, 27, 28},
			csMove:    {24, 23, 22, 21, 20},
			csRetreat: {0, 1, 2, 3},
			csProne:   {5, 6, 7, 8, 9},
			csAttack: {
				60, 61, 62, 63, 64,
				55, 56, 57, 58, 59,
			},
			csVictory: {46, 47, 48, 49},
			csUse:     {15, 16, 17, 18, 19},
			csHurt:    {40, 41, 42, 43},
			csDie:     {35, 36, 37, 38},
			csDeath:   {35, 36, 37, 38},
		},
		FacingDirection:    CharacterFacingDirection[2],
		EntityDef:          Entities["thief"],
		CombatEntityDef:    Entities["combat_thief"],
		DefaultState:       "wait",
		DefaultCombatState: csNpcStand,
	},
	"mage": {
		Id: "mage",
		Animations: map[string][]int{
			"up":    {112, 113, 114, 115},
			"right": {116, 117, 118, 119},
			"down":  {120, 121, 122, 123},
			"left":  {124, 125, 126, 127},

			csMove:    {15, 16, 17, 18, 19},
			csStandby: {25, 26, 27, 28},
			csRetreat: {0, 1, 2, 3},
			csProne:   {5, 6},
			csAttack: {
				50, 51, 52, 53, 54,
				45, 46, 47, 48, 49,
			},
			csVictory: {46, 47, 48, 49},
			csUse:     {10, 11, 12, 13, 14},
			csHurt:    {40, 41, 42, 43},
			csDie:     {35, 36, 37, 38},
			csDeath:   {35, 36, 37, 38},
		},
		FacingDirection: CharacterFacingDirection[2],
		EntityDef:       Entities["mage"],
		CombatEntityDef: Entities["combat_mage"],
		DefaultState:    "wait",
	},
	"sleeper": {
		Id: "sleeper",
		Animations: map[string][]int{
			"left": {13},
		},
		FacingDirection: CharacterFacingDirection[3],
		EntityDef:       Entities["hero"],
		CombatEntityDef: Entities["empty"],
		DefaultState:    "wait",
	},
	"npc1": {
		Id:              "npc1",
		FacingDirection: CharacterFacingDirection[2],
		EntityDef:       Entities["npc1"],
		CombatEntityDef: Entities["empty"],
		DefaultState:    "wait",
	},
	"npc2": {
		Id: "npc2",
		Animations: map[string][]int{
			"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
		},
		FacingDirection: CharacterFacingDirection[2],
		EntityDef:       Entities["npc2"],
		CombatEntityDef: Entities["empty"],
		DefaultState:    "wait",
	},
	"guard": {
		Id: "guard",
		Animations: map[string][]int{
			"up": {48, 49, 50, 51}, "right": {52, 53, 54, 55}, "down": {56, 57, 58, 59}, "left": {60, 61, 62, 63},
		},
		FacingDirection: CharacterFacingDirection[2],
		EntityDef:       Entities["npc2"],
		CombatEntityDef: Entities["empty"],
		DefaultState:    "wait",
	},
	"prisoner": {
		Id: "prisoner",
		Animations: map[string][]int{
			"up": {80, 81, 82, 83}, "right": {84, 85, 86, 87}, "down": {88, 89, 90, 91}, "left": {92, 93, 94, 95},
		},
		FacingDirection: CharacterFacingDirection[2],
		EntityDef:       Entities["prisoner"],
		CombatEntityDef: Entities["empty"],
		DefaultState:    "wait",
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
	"goblin": {
		Id:                 "goblin",
		FacingDirection:    CharacterFacingDirection[2],
		EntityDef:          Entities["goblin"],
		DefaultState:       "wait",
		DefaultCombatState: csStandby,
		Animations: map[string][]int{
			csHurt: {0, 1},
		},
	},
}
