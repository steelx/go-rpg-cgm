package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/dice"
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
	var HeroDef = combat.ActorDef{
		Stats: combat.DefaultStats,
		StatGrowth: map[string]func() int{
			"HpMax":        dice.Create("4d50+100"),
			"MpMax":        dice.Create("2d50+100"),
			"Strength":     combat.StatsGrowth.Fast,
			"Speed":        combat.StatsGrowth.Fast,
			"Intelligence": combat.StatsGrowth.Med,
		},
	}

	PrintLevelUp := func(levelUp combat.LevelUp) {
		stats := levelUp.BaseStats

		fmt.Printf("HP:+%v MP:+%v \n", stats["HpMax"], stats["MpMax"])

		fmt.Printf("str:+%v spd:+%v int:+%v \n",
			stats["Strength"],
			stats["Speed"],
			stats["Intelligence"])
		fmt.Println("--^^--")
	}

	ApplyXP := func(char combat.Actor, xp float64) {
		char.AddXP(xp)
		fmt.Println("==XP applied==", char.XP)

		fmt.Println("char.ReadyToLevelUp()", char.ReadyToLevelUp())
		for char.ReadyToLevelUp() {
			levelup := char.CreateLevelUp()
			fmt.Printf("Level Up! (Level %v) \n", char.Level+levelup.Level)
			char.ApplyLevel(levelup)
			PrintLevelUp(levelup)
		}
	}

	hero := combat.ActorCreate(HeroDef)

	ApplyXP(hero, 10001)

	fmt.Println("FINAL")
	fmt.Println("HpMax", hero.Stats.Get("HpMax"))
	fmt.Println("HpMax", hero.Stats.Get("MpMax"))
	fmt.Println("Intelligence", hero.Stats.Get("Intelligence"))
	fmt.Println("Speed", hero.Stats.Get("Speed"))

	pixelgl.Run(run)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup(win *pixelgl.Window) {
	stack = gui.StateStackCreate(win)

	var introScene = []interface{}{
		game_map.BlackScreen("blackscreen"),
		game_map.Wait(1),
		game_map.KillState("blackscreen"),
		game_map.PlayBGSound("../sound/rain.mp3"),
		game_map.TitleCaptionScreen("title", "Chandragupta Maurya", 3),
		game_map.SubTitleCaptionScreen("subtitle", "A jRPG game in GO", 2),
		//game_map.Wait(2),
		game_map.KillState("title"),
		game_map.KillState("subtitle"),
		game_map.Scene("map_player_house", true, win),
		game_map.RunActionAddNPC("map_player_house", "sleeper", 14, 19, 3),
		game_map.RunActionAddNPC("map_player_house", "guard", 19, 23, 0),
		game_map.PlaySound("../sound/door_break.mp3", 1),
		game_map.MoveNPC("guard", "map_player_house", []string{
			"up", "up", "up", "left", "left", "left",
		}),
		game_map.Say("map_player_house", "guard", "You'r coming with me!!", 3),
		game_map.StopBGSound(),
		game_map.PlaySound("../sound/wagon.mp3", 4),
		game_map.BlackScreen("blackscreen"),
		game_map.Wait(3),
		game_map.KillState("blackscreen"),
		game_map.ReplaceScene("map_player_house", "map_jail", 31, 21, false, win),
		game_map.Wait(1),
		game_map.Say("map_jail", "hero", "Where am I...", 2),
		game_map.Say("map_jail", "hero", "I should keep looking for ways out", 2),
		game_map.Wait(1),
		game_map.HandOffToMainStack("map_jail"),
	}

	var storyboardI = game_map.StoryboardCreate(stack, win, introScene, false)
	stack.PushFitted(200, 1300, "storyboardI stack pop out.. :)")
	stack.Push(storyboardI)
	stack.Push(gui.TitleScreenCreate(stack, win))

}

//=============================================================
// Game loop
//=============================================================
func gameLoop(win *pixelgl.Window) {
	last := time.Now()
	menu := game_map.InGameMenuStateCreate(stack, win)
	stack.Globals["menu"] = menu

	//set fullscreen
	//win.SetMonitor(globals.Global.PrimaryMonitor)

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
