package game_map

import (
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
)

func mapSewer(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("jail.tmx")
	logFatalErr(err)

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
		HiddenLayer:        "01 detail",
		Actions:            nil,
		TriggerTypes:       nil,
		Triggers:           nil,
		OnWake:             nil,
	}
}
