package main

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_stacks"
	"sort"
	"time"
)

const camZoom = 1.0

var (
	CastleRoomMap = &game_map.GameMap{}
	camPos        = pixel.ZV
	//camSpeed    = 1000.0
	//camZoomSpeed = 1.2
	frameRate  = 15 * time.Millisecond
	textStacks state_stacks.StateStack
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
	setup()
	globals.PrintMemoryUsage()
	gameLoop(win)
}

func main() {
	pixelgl.Run(run)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup() {
	textStacks = state_stacks.StateStackCreate()

	// Init map
	m, err := tilepix.ReadFile("small_room.tmx")
	globals.PanicIfErr(err)
	CastleRoomMap.Create(m)

	//Actions & Triggers
	gUpDoorTeleport := ActionTeleport(*CastleRoomMap, globals.Direction{7, 2})
	gDownDoorTeleport := ActionTeleport(*CastleRoomMap, globals.Direction{9, 10})
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
			textStacks.AddFitted(300, 250, "Dude, snakes.. run!")
		},
	)

	CastleRoomMap.SetTrigger(7, 2, gTriggerTop)
	CastleRoomMap.SetTrigger(9, 10, gTriggerBottom)
	CastleRoomMap.SetTrigger(8, 6, gTriggerFlowerPot)

	//GameMap.GetEntityAtPos needs this
	CastleRoomMap.Entities = []*game_map.Entity{Hero.Entity, NPC2.Entity, NPC1.Entity}
}

//=============================================================
// Game loop
//=============================================================
func gameLoop(win *pixelgl.Window) {
	last := time.Now()

	choices := []string{"Menu 1", "lola", "Menu 2", "Menu 03", "Menu 04", "Menu 05", "Menu 06", "Menu 007", "", "", "", "Menu @_@"}
	textStacks.AddSelectionMenu(
		-100, 250, 400, 200,
		"Select from the list below",
		choices, func(i int, item string) {
			fmt.Println(i, item)
		})

	textStacks.AddFixed(
		-150, 10, 300, 100,
		"A nation can survive its fools, and even the ambitious. But it cannot survive treason from within. An enemy at the gates is less formidable, for he is known and carries his banner openly. But the traitor moves amongst those within the gate freely, his sly whispers rustling through all the alleys, heard in the very halls of government itself. For the traitor appears not a traitor; he speaks in accents familiar to his victims, and he wears their face and their arguments, he appeals to the baseness that lies deep in the hearts of all men. He rots the soul of a nation, he works secretly and unknown in the night to undermine the pillars of the city, he infects the body politic so that it can no longer resist. A murderer is less to fear. Jai Hind I Love India <3 ",
		"Ajinkya", globals.AvatarPng)

	textStacks.AddFitted(100, 100, "Hello! if you smell the rock was cookin")
	textStacks.AddFitted(200, 200, "1111 if you smell the rock was cookin")
	textStacks.AddFitted(300, 250, "Pop pop pop. mark me unread HIT spacebar")

	progressBar := gui.ProgressBarCreate(200, 0)
	//progressBar.SetValue(90)

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

			err := CastleRoomMap.DrawAfter(func(canvas *pixelgl.Canvas, layer int) {
				gameCharacters := [3]game_map.Character{*Hero, *NPC2, *NPC1}

				sort.Slice(gameCharacters[:], func(i, j int) bool {
					return gameCharacters[i].Entity.TileY < gameCharacters[j].Entity.TileY
				})

				if layer == 2 {
					for _, gCharacter := range gameCharacters {
						gCharacter.Entity.TeleportAndDraw(*CastleRoomMap, canvas)
						gCharacter.Controller.Update(dt)
					}
				}
			})
			globals.PanicIfErr(err)

			textStacks.Render(win)
			textStacks.Update(dt)

			progressBar.Render(win)

			// Camera
			CastleRoomMap.CamToTile(Hero.Entity.TileX, Hero.Entity.TileY)
			camPos = pixel.V(CastleRoomMap.CamX, CastleRoomMap.CamY)
			cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
			win.SetMatrix(cam)

			if win.JustPressed(pixelgl.KeyE) {
				tileX, tileY := Hero.Entity.Map.GetTileIndex(Hero.GetFacedTileCoords())
				trigger := Hero.Entity.Map.GetTrigger(tileX, tileY)
				if trigger.OnUse != nil {
					trigger.OnUse(Hero.Entity)
				}
			}
		}

		win.Update()
	}
}
