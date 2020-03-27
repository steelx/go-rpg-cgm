package game_map

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/world"
	"github.com/steelx/tilepix"
	"reflect"
)

func mapArena(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("map_arena.tmx")
	logFatalErr(err)

	worldV := reflect.ValueOf(gStack.Globals["world"])
	worldI := worldV.Interface().(*combat.WorldExtended)

	talkRecruit := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {

		char := gameMap.GetNPC(tileX, tileY)
		if char == nil {
			logrus.Info("Character not found at tile:", tileX, tileY)
			return
		}

		actorDef, ok := combat.PartyMembersDefinitions[char.Id]
		if !ok {
			logrus.Info("Missing actor definition at party_members_definitions.go")
			return
		}

		x, y := gameMap.GetTileIndex(tileX, tileY)
		playKeyItemFound := PlayBGSound("../sound/key_item.mp3")
		recruitNpc := func() {
			gStack.Pop() //remove selection menu

			gameMap.WriteTile(tileX, tileY, false)
			gameMap.RemoveTrigger(tileX, tileY)
			gameMap.RemoveNPC(tileX, tileY)
			playKeyItemFound()
			gStack.PushFitted(x, y, fmt.Sprintf(`"%s" joined your team, who is a "%s"`, actorDef.Name, actorDef.Id))
			worldI.Party.Add(combat.ActorCreate(actorDef))
		}

		choices := []string{
			"Sure",
			"See ya later!",
		}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				recruitNpc()
			}
		}

		gStack.PushSelectionMenu(x, y, 300, 70, "Can I join your party!", choices, onSelection, true)
	}

	addChest := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)

		loot := []world.ItemIndex{
			{1, 1},
			{2, 1},
			{6, 1},
			{10, 1},
		}

		playKeyItemFound := PlayBGSound("../sound/key_item.mp3")
		OnOpenChest := func() {
			gStack.Pop() //remove selection menu

			gameMap.WriteTile(tileX, tileY, false)
			gameMap.RemoveTrigger(tileX, tileY)

			playKeyItemFound()
			//gStack.PushFitted(x, y, "The chest is empty! lol")

			//Add Loot to world items
			for _, v := range loot {
				worldI.AddItem(v.Id, v.Count)

				name := world.ItemsDB[v.Id].Name
				message := fmt.Sprintf("Got %s", name)
				gStack.PushFitted(x, y, message)
			}
			chest := gameMap.GetNPC(tileX, tileY)
			chest.Entity.SetFrame(1)
		}

		choices := []string{
			"Open it",
			"Leave",
		}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				OnOpenChest()
			}
		}

		gStack.PushSelectionMenu(x, y, 300, 70, "You found a treasure chest", choices, onSelection, true)
	}

	enterTown := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)

		loadTownMap := func() {
			gStack.Pop() //remove selection menu

			storyboardEvents := []interface{}{
				Wait(0),
				FadeOutCharacter("handin", "hero", 1),

				BlackScreen("blackscreen"),
				Wait(1),
				KillState("blackscreen"),

				FadeOutMap("handin", 1),

				ReplaceScene("handin", "map_town", 1, 106, false, gStack.Win),
				PlayBGSound("../sound/reveal.mp3"),
				HandOffToMainStack("map_town"),
			}
			storyboard := StoryboardCreate(gStack, gStack.Win, storyboardEvents, true)
			gStack.Push(storyboard)
		}

		choices := []string{
			"Go to Town",
			"Stay at Arena",
		}
		onSelection := func(index int, c interface{}) {
			if index == 0 {
				loadTownMap()
			}
		}

		gStack.PushSelectionMenu(x, y, 400, 70, "You get spotted by Goblins", choices, onSelection, true)
	}

	enterArena := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		prevState := gStack.Pop()
		gStack.Push(ArenaStateCreate(gStack, prevState))
	}

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
		HiddenLayer:        "",

		OnWake: map[string][]TriggerParam{
			"AddNPC": {
				{Id: "mage", X: 26, Y: 14},
				{Id: "thief", X: 27, Y: 14},
			},
			"AddChest": {
				{Id: "chest", X: 17, Y: 14},
			},
		},

		Actions: map[string]MapAction{
			"talk_recruit": {
				Id:     "RunScript",
				Script: talkRecruit,
			},
			"add_chest": {
				Id:     "RunScript",
				Script: addChest,
			},
			"enter_town": {
				Id:     "RunScript",
				Script: enterTown,
			},
			"enter_arena": {
				Id:     "RunScript",
				Script: enterArena,
			},
		},
		TriggerTypes: map[string]TriggerType{
			"talk_recruit_at_alley": {
				OnUse: "talk_recruit",
			},
			"add_chest_1": {
				OnUse: "add_chest",
			},
			"enter_town_at_gate": {
				OnEnter: "enter_town",
			},
			"enter_arena_at_door": {
				OnUse: "enter_arena",
			},
		},
		Triggers: []TriggerParam{
			{Id: "talk_recruit_at_alley", X: 26, Y: 14},
			{Id: "talk_recruit_at_alley", X: 27, Y: 14},
			{Id: "add_chest_1", X: 17, Y: 14},
			{Id: "enter_town_at_gate", X: 22, Y: 42},
			{Id: "enter_town_at_gate", X: 23, Y: 42},
			{Id: "enter_town_at_gate", X: 24, Y: 42},
			{Id: "enter_arena_at_door", X: 22, Y: 13},
			{Id: "enter_arena_at_door", X: 23, Y: 13},
			{Id: "enter_arena_at_door", X: 24, Y: 13},
		},
	}
}
