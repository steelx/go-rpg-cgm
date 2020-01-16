package game_map

import (
	"github.com/steelx/go-rpg-cgm/utilz"
	"reflect"
)

var LIST = make(map[string]func(gMap *GameMap, a ...interface{}) func(b interface{}))

func init() {
	LIST["AddNPC"] = AddNPC_I
}

//ActionTeleport : *GameMap, globals.Direction => *Entity => ()
func ActionTeleport(gMap *GameMap, a ...interface{}) func(entity *Entity) {
	aVal := reflect.ValueOf(a[0])
	to := aVal.Interface().(utilz.Direction)
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
func AddNPC_I(gMap *GameMap, a ...interface{}) func(b interface{}) {
	a0 := reflect.ValueOf(a[0])
	a1 := reflect.ValueOf(a[1])
	x := a0.Interface().(float64)
	y := a1.Interface().(float64)
	return func(b interface{}) {
		bVal := reflect.ValueOf(b)
		char := bVal.Interface().(*Character)
		char.Entity.SetTilePos(x, y)
		gMap.AddNPC(char)
		gMap.GoToTile(x, y)
	}
}
