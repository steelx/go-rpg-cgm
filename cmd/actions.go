package main

import (
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
)

func ActionTeleport(gMap game_map.GameMap, to globals.Direction) func(entity *game_map.Entity) {
	return func(entity *game_map.Entity) {
		entity.TileX = to.X
		entity.TileY = to.Y
		entity.TeleportAndDraw(gMap, gMap.Canvas)
	}
}
