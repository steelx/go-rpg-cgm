package game_map

import (
	"github.com/bcvery1/tilepix"
	"github.com/steelx/go-rpg-cgm/gui"
	"log"
	"reflect"
)

type MapAction struct {
	Id     string
	Script func(gameMap *GameMap, entity *Entity, x, y float64)
}

type TriggerType struct {
	OnUse   string
	OnEnter string
	OnExit  string
}

type TriggerParam struct {
	Id   string
	X, Y float64
}

type MapInfo struct {
	Tilemap            *tilepix.Map
	CollisionLayer     int
	CollisionLayerName string
	Actions            map[string]MapAction   //"break_wall_script" : { Id = "RunScript", Scripts : []{ CrumbleScript } }
	TriggerTypes       map[string]TriggerType //"cracked_stone" : { OnUse = "break_wall_script" }
	Triggers           []TriggerParam         //[]{Id = "cracked_stone", x = 60, y = 11}
}

var MapsDB map[string]func(gStack *gui.StateStack) MapInfo

func init() {
	MapsDB = make(map[string]func(gStack *gui.StateStack) MapInfo)
	MapsDB["player_room"] = playerHouseMap
	MapsDB["small_room"] = smallRoomMap
	MapsDB["jail_room"] = jailRoomMap
}

//player render rule is we render them with Collision Layer
func playerHouseMap(gStack *gui.StateStack) MapInfo {
	//gStack could be global stack in future

	gMap, err := tilepix.ReadFile("sontos_house.tmx")
	logFatalErr(err)
	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "2-collision",
	}
}

func jailRoomMap(gStack *gui.StateStack) MapInfo {
	//exploreState.Map.WriteTile(35, 22, false)
	gMap, err := tilepix.ReadFile("jail.tmx")
	logFatalErr(err)

	boneScript := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)
		boneItemId := 4
		menu_ := gStack.Globals["menu"]
		menuV := reflect.ValueOf(menu_)
		menuI := menuV.Interface().(*InGameMenuState)
		giveBone := func(gMap *GameMap) {
			//player picked up the bone
			gStack.Pop() //remove selection menu
			gStack.PushFitted(x, y, `Found key item: "Calcified bone"`)
			//play sound skeleton_collapsed - pending
			menuI.World.AddKeyItem(boneItemId)
		}

		choices := []string{"Hit space to add it to your Inventory"}
		onSelection := func(index int, c string) {
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
	}

	breakWallScript := func(gameMap *GameMap, entity *Entity, tileX, tileY float64) {
		x, y := gameMap.GetTileIndex(tileX, tileY)
		onPush := func(gMap *GameMap) {
			// The player's pushing the wall.
			gStack.Pop() //remove selection menu
			gStack.PushFitted(x, y, "The wall crumbles.")
			//play sound wall_crumbles - pending

			//see below Triggers - "cracked_stone"
			gMap.RemoveTrigger(35, 22)
			gMap.WriteTile(35, 22, false)
		}
		choices := []string{
			"Push the wall",
			"Get back!",
		}
		onSelection := func(index int, c string) {
			if index == 0 {
				onPush(gameMap)
			}
		}
		gStack.PushSelectionMenu(x, y, 400, 100, "The wall here is crumbling. Push it?", choices, onSelection, false)
		//gStack.PushFITMenu(x, y, "The wall here is crumbling. Push it?", choices, onSelection)
	}

	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
		Actions: map[string]MapAction{
			"break_wall_script": {
				Id:     "RunScript",
				Script: breakWallScript,
			},
			"bone_script": {
				Id:     "RunScript",
				Script: boneScript,
			},
		},
		TriggerTypes: map[string]TriggerType{
			"cracked_stone": {
				OnUse: "break_wall_script",
			},
			"calcified_bone": {
				OnUse: "bone_script",
			},
		},
		Triggers: []TriggerParam{
			{Id: "cracked_stone", X: 35, Y: 22},
			{Id: "calcified_bone", X: 41, Y: 22},
			{Id: "calcified_bone", X: 42, Y: 22},
		},
	}
}

func smallRoomMap(gStack *gui.StateStack) MapInfo {
	gMap, err := tilepix.ReadFile("small_room.tmx")
	logFatalErr(err)
	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     3,
		CollisionLayerName: "collision",
		Actions:            nil,
		TriggerTypes:       nil,
		Triggers:           nil,
	}
}

func logFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
