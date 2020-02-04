package game_map

import (
	"github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/tilepix"
	"reflect"
)

func mapJail(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("jail.tmx")
	logFatalErr(err)

	boneItemId := 4

	worldV := reflect.ValueOf(gStack.Globals["world"])
	worldI := worldV.Interface().(*combat.WorldExtended)

	playKeyItemFound := PlayBGSound("../sound/key_item.mp3")
	playSkeletonDestroyed := PlayBGSound("../sound/skeleton_destroy.mp3")
	boneScript := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)
		giveBone := func(gMap *GameMap) {
			playKeyItemFound()
			gStack.Pop() //remove selection menu
			gStack.PushFitted(x, y, `Found key item: "Calcified bone"`)
			worldI.AddKeyItem(boneItemId)
		}

		choices := []string{"Hit space to add it to your Inventory"}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				giveBone(gameMap)
			}
		}
		gStack.PushSelectionMenu(
			x, y, 400, 70,
			"The skeleton collapsed into dust.", choices, onSelection, false)
		//since skeleton occupied 2 tiles on Tiled
		gameMap.RemoveTrigger(41, 22)
		gameMap.RemoveTrigger(42, 22)
		//removed collision from skeleton tile
		gameMap.WriteTile(41, 22, false)
		gameMap.WriteTile(42, 22, false)
		playSkeletonDestroyed()
	}

	playCrumble := PlayBGSound("../sound/crumble.mp3")
	breakWallScript := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)
		onPush := func(gMap *GameMap) {
			// The player's pushing the wall.
			gStack.Pop() //remove selection menu
			gStack.PushFitted(x, y, "The wall crumbles.")
			playCrumble()

			//see below Triggers - "cracked_stone"
			gMap.RemoveTrigger(35, 22)
			gMap.WriteTile(35, 22, false)
			gMap.SetHiddenTileVisible(35, 22)
		}
		choices := []string{
			"Push the wall",
			"Get back!",
		}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				onPush(gameMap)
			}
		}
		gStack.PushSelectionMenu(x, y, 400, 100, "The wall here is crumbling. Push it?", choices, onSelection, false)
		//gStack.PushFITMenu(x, y, "The wall here is crumbling. Push it?", choices, onSelection)
	}

	moveGregor := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		if worldI.HasKey(boneItemId) {
			prisoner, ok := gameMap.NPCbyId["prisoner"]
			if !ok {
				logrus.Fatal("GameMap prisoner not found!")
				return
			}
			//started at 20, 29
			prisoner.FollowPath(
				[]string{
					"up", "up", "up", "up", "up", "up", "up", "up", "up", "up",
					"right", "right", "right", "right", "right",
					"down", "down", "down", "down",
				},
			)
			gameMap.RemoveTrigger(tileX, tileY)
		}
	}

	avatarPng, err := utilz.LoadPicture("../resources/avatar.png")
	utilz.PanicIfErr(err)

	talkGregor := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		prisoner, ok := gameMap.NPCbyId["prisoner"]
		if !ok {
			logrus.Fatal("GameMap prisoner not found!")
			return
		}

		if prisoner.Entity.TileX == 25 && prisoner.Entity.TileY == 22 {
			speech := []string{
				"You're another black blood aren't you?",
				"Tomorrow morning, they'll kill you, just like the others.",
				"If I was you, I'd try and escape. Pry the drain open, with that big bone you're holding.",
			}
			//prisoner.TalkIndex
			textMsg := speech[prisoner.TalkIndex]

			x, y := gameMap.GetTileIndex(tileX, tileY)
			gStack.PushFixed(x, y, 400, 100, textMsg, "Prisoner", avatarPng)
			prisoner.TalkIndex++
			if prisoner.TalkIndex >= len(speech) {
				prisoner.TalkIndex = 0
			}
		}
	}

	playUnlock := PlayBGSound("../sound/unlock.mp3")
	grillOnUse := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		if !worldI.HasKey(boneItemId) {
			return
		}

		x, y := gameMap.GetTileIndex(tileX, tileY)
		grillOpen := func(gMap *GameMap) {
			gStack.Pop()
			gStack.PushFitted(x, y, "The grill opened, and leads a way inside the sewers")
			playUnlock()

			gMap.RemoveTrigger(32, 15)
			gMap.RemoveTrigger(33, 15)
			gMap.WriteTile(32, 15, false)
			gMap.WriteTile(33, 15, false)
			gMap.SetHiddenTileVisible(32, 15)
			gMap.SetHiddenTileVisible(33, 15)

			//now we add new trigger onEnter
			gMap.AddTrigger("grill_when_open", 32, 15)
			gMap.AddTrigger("grill_when_open", 33, 15)
		}

		choices := []string{"Pry open the grill", "Leave it alone"}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				grillOpen(gameMap)
			}
		}
		gStack.PushSelectionMenu(
			x, y, 400, 110,
			"Do you want to dig teh grill?", choices, onSelection, false)
	}

	//jail break
	jailBreakCutsceneEvents := []interface{}{
		Wait(0),
		FadeOutCharacter("handin", "hero", 2),
		//Wait(1),
		//
		//MoveCamToTile("handin", 32, 15, 16, 31, 2),
		//RunActionAddNPC("handin", "guard", 16, 31, 0),
		//MoveNPC("prisoner", "handin", []string{
		//	"up", "up", "up", "up",
		//	"left", "left", "left", "left", "left",
		//	"down", "down", "down", "down", "down", "down",
		//}),
		//Wait(0),
		//MoveNPC("guard", "handin", []string{
		//	"right", "right", "right", "right", "right",
		//	"up",
		//}),
		//WriteTile("handin", 21, 30, false), //at jail door
		//PlaySound("../sound/unlock.mp3", 1),
		//SetHiddenTileVisible("handin", 21, 29), //jail door
		//SetHiddenTileVisible("handin", 21, 28), //jail door
		//MoveNPC("guard", "handin", []string{
		//	"up", "up", "up", "up", "up", "up",
		//}),
		//Say("handin", "guard", "Has the other prisoner gone?", 2),
		//Say("handin", "prisoner", "Yeah.", 1),
		//Wait(1),
		//Say("handin", "guard", "Hmm", 1),
		//Say("handin", "guard", "Dhananand wants to see you in the Tower", 2),
		//Wait(1),
		//FadeOutMap("handin", 2),
		//BlackScreen("blackscreen"),
		//Wait(1),
		//KillState("blackscreen"),

		ReplaceScene("handin", "map_sewer", 3, 5, false, gStack.Win),
		PlayBGSound("../sound/reveal.mp3"),
		HandOffToMainStack("map_sewer"),
	}
	grillOnEnter := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)
		grillEnter := func(gMap *GameMap) {
			gStack.Pop()
			gMap.RemoveTrigger(32, 15)
			gMap.RemoveTrigger(33, 15)

			jailBreakCutscene := StoryboardCreate(gStack, gStack.Win, jailBreakCutsceneEvents, true)
			gStack.Push(jailBreakCutscene)
		}
		choices := []string{"HIT space to enter the Tunnel"}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				grillEnter(gameMap)
			}
		}
		gStack.PushSelectionMenu(
			x, y, 400, 70,
			"There's a tunnel behind the grate", choices, onSelection, false)
	}

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
		HiddenLayer:        "01 detail",
		OnWake: map[string][]TriggerParam{
			"AddNPC": {
				{Id: "prisoner", X: 20, Y: 29},
			},
		},
		Actions: map[string]MapAction{
			"break_wall_script": {
				Id:     "RunScript",
				Script: breakWallScript,
			},
			"bone_script": {
				Id:     "RunScript",
				Script: boneScript,
			},
			"move_gregor": {
				Id:     "RunScript",
				Script: moveGregor,
			},
			"talk_gregor": {
				Id:     "RunScript",
				Script: talkGregor,
			},
			"grill_on_use": {
				Id:     "RunScript",
				Script: grillOnUse,
			},
			"grill_on_enter": {
				Id:     "RunScript",
				Script: grillOnEnter,
			},
		},
		TriggerTypes: map[string]TriggerType{
			"cracked_stone": {
				OnUse: "break_wall_script",
			},
			"calcified_bone": {
				OnUse: "bone_script",
			},
			"gregor_talk_trigger": {
				OnUse: "talk_gregor",
			},
			"gregor_move_trigger": {
				OnExit: "move_gregor",
			},
			"grill_when_closed": {
				OnUse: "grill_on_use",
			},
			"grill_when_open": {
				OnEnter: "grill_on_enter",
			},
		},
		Triggers: []TriggerParam{
			{Id: "cracked_stone", X: 35, Y: 22},
			{Id: "calcified_bone", X: 41, Y: 22},
			{Id: "calcified_bone", X: 42, Y: 22},
			{Id: "gregor_move_trigger", X: 36, Y: 22}, //at cracked_stone
			{Id: "gregor_talk_trigger", X: 25, Y: 23}, //at prisoner door
			{Id: "gregor_talk_trigger", X: 26, Y: 23}, //at prisoner door
			{Id: "grill_when_closed", X: 32, Y: 15},   //at grill closed
			{Id: "grill_when_closed", X: 33, Y: 15},
		},
	}
}
