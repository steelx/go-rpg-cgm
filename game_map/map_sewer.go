package game_map

import (
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
)

func mapSewer(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("map_sewer.tmx")
	logFatalErr(err)
	//end game 39, 4-5-6

	loadArenaMapEvents := []interface{}{
		Wait(1),
		FadeOutCharacter("handin", "hero", 2),

		BlackScreen("blackscreen"),
		Wait(1),
		KillState("blackscreen"),

		ReplaceScene("handin", "map_arena", 9, 7, false, gStack.Win),
		PlayBGSound("../sound/reveal.mp3"),
		HandOffToMainStack("map_arena"),
	}

	endOnEnter := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)

		endEnter := func(gMap *GameMap) {
			gStack.Pop() //remove SelectionMenu

			loadArenaMap := StoryboardCreate(gStack, gStack.Win, loadArenaMapEvents, true)
			gStack.Push(loadArenaMap)
		}

		choices := []string{"Hit space to exit"}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				endEnter(gameMap)
			}
		}
		gStack.PushSelectionMenu(
			x, y, 400, 70,
			"You have reached the end of sewer", choices, onSelection, false)
	}

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     1,
		CollisionLayerName: "01 collision",
		Actions: map[string]MapAction{
			"end_on_enter": {
				Id:     "RunScript",
				Script: endOnEnter,
			},
		},
		TriggerTypes: map[string]TriggerType{
			"end_of_sewer": {
				OnEnter: "end_on_enter",
			},
		},
		Triggers: []TriggerParam{
			{Id: "end_of_sewer", X: 39, Y: 4},
			{Id: "end_of_sewer", X: 39, Y: 5},
			{Id: "end_of_sewer", X: 39, Y: 6},
		},
		OnWake: nil,
	}
}
