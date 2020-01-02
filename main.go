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
	CastleRoomMap = &GameMap{}
	gHero         *Character
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
	CastleRoomMap.CamToTile(5, 6)

	// Camera
	camPos = pixel.V(CastleRoomMap.mCamX, CastleRoomMap.mCamY)
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
	global.gWin.SetMatrix(cam)

	pic, err := LoadPicture("./resources/walk_cycle.png")
	panicIfErr(err)

	//quickTeleport := ActionTeleport(*CastleRoomMap, 4, 3)
	gUpDoorTeleport := ActionTeleport(*CastleRoomMap, 7, 2)
	gDownDoorTeleport := ActionTeleport(*CastleRoomMap, 9, 10)
	gTriggerTop := TriggerCreate(gDownDoorTeleport, nil, nil)
	gTriggerBottom := TriggerCreate(
		gUpDoorTeleport,
		func() {
			tileX, tileY := CastleRoomMap.GetTileIndex(9, 10)
			fmt.Println(tileX, tileY)
			//endTxt := text.New(pixel.V(tileX, tileY), basicAtlas)
			//fmt.Fprintln(endTxt, "THE END.")
			//endTxt.Draw(CastleRoomMap.canvas, pixel.IM.Scaled(endTxt.Bounds().Center(), 2))
		},
		nil,
	)
	gTriggerFlowerPot := TriggerCreate(
		nil,
		nil,
		gDownDoorTeleport,
	)

	tileX, tileY := CastleRoomMap.GetTileIndex(7, 2)
	CastleRoomMap.SetTrigger(tileX, tileY, gTriggerTop)

	tileX, tileY = CastleRoomMap.GetTileIndex(9, 10)
	CastleRoomMap.SetTrigger(tileX, tileY, gTriggerBottom)

	tileX, tileY = CastleRoomMap.GetTileIndex(8, 6)
	CastleRoomMap.SetTrigger(tileX, tileY, gTriggerFlowerPot)

	gHero = &Character{
		mAnimUp:    []int{16, 17, 18, 19},
		mAnimRight: []int{20, 21, 22, 23},
		mAnimDown:  []int{24, 25, 26, 27},
		mAnimLeft:  []int{28, 29, 30, 31},
		mFacing:    CharacterFacingDirection[2],
		mEntity: CreateEntity(CharacterDefinition{
			texture: pic, width: 16, height: 24,
			startFrame: 24,
			tileX:      4,
			tileY:      4,
		}),
		mController: StateMachineCreate(
			map[string]func() State{
				"wait": func() State {
					return WaitStateCreate(gHero, CastleRoomMap)
				},
				"move": func() State {
					return MoveStateCreate(gHero, CastleRoomMap)
				},
			},
		),
	}
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
			gHero.mController.Update(dt)
		}

		global.gWin.Update()
	}
}
