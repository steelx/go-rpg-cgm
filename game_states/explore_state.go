package game_states

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"log"
	"sort"
)

type ExploreState struct {
	Stack    *gui.StateStack
	MapDef   *tilepix.Map
	Map      *game_map.GameMap
	Hero     *game_map.Character
	win      *pixelgl.Window
	startPos pixel.Vec
}

func ExploreStateCreate(stack *gui.StateStack,
	tilemap *tilepix.Map, collisionLayer int, collisionLayerName string,
	startPos pixel.Vec, heroPng pixel.Picture, window *pixelgl.Window) ExploreState {

	es := ExploreState{
		Stack:    stack,
		startPos: startPos,
		MapDef:   tilemap,
	}

	es.win = window
	es.Map = game_map.MapCreate(es.MapDef, collisionLayer, collisionLayerName)

	es.Hero = game_map.CharacterCreate("Ajinkya",
		[][]int{{16, 17, 18, 19}, {20, 21, 22, 23}, {24, 25, 26, 27}, {28, 29, 30, 31}},
		game_map.CharacterFacingDirection[2],
		game_map.CharacterDefinition{
			Texture: heroPng, Width: 16, Height: 24,
			StartFrame: 24,
			TileX:      startPos.X,
			TileY:      startPos.Y,
			Map:        es.Map,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return character_states.WaitStateCreate(es.Hero, es.Map)
			},
			"move": func() state_machine.State {
				return character_states.MoveStateCreate(es.Hero, es.Map)
			},
		},
	)
	es.Hero.Controller.Change("wait", globals.Direction{0, 0})

	es.Map.GoToTile(startPos.X, startPos.Y)

	return es
}

func (es *ExploreState) HideHero() {
	es.Hero.Entity.TileX = -1
	es.Hero.Entity.TileY = -1
}
func (es *ExploreState) ShowHero() {
	es.Hero.Entity.TileX = es.startPos.X
	es.Hero.Entity.TileY = es.startPos.Y
	es.Map.GoToTile(es.startPos.X, es.startPos.Y)
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
		gameCharacters := append([]*game_map.Character{es.Hero}, es.Map.NPCs...)
		//sort players as per visible to screen Y position
		sort.Slice(gameCharacters[:], func(i, j int) bool {
			return gameCharacters[i].Entity.TileY < gameCharacters[j].Entity.TileY
		})

		if layer == 2 {
			for _, gCharacter := range gameCharacters {
				gCharacter.Entity.TeleportAndDraw(*es.Map, canvas)
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
	//es.Hero.Controller.Update(globals.Global.DeltaTime)
	//use key
	if win.JustPressed(pixelgl.KeyE) {
		// which way is the player facing?
		tileX, tileY := es.Hero.Entity.Map.GetTileIndex(es.Hero.GetFacedTileCoords())
		trigger := es.Hero.Entity.Map.GetTrigger(tileX, tileY)
		if trigger.OnUse != nil {
			trigger.OnUse(es.Hero.Entity)
		}
	}
}

func (es *ExploreState) AddNPC(NPC *game_map.Character) {
	es.Map.AddNPC(NPC)
	NPC.Controller.Change("wait", globals.Direction{0, 0})
}
