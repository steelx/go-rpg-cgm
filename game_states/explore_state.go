package game_states

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/gui"
	"log"
	"sort"
)

type ExploreState struct {
	Stack       *gui.StateStack
	MapDef      *tilepix.Map
	Map         *game_map.GameMap
	Hero        *game_map.Character
	win         *pixelgl.Window
	startPos    pixel.Vec
	heroVisible bool
}

func ExploreStateCreate(stack *gui.StateStack, mapInfo game_map.MapInfo, window *pixelgl.Window) ExploreState {

	es := ExploreState{
		Stack:       stack,
		MapDef:      mapInfo.Tilemap,
		heroVisible: true,
	}

	es.win = window
	es.Map = game_map.MapCreate(mapInfo)

	es.Hero = character_states.Characters["hero"](es.Map)
	es.Map.NPCbyId[es.Hero.Id] = es.Hero
	//es.Hero.Controller.Change("wait", globals.Direction{0, 0})

	es.startPos = pixel.V(es.Hero.Entity.TileX, es.Hero.Entity.TileY)
	es.Map.GoToTile(es.startPos.X, es.startPos.Y)

	return es
}

func (es *ExploreState) HideHero() {
	es.heroVisible = false
	es.Hero.Entity.TileX = 0
	es.Hero.Entity.TileY = 0
}
func (es *ExploreState) ShowHero(tileX, tileY float64) {
	es.heroVisible = true
	es.Hero.Entity.TileX = tileX
	es.Hero.Entity.TileY = tileY
	es.Map.GoToTile(es.Hero.Entity.TileX, es.Hero.Entity.TileY)
}

func (es ExploreState) Enter() {

}
func (es ExploreState) Exit() {

}

func (es *ExploreState) Update(dt float64) bool {
	// Update the camera according to player position
	playerPosX, playerPosY := es.Hero.Entity.TileX, es.Hero.Entity.TileY
	es.Map.GoToTile(playerPosX, playerPosY)

	gameCharacters := append([]*game_map.Character{es.Hero}, es.Map.NPCs...)
	for _, gCharacter := range gameCharacters {
		gCharacter.Controller.Update(dt)
	}
	return true
}

func (es ExploreState) Render(win *pixelgl.Window) {
	//Map & Characters
	err := es.Map.DrawAfter(func(canvas *pixelgl.Canvas, layer int) {
		var gameCharacters []*game_map.Character
		gameCharacters = append(gameCharacters, es.Map.NPCs...)
		if es.heroVisible {
			gameCharacters = append([]*game_map.Character{es.Hero}, gameCharacters...)
		}
		//sort players as per visible to screen Y position
		sort.Slice(gameCharacters[:], func(i, j int) bool {
			return gameCharacters[i].Entity.TileY < gameCharacters[j].Entity.TileY
		})

		if layer == es.Map.MapInfo.CollisionLayer {
			for _, gCharacter := range gameCharacters {
				//gCharacter.Entity.TeleportAndDraw(es.Map, canvas) //probably can remove now
				gCharacter.Entity.Render(es.Map, canvas)
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	//Camera
	camPos := pixel.V(es.Map.CamX, es.Map.CamY)
	cam := pixel.IM.Scaled(camPos, 1.0).Moved(es.win.Bounds().Center().Sub(camPos))
	win.SetMatrix(cam)
}

func (es ExploreState) HandleInput(win *pixelgl.Window) {
	//use key
	if win.JustPressed(pixelgl.KeyE) {
		// which way is the player facing?
		//tileX, tileY := es.Map.GetTileIndex(es.Hero.GetFacedTileCoords())
		tileX, tileY := es.Hero.GetFacedTileCoords()
		trigger := es.Map.GetTrigger(tileX, tileY)
		if trigger.OnUse != nil {
			trigger.OnUse(es.Map, es.Hero.Entity, tileX, tileY)
		}
	}
}

func (es *ExploreState) AddNPC(NPC *game_map.Character) {
	es.Map.AddNPC(NPC)
	//NPC.Controller.Change("wait", globals.Direction{0, 0})
}
