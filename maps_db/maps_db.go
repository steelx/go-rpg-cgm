package maps_db

import (
	"github.com/bcvery1/tilepix"
	"log"
)

type MapDef struct {
	gMap               *tilepix.Map
	collisionLayer     int
	collisionLayerName string
}

var MapsDB map[string]func() (*tilepix.Map, int, string)

func init() {
	MapsDB = make(map[string]func() (*tilepix.Map, int, string))
	MapsDB["player_room"] = playerHouseMap
	MapsDB["small_room"] = smallRoomMap
	MapsDB["jail_room"] = jailRoomMap
}

//player render rule is we render them with Collision Layer
func playerHouseMap() (gMap *tilepix.Map, collisionLayer int, collisionLayerName string) {
	collisionLayer, collisionLayerName = 2, "2-collision"
	gMap, err := tilepix.ReadFile("sontos_house.tmx")
	logFatalErr(err)
	return
}

func jailRoomMap() (gMap *tilepix.Map, collisionLayer int, collisionLayerName string) {
	collisionLayer, collisionLayerName = 2, "02 collision"
	gMap, err := tilepix.ReadFile("jail.tmx")
	logFatalErr(err)
	return
}

func smallRoomMap() (gMap *tilepix.Map, collisionLayer int, collisionLayerName string) {
	collisionLayer, collisionLayerName = 3, "collision"
	gMap, err := tilepix.ReadFile("small_room.tmx")
	logFatalErr(err)
	return
}

func logFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
