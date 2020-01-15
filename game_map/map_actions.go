package game_map

func RunScript(script func(gMap *GameMap, entity *Entity, x, y float64)) func(gMap *GameMap, entity *Entity, x, y float64) {

	return func(gMap *GameMap, entity *Entity, x, y float64) {
		script(gMap, entity, x, y)
	}
}
