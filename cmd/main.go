package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/utilz"
	"time"
)

var (
	stack *gui.StateStack
	//exploreState game_states.ExploreState
	//camSpeed    = 1000.0
	//camZoomSpeed = 1.2
	frameRate = 15 * time.Millisecond
	gWorld    *combat.WorldExtended
)

func run() {
	globals.Global.PrimaryMonitor = pixelgl.PrimaryMonitor()
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

	utilz.PrintMemoryUsage()
	// Setup world etc.
	setup(win)
	utilz.PrintMemoryUsage()
	gameLoop(win)
}

func main() {
	pixelgl.Run(run)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup(win *pixelgl.Window) {
	//set fullscreen
	//win.SetMonitor(globals.Global.PrimaryMonitor)

	stack = gui.StateStackCreate(win)

	gWorld = combat.WorldExtendedCreate()
	gWorld.Party.Add(combat.ActorCreate(combat.HeroDef))
	gWorld.Party.Add(combat.ActorCreate(combat.ThiefDef))
	gWorld.AddItem(11, 3)
	gWorld.AddItem(12, 2)

	stack.Globals["world"] = gWorld

	var storyboardI = game_map.StoryboardCreate(stack, win, game_map.IntroScene, false)
	stack.Push(storyboardI)

	stack.Push(gui.TitleScreenCreate(stack, win))

}

//=============================================================
// Game loop
//=============================================================
func gameLoop(win *pixelgl.Window) {
	last := time.Now()

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
			stack.DeltaTime = dt

			//update StateStack
			stack.Update(dt)
			gWorld.Update(dt)

			//<-- this would render only 1 stack at a time
			switch top := stack.States[stack.GetLastIndex()].(type) {
			case *game_map.InGameMenuState:
				top.Render(win)

			case *gui.TitleScreen:
				top.Render(win)

			default:
				stack.Render(win) //else render all
			}

		}

		win.Update()
	}
}
