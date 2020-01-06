package main

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"sort"
	"time"
)

const camZoom = 1.0

var (
	basicAtlas14  *text.Atlas
	basicAtlas12  *text.Atlas
	CastleRoomMap = &GameMap{}
	camPos        = pixel.ZV
	//camSpeed    = 1000.0
	//camZoomSpeed = 1.2
	frameRate                              = 15 * time.Millisecond
	avatarPng, continueCaretPng, cursorPng pixel.Picture
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "GP RPG",
		Bounds:      pixel.R(0, 0, global.gWindowWidth, global.gWindowHeight),
		VSync:       global.gVsync,
		Undecorated: global.gUndecorated,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	global.gWin = win

	PrintMemoryUsage()
	// Setup world etc.
	setup()
	PrintMemoryUsage()
	gameLoop()
}
func init() {
	fontFace14, err := loadTTF("./resources/font/joystix.ttf", 14)
	panicIfErr(err)
	fontFace12, err := loadTTF("./resources/font/joystix.ttf", 12)
	panicIfErr(err)
	basicAtlas14 = text.NewAtlas(fontFace14, text.ASCII)
	basicAtlas12 = text.NewAtlas(fontFace12, text.ASCII)

	//images for Textbox & Panel
	avatarPng, err = LoadPicture("./resources/avatar.png")
	panicIfErr(err)
	continueCaretPng, err = LoadPicture("./resources/continue_caret.png")
	panicIfErr(err)
	cursorPng, err = LoadPicture("./resources/cursor.png")
	panicIfErr(err)
}
func main() {
	pixelgl.Run(run)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup() {
	// Init map
	m, err := tilepix.ReadFile("small_room.tmx")
	panicIfErr(err)
	CastleRoomMap.Create(m)

	//Actions & Triggers
	gUpDoorTeleport := ActionTeleport(*CastleRoomMap, Direction{7, 2})
	gDownDoorTeleport := ActionTeleport(*CastleRoomMap, Direction{9, 10})
	gTriggerTop := TriggerCreate(gDownDoorTeleport, nil, nil)
	gTriggerBottom := TriggerCreate(
		gUpDoorTeleport,
		nil,
		nil,
	)
	gTriggerFlowerPot := TriggerCreate(
		nil,
		nil,
		func(entity *Entity) {
			//story01a.gMap = entity.gMap
			//story01a.Render()
			//story01a = story01a.Play("space")

		},
	)

	CastleRoomMap.SetTrigger(7, 2, gTriggerTop)
	CastleRoomMap.SetTrigger(9, 10, gTriggerBottom)
	CastleRoomMap.SetTrigger(8, 6, gTriggerFlowerPot)

	CastleRoomMap.mEntities = []*Entity{gHero.mEntity, gNPC2.mEntity, gNPC1.mEntity}
}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()

	pic, err := LoadPicture("./resources/simple_panel.png")
	panicIfErr(err)

	menu := SelectionMenuCreate([]string{"Menu 1-", "", "Menu 2", "Menu 03", "Menu 007"}, 1, pixel.V(200, 100), func(i int, item string) {
		fmt.Println(i, item)
	})

	tBox := TextboxCreate(
		"A nation can survive its fools, and even the ambitious. But it cannot survive treason from within. An enemy at the gates is less formidable, for he is known and carries his banner openly. But the traitor moves amongst those within the gate freely, his sly whispers rustling through all the alleys, heard in the very halls of government itself. For the traitor appears not a traitor; he speaks in accents familiar to his victims, and he wears their face and their arguments, he appeals to the baseness that lies deep in the hearts of all men. He rots the soul of a nation, he works secretly and unknown in the night to undermine the pillars of the city, he infects the body politic so that it can no longer resist. A murderer is less to fear. Jai Hind I Love India <3 ",
		basicAtlas12,
		PanelCreate(pic, pixel.V(-150, 200), 300, 100),
		continueCaretPng,
		"Ajinkya",
		avatarPng,
	)

	tick := time.Tick(frameRate)
	for !global.gWin.Closed() {

		if global.gWin.JustPressed(pixelgl.KeyQ) {
			break
		}
		if global.gWin.JustPressed(pixelgl.KeySpace) {
			tBox.Next()
		}

		global.gWin.Clear(global.gClearColor)

		select {
		case <-tick:
			dt := time.Since(last).Seconds()
			last = time.Now()

			err := CastleRoomMap.DrawAfter(func(canvas *pixelgl.Canvas, layer int) {
				gameCharacters := [3]Character{*gHero, *gNPC2, *gNPC1}

				sort.Slice(gameCharacters[:], func(i, j int) bool {
					return gameCharacters[i].mEntity.mTileY < gameCharacters[j].mEntity.mTileY
				})

				if layer == 2 {
					for _, gCharacter := range gameCharacters {
						gCharacter.mEntity.TeleportAndDraw(*CastleRoomMap, canvas)
						gCharacter.mController.Update(dt)
					}
				}
			})
			panicIfErr(err)

			tBox.DrawTextWithPanel()
			menu.Render()
			menu.HandleInput()

			// Camera
			CastleRoomMap.CamToTile(gHero.mEntity.mTileX, gHero.mEntity.mTileY)
			camPos = pixel.V(CastleRoomMap.mCamX, CastleRoomMap.mCamY)
			cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
			global.gWin.SetMatrix(cam)

			if global.gWin.JustPressed(pixelgl.KeyE) {
				tileX, tileY := gHero.mEntity.gMap.GetTileIndex(gHero.GetFacedTileCoords())
				trigger := gHero.mEntity.gMap.GetTrigger(tileX, tileY)
				if trigger.OnUse != nil {
					trigger.OnUse(gHero.mEntity)
				}
			}
		}

		global.gWin.Update()
	}
}
