package main

func ActionTeleport(gMap GameMap, tileX, tileY float64) func(entity *Entity) {
	return func(entity *Entity) {
		entity.mTileX = tileX
		entity.mTileY = tileY
		entity.TeleportAndDraw(gMap, gMap.canvas)
	}
}
