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
	fxRestoreHpPng = "../resources/fx_restore_hp.png"
	fxRestoreMpPng = "../resources/fx_restore_mp.png"
	fxRevivePng    = "../resources/fx_revive.png"
	fxUseItemPng   = "../resources/fx_use_item.png"
	fxFirePng      = "../resources/fx_fire.png"
	fxIcePng       = "../resources/fx_ice.png"
	fxElectricPng  = "../resources/fx_electric.png"
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
	"fx_restore_hp": {
		Texture:    fxRestoreHpPng,
		Width:      16,
		Height:     16,
		StartFrame: 0,
		Frames:     []int{0, 1, 2, 3, 4},
	},
	"fx_restore_mp": {
		Texture:    fxRestoreMpPng,
		Width:      16,
		Height:     16,
		StartFrame: 0,
		Frames:     []int{0, 1, 2, 3, 4, 5},
	},
	"fx_revive": {
		Texture:    fxRevivePng,
		Width:      16,
		Height:     16,
		StartFrame: 0,
		Frames:     []int{0, 1, 2, 3, 4, 5, 6, 7},
	},
	"fx_use_item": {
		Texture:    fxUseItemPng,
		Width:      16,
		Height:     16,
		StartFrame: 0,
		Frames:     []int{0, 1, 2, 3, 3, 2, 1, 0},
	},
	"fx_fire": {
		Texture:    fxFirePng,
		Width:      32,
		Height:     48,
		StartFrame: 1,
		Frames:     []int{0, 1, 2},
	},
	"fx_electric": {
		Texture:    fxElectricPng,
		Width:      32,
		Height:     16,
		StartFrame: 1,
		Frames:     []int{0, 1, 2},
	},
	"fx_ice_1": {
		Texture:    fxIcePng,
		Width:      16,
		Height:     16,
		StartFrame: 1,
		Frames:     []int{0, 1, 2, 3},
	},
	"fx_ice_2": {
		Texture:    fxIcePng,
		Width:      16,
		Height:     16,
		StartFrame: 5,
		Frames:     []int{4, 5, 6, 7},
	},
	"fx_ice_3": {
		Texture:    fxIcePng,
		Width:      16,
		Height:     16,
		StartFrame: 9,
		Frames:     []int{8, 9, 10, 11},
	},
	"fx_ice_spark": {
		Texture:    fxIcePng,
		Width:      16,
		Height:     16,
		StartFrame: 13,
		Frames:     []int{12, 13, 14, 15},
	},
}
