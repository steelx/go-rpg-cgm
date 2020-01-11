package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/game_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/storyboard"
	"time"
)

var (
	stack        *gui.StateStack
	exploreState game_states.ExploreState
	//camSpeed    = 1000.0
	//camZoomSpeed = 1.2
	frameRate = 15 * time.Millisecond
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "GP RPG",
		Bounds:      pixel.R(0, 0, globals.Global.WindowWidth, globals.Global.WindowHeight),
		VSync:       globals.Global.Vsync,
		Undecorated: globals.Global.Undecorated,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	globals.Global.Win = win

	globals.PrintMemoryUsage()
	// Setup world etc.
	setup(win)
	globals.PrintMemoryUsage()
	gameLoop(win)
}

func main() {
	pixelgl.Run(run)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup(win *pixelgl.Window) {
	stack = gui.StateStackCreate(win)

	// Init map & Add Stacks
	walkCyclePng, err := globals.LoadPicture("../resources/walk_cycle.png")
	globals.PanicIfErr(err)
	exploreState = game_states.ExploreStateCreate(stack,
		globals.CastleMapDef, pixel.V(2, 4), walkCyclePng, win,
	)

	//Add NPCs
	var NPC1 *game_map.Character
	NPC1 = game_map.CharacterCreate(
		"Aghori Baba", nil, game_map.CharacterFacingDirection[2],
		game_map.CharacterDefinition{
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 46,
			TileX:      9,
			TileY:      4,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return character_states.NPCWaitStateCreate(NPC1, exploreState.Map)
			},
		},
	)

	var NPC2 *game_map.Character
	NPC2 = game_map.CharacterCreate(
		"Bhadrasaal", [][]int{{48, 49, 50, 51}, {52, 53, 54, 55}, {56, 57, 58, 59}, {60, 61, 62, 63}},
		game_map.CharacterFacingDirection[2],
		game_map.CharacterDefinition{
			Texture: walkCyclePng, Width: 16, Height: 24,
			StartFrame: 56,
			TileX:      3,
			TileY:      8,
		},
		map[string]func() state_machine.State{
			"wait": func() state_machine.State {
				return character_states.NPCStrollWaitStateCreate(NPC2, exploreState.Map)
			},
			"move": func() state_machine.State {
				return character_states.MoveStateCreate(NPC2, exploreState.Map)
			},
		},
	)

	exploreState.AddNPC(NPC1)
	exploreState.AddNPC(NPC2)

	stack.Push(&exploreState)

	var introScene = []interface{}{
		storyboard.Wait(0),
		storyboard.BlackScreen("blackscreen"),
		storyboard.Wait(2),
	}

	var storyboardI = storyboard.Create(stack, win, introScene)
	stack.Push(&storyboardI)
}

//=============================================================
// Game loop
//=============================================================
func gameLoop(win *pixelgl.Window) {
	last := time.Now()

	//initial map Camera
	exploreState.Map.GoToTile(4, 4)
	camPos := pixel.V(exploreState.Map.CamX, exploreState.Map.CamY)
	cam := pixel.IM.Scaled(camPos, 1.0).Moved(win.Bounds().Center().Sub(camPos))
	win.SetMatrix(cam)

	tick := time.Tick(frameRate)
	for !win.Closed() {

		if win.JustPressed(pixelgl.KeyQ) {
			break
		}

		win.Clear(globals.Global.ClearColor)

		select {
		case <-tick:
			dt := time.Since(last).Seconds()
			last = time.Now()

			//update StateStack
			stack.Update(dt)
			stack.Render(win)

			//Fullscreen Layout Menu
			if win.JustPressed(pixelgl.KeyLeftAlt) {
				menu := game_states.InGameMenuStateCreate(stack, win)
				stack.Push(menu)
			}
		}

		win.Update()
	}
}
