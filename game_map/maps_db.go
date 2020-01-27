package game_map

import (
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/tilepix"
	"log"
)

var MapsDB map[string]func(gStack *gui.StateStack) MapInfo

func init() {
	MapsDB = make(map[string]func(gStack *gui.StateStack) MapInfo)
	MapsDB["map_player_house"] = mapPlayerHouse
	MapsDB["small_room"] = smallRoomMap
	MapsDB["map_jail"] = mapJail
	MapsDB["map_sewer"] = mapSewer
	MapsDB["map_arena"] = mapArena
}

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
	Tilemap                         *tilepix.Map
	CollisionLayer                  int
	CollisionLayerName, HiddenLayer string
	Actions                         map[string]MapAction   //"break_wall_script" : { Id = "RunScript", Scripts : []{ CrumbleScript } }
	TriggerTypes                    map[string]TriggerType //"cracked_stone" : { OnUse = "break_wall_script" }
	Triggers                        []TriggerParam         //[]{Id = "cracked_stone", x = 60, y = 11}
	OnWake                          map[string]TriggerParam
}

//player render rule is we render them with Collision Layer
func mapPlayerHouse(gStack *gui.StateStack) MapInfo {
	//gStack could be global stack in future

	gMap, err := tilepix.ReadFile("sontos_house.tmx")
	logFatalErr(err)
	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "2-collision",
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
