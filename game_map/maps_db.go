package game_map

import (
	"github.com/bcvery1/tilepix"
	"log"
)

type MapAction struct {
	Id      string
	Scripts []func(a ...interface{}) func(b ...interface{})
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

var MapsDB map[string]func() MapInfo

func init() {
	MapsDB = make(map[string]func() MapInfo)
	MapsDB["player_room"] = playerHouseMap
	MapsDB["small_room"] = smallRoomMap
	MapsDB["jail_room"] = jailRoomMap
}

//player render rule is we render them with Collision Layer
func playerHouseMap() MapInfo {
	gMap, err := tilepix.ReadFile("sontos_house.tmx")
	logFatalErr(err)
	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "2-collision",
	}
}

func jailRoomMap() MapInfo {
	//exploreState.Map.WriteTile(35, 22, false)

	gMap, err := tilepix.ReadFile("jail.tmx")
	logFatalErr(err)
	return MapInfo{
		Tilemap:            gMap,
		CollisionLayer:     2,
		CollisionLayerName: "02 collision",
	}
}

func smallRoomMap() MapInfo {
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
