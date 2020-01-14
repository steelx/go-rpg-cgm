package actions

import (
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"reflect"
)

//var Actions = make(map[string]func(gMap *game_map.GameMap, a ...interface{}) func(entity *game_map.Entity))
//
//func init() {
//	Actions["teleport"] = ActionTeleport
//}

//ActionTeleport : *GameMap, globals.Direction => *Entity => ()
func ActionTeleport(gMap *game_map.GameMap, a ...interface{}) func(entity *game_map.Entity) {
	aVal := reflect.ValueOf(a[0])
	to := aVal.Interface().(globals.Direction)
	return func(entity *game_map.Entity) {
		entity.TileX = to.X
		entity.TileY = to.Y
		entity.TeleportAndDraw(gMap, gMap.Canvas)
	}
}

func ActionAddNPC(gMap *game_map.GameMap, x, y float64) func(char *game_map.Character) {
	return func(char *game_map.Character) {
		char.Entity.SetTilePos(x, y)
		gMap.AddNPC(char)
	}
}
