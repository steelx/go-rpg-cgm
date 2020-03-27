package game_map

import (
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
)

func mapTown(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("map_town.tmx")
	logFatalErr(err)

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 Collision",
		HiddenLayer:        "",
		Actions:            nil,
		TriggerTypes:       nil,
		Triggers:           nil,
		OnWake:             nil,
	}
}
