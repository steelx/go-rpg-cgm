package game_map

import (
	"reflect"
)

var LIST = make(map[string]func(gMap *GameMap, a ...interface{}) func(b ...interface{}))

func init() {
	LIST["AddNPC"] = addNPC_I
	LIST["AddChest"] = addChest_I
}

//ActionTeleport : *GameMap, Direction => *Entity => ()
func ActionTeleport(gMap *GameMap, a ...interface{}) func(entity *Entity) {
	aVal := reflect.ValueOf(a[0])
	to := aVal.Interface().(Direction)
	return func(entity *Entity) {
		entity.TileX = to.X
		entity.TileY = to.Y
		entity.TeleportAndDraw(gMap, gMap.Canvas)
	}
}

func AddNPC(gMap *GameMap, x, y float64) func(char *Character) {
	return func(char *Character) {
		char.Entity.SetTilePos(x, y)
		gMap.AddNPC(char)
		gMap.GoToTile(x, y)
	}
}

//len a MUST BE 2
func addNPC_I(gMap *GameMap, a ...interface{}) func(b ...interface{}) {
	a0 := reflect.ValueOf(a[0])
	a1 := reflect.ValueOf(a[1])
	x := a0.Interface().(float64)
	y := a1.Interface().(float64)
	return func(b ...interface{}) {
		bVal := reflect.ValueOf(b[0])
		char := bVal.Interface().(*Character)
		char.Entity.SetTilePos(x, y)
		gMap.AddNPC(char)
		gMap.GoToTile(x, y)
	}
}

func addChest_I(gMap *GameMap, a ...interface{}) func(b ...interface{}) {
	a0 := reflect.ValueOf(a[0])
	a1 := reflect.ValueOf(a[1])
	tileX := a0.Interface().(float64)
	tileY := a1.Interface().(float64)
	return func(b ...interface{}) {
		// Add Chest on map
		b0 := reflect.ValueOf(b[0])
		char := b0.Interface().(*Character)
		char.Entity.SetTilePos(tileX, tileY)
		gMap.AddNPC(char)
	}
}
