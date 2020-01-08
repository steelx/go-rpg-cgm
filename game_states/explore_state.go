package game_states

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/state_stacks"
	"log"
	"math"
	"sort"
)

type ExploreState struct {
	Stack  *state_stacks.StateStack
	MapDef *tilepix.Map
	Map    *game_map.GameMap
	Hero   *game_map.Character
	win    *pixelgl.Window
}

func ExploreStateCreate(
	stack *state_stacks.StateStack, tilemap *tilepix.Map, startPos pixel.Vec, heroPng pixel.Picture, window *pixelgl.Window) ExploreState {
	es := ExploreState{
		Stack:  stack,
		MapDef: tilemap,
	}

	es.win = window
	es.Map = game_map.MapCreate(es.MapDef)

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

func (es ExploreState) Enter(data globals.Direction) {

}
func (es ExploreState) Exit() {

}

func (es *ExploreState) Update(dt float64) {
	// Update the camera according to player position
	playerPosX, playerPosY := es.Hero.Entity.TileX, es.Hero.Entity.TileY
	es.Map.CamX = math.Floor(playerPosX)
	es.Map.CamY = math.Floor(playerPosY)

	gameCharacters := append([]*game_map.Character{es.Hero}, es.Map.NPCs...)
	for _, gCharacter := range gameCharacters {
		gCharacter.Controller.Update(dt)
	}
}

func (es ExploreState) Render() {
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
	es.win.SetMatrix(cam)
}
func (es ExploreState) HandleInput(window *pixelgl.Window) {
	if window.JustPressed(pixelgl.KeyE) {
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
