package game_map

import (
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
)

func mapSewer(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("map_sewer.tmx")
	logFatalErr(err)
	//end game 39, 4-5-6

	//gameOverEvents := []interface{}{
	//	Wait(1),
	//	TitleCaptionScreen("title", "Game Over", 3),
	//	SubTitleCaptionScreen("subtitle", "picture abhi baki hain mere dost...", 2),
	//	Wait(3),
	//	KillState("title"),
	//	KillState("subtitle"),
	//	BlackScreen("blackscreen"),
	//	Wait(1),
	//	KillState("blackscreen"),
	//}

	endOnEnter := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)

		endEnter := func(gMap *GameMap) {
			gStack.Pop() //remove SelectionMenu
			gStack.Pop() //remove map Sewer

			gameOver := GameOverStateCreate(gStack, []gui.CaptionStyle{
				{"Game Over", 3},
				{"will be continued...", 1.8},
				{"press Q to quit", 1},
			})
			//gameOverSB := StoryboardCreate(gStack, globals.Global.Win, gameOverEvents, false)
			gStack.Push(gameOver)
		}

		choices := []string{"Hit space to exit"}
		onSelection := func(index int, c string) {
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
