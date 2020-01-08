package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/game_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"time"
)

const camZoom = 1.0

var (
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
	// Init map
	choices := []string{"Menu 1", "lola", "Menu 2", "Menu 03", "Menu 04", "Menu 05", "Menu 06", "Menu 007", "", "", "", "Menu @_@"}
	//textStacks = gui.StateStackCreate()

	// Init map
	walkCyclePng, err := globals.LoadPicture("../resources/walk_cycle.png")
	globals.PanicIfErr(err)
	exploreState = game_states.ExploreStateCreate(
		globals.CastleMapDef, pixel.V(2, 4), walkCyclePng, win,
	)

	exploreState.Stack.PushSelectionMenu(
		-100, 250, 400, 200,
		"Select from the list below",
		choices, func(i int, item string) {
			fmt.Println(i, item)
		})

	exploreState.Stack.PushFixed(
		-150, 10, 300, 100,
		"A nation can survive its fools, and even the ambitious. But it cannot survive treason from within. An enemy at the gates is less formidable, for he is known and carries his banner openly. But the traitor moves amongst those within the gate freely, his sly whispers rustling through all the alleys, heard in the very halls of government itself. For the traitor appears not a traitor; he speaks in accents familiar to his victims, and he wears their face and their arguments, he appeals to the baseness that lies deep in the hearts of all men. He rots the soul of a nation, he works secretly and unknown in the night to undermine the pillars of the city, he infects the body politic so that it can no longer resist. A murderer is less to fear. Jai Hind I Love India <3 ",
		"Ajinkya", globals.AvatarPng)

	exploreState.Stack.PushFitted(100, 100, "Hello! if you smell the rock was cookin")
	exploreState.Stack.PushFitted(200, 200, "1111 if you smell the rock was cookin")
	exploreState.Stack.PushFitted(300, 250, "Pop pop pop. mark me unread HIT spacebar")
	exploreState.Stack.Push(gui.ProgressBarCreate(exploreState.Stack, 200, -50))

	//Actions & Triggers
	gUpDoorTeleport := ActionTeleport(*exploreState.Map, globals.Direction{7, 2})
	gDownDoorTeleport := ActionTeleport(*exploreState.Map, globals.Direction{9, 10})
	gTriggerTop := game_map.TriggerCreate(gDownDoorTeleport, nil, nil)
	gTriggerBottom := game_map.TriggerCreate(
		gUpDoorTeleport,
		nil,
		nil,
	)
	gTriggerFlowerPot := game_map.TriggerCreate(
		nil,
		nil,
		func(entity *game_map.Entity) {
			exploreState.Stack.PushFitted(300, 250, "Dude, snakes.. run!")
		},
	)

	exploreState.Map.SetTrigger(7, 2, gTriggerTop)
	exploreState.Map.SetTrigger(9, 10, gTriggerBottom)
	exploreState.Map.SetTrigger(8, 6, gTriggerFlowerPot)

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

			exploreState.Update(dt)
			exploreState.HandleInput(win)
			exploreState.Render()

			exploreState.Stack.Render(win)
			exploreState.Stack.Update(dt)

			if win.JustPressed(pixelgl.KeyF) {
				fade := gui.FadeScreenCreate(exploreState.Stack, 1, 0, 3, pixel.V(exploreState.Map.CamX, exploreState.Map.CamY))
				exploreState.Stack.Push(&fade)
			}
		}

		win.Update()
	}
}
