package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"time"
)

var (
	stack *gui.StateStack
	//exploreState game_states.ExploreState
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

	//player_room, collision, collisionLayerName := maps_db.MapsDB["player_room"]()
	//exploreState = game_states.ExploreStateCreate(stack, player_room, collision, collisionLayerName, win)
	//
	////Add NPCs
	//exploreState.AddNPC(character_states.NPC1(exploreState.Map))
	//exploreState.AddNPC(character_states.NPC2(exploreState.Map))
	//stack.Push(&exploreState)

	var introScene = []interface{}{
		//game_map.BlackScreen("blackscreen"),
		//game_map.Wait(1),
		//game_map.KillState("blackscreen"),
		//game_map.TitleCaptionScreen("title", "Chandragupta Maurya", 3),
		//game_map.SubTitleCaptionScreen("subtitle", "A jRPG game in GO", 2),
		//game_map.Wait(3),
		//game_map.KillState("title"),
		//game_map.KillState("subtitle"),
		game_map.Scene("player_room", true, win),
		game_map.RunActionAddNPC("player_room", "sleeper", 14, 19, 3),
		game_map.RunActionAddNPC("player_room", "guard", 19, 23, 0),
		game_map.Say("player_room", "guard", "..door smashed", 1.5),
		//play sound door_smashed - pending
		game_map.MoveNPC("guard", "player_room", []string{
			"up", "up", "up", "left", "left", "left",
		}),
		game_map.Say("player_room", "guard", "You'r coming with me!!", 3),
		game_map.BlackScreen("blackscreen"),
		game_map.Wait(1),
		game_map.KillState("blackscreen"),
		game_map.ReplaceScene("player_room", "jail_room", 31, 21, false, win),
		game_map.Wait(1),
		game_map.Say("jail_room", "hero", "Where am I...", 1.5),
		game_map.Say("jail_room", "hero", "I should stay calm..", 2.5),
		game_map.Wait(1),
		game_map.HandOffToMainStack("jail_room"),
	}

	var storyboardI = game_map.Create(stack, win, introScene, false)
	stack.PushFitted(200, 1300, "storyboardI stack pop out.. :)")
	stack.Push(storyboardI)

}

//=============================================================
// Game loop
//=============================================================
func gameLoop(win *pixelgl.Window) {
	last := time.Now()
	menu := game_map.InGameMenuStateCreate(stack, win)
	stack.Globals["menu"] = menu

	tick := time.Tick(frameRate)
	for !win.Closed() {

		if win.JustPressed(pixelgl.KeyQ) {
			break
		}
		//Fullscreen Layout Menu
		if win.JustPressed(pixelgl.KeyLeftAlt) {
			//In Game Menu
			stack.Push(menu)
		}

		win.Clear(globals.Global.ClearColor)

		select {
		case <-tick:
			dt := time.Since(last).Seconds()
			last = time.Now()
			globals.Global.DeltaTime = dt

			//update StateStack
			stack.Update(dt)
			//stack.Render(win)

			//<-- this would render only 1 stack at a time
			if len(stack.States) > 1 {
				switch top := stack.States[stack.GetLastIndex()].(type) {
				case *game_map.InGameMenuState:
					top.Render(win)

				default:
					stack.Render(win) //else render all
				}
			} else {
				stack.Render(win)
			}

		}

		win.Update()
	}
}
