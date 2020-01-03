package main

func ActionTeleport(gMap GameMap, to Direction) func(entity *Entity) {
	return func(entity *Entity) {
		entity.mTileX = to.x
		entity.mTileY = to.y
		entity.TeleportAndDraw(gMap, gMap.canvas)
	}
}
