package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
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

	//stack.Push(&exploreState)

	var introScene = []interface{}{
		storyboard.Scene("sontos_house", 20, 20, false, win),
		storyboard.Wait(0),
		storyboard.BlackScreen("blackscreen"),
		storyboard.Wait(1),
		storyboard.FadeScreen("fadeWhite", 1, 0, 2),
		storyboard.TitleCaptionScreen("title", "Chandragupta Maurya", 3),
		storyboard.SubTitleCaptionScreen("subtitle", "A jRPG game in GO", 2),
		storyboard.Wait(2),
	}

	var storyboardI = storyboard.Create(stack, win, introScene)
	stack.Push(storyboardI)

}

//=============================================================
// Game loop
//=============================================================
func gameLoop(win *pixelgl.Window) {
	last := time.Now()

	//initial map Camera

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
			globals.Global.DeltaTime = dt

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
