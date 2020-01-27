package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
	"reflect"
)

func mapArena(gStack *gui.StateStack) MapInfo {

	gMap, err := tilepix.ReadFile("map_arena.tmx")
	logFatalErr(err)

	talkRecruit := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {

		x, y := gameMap.GetTileIndex(tileX, tileY)
		worldV := reflect.ValueOf(gStack.Globals["world"])
		worldI := worldV.Interface().(*combat.WorldExtended)

		char := gameMap.GetNPC(tileX, tileY)
		if char == nil {
			fmt.Println("Character not found at tile:", tileX, tileY)
			return
		}

		actorDef, ok := combat.PartyMembersDefinitions[char.Id]
		if !ok {
			fmt.Println("Missing actor definition at party_members_definitions.go")
			return
		}

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

		gStack.PushSelectionMenu(x, y, 300, 100, "Can I join your party!", choices, onSelection, true)
	}

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
		HiddenLayer:        "",

		OnWake: map[string][]TriggerParam{
			"AddNPC": {
				{Id: "mage", X: 36, Y: 12},
				{Id: "thief", X: 37, Y: 10},
			},
		},

		Actions: map[string]MapAction{
			"talk_recruit": {
				Id:     "RunScript",
				Script: talkRecruit,
			},
		},
		TriggerTypes: map[string]TriggerType{
			"talk_recruit_at_alley": {
				OnUse: "talk_recruit",
			},
		},
		Triggers: []TriggerParam{
			{"talk_recruit_at_alley", 36, 12},
			{"talk_recruit_at_alley", 37, 10},
		},
	}
}
