package game_map

import (
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
)

func mapArena(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("map_arena.tmx")
	logFatalErr(err)

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
		HiddenLayer:        "",
		Actions:            nil,
		TriggerTypes:       nil,
		Triggers:           nil,
		OnWake:             nil,
	}
}
