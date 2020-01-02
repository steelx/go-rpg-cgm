package main

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"time"
)

var (
	basicAtlas    = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	endTxt        = &text.Text{}
	CastleRoomMap = &GameMap{}
	camPos        = pixel.ZV
	//camSpeed    = 1000.0
	camZoom = 1.8
	//camZoomSpeed = 1.2
	frameRate = 15 * time.Millisecond
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "GP RPG",
		Bounds:      pixel.R(0, 0, float64(global.gWindowWidth), float64(global.gWindowHeight)),
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
func main() {
	pixelgl.Run(run)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup() {
	// Camera setup

	// Init map
	m, err := tilepix.ReadFile("small_room.tmx")
	panicIfErr(err)

	CastleRoomMap.Create(m)
	CastleRoomMap.CamToTile(11, 5)

	// Camera
	camPos = pixel.V(CastleRoomMap.mCamX, CastleRoomMap.mCamY)
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
	global.gWin.SetMatrix(cam)

	//Actions & Triggers
	gUpDoorTeleport := ActionTeleport(*CastleRoomMap, 7, 2)
	gDownDoorTeleport := ActionTeleport(*CastleRoomMap, 9, 10)
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
			endTxt = text.New(pixel.V(220, 100), basicAtlas)
			fmt.Fprintln(endTxt, "Pot is full of snakes!")
		},
	)

	tileX, tileY := CastleRoomMap.GetTileIndex(7, 2)
	CastleRoomMap.SetTrigger(tileX, tileY, gTriggerTop)

	tileX, tileY = CastleRoomMap.GetTileIndex(9, 10)
	CastleRoomMap.SetTrigger(tileX, tileY, gTriggerBottom)

	tileX, tileY = CastleRoomMap.GetTileIndex(8, 6)
	CastleRoomMap.SetTrigger(tileX, tileY, gTriggerFlowerPot)

	// gHero Init
	gHero.mController.Change("wait", Direction{0, 0})
}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()

	tick := time.Tick(frameRate)
	for !global.gWin.Closed() {

		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}
		if global.gWin.Pressed(pixelgl.KeySpace) {
			tileX, tileY := CastleRoomMap.GetTileIndex(gHero.GetFacedTileCoords())
			trigger := CastleRoomMap.GetTrigger(tileX, tileY)
			if trigger.OnUse != nil {
				trigger.OnUse(gHero.mEntity)
			}
		}

		global.gWin.Clear(global.gClearColor)

		select {
		case <-tick:
			dt := time.Since(last).Seconds()
			last = time.Now()
			CastleRoomMap.DrawAfter(1, func(canvas *pixelgl.Canvas) {
				gHero.mEntity.TeleportAndDraw(*CastleRoomMap, canvas)
			})
			endTxt.Draw(global.gWin, pixel.IM.Scaled(pixel.V(300, 300), 1))
			gHero.mController.Update(dt)
		}

		global.gWin.Update()
	}
}
