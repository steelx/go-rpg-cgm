package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/salviati/go-tmx/tmx"
	"time"
)

var (
	Map         = &GameMap{}
	CastleRoom1 = &GameMap{}
	gHero       Character
	camPos      = pixel.ZV
	//camSpeed    = 1000.0
	camZoom = 2.0
	//camZoomSpeed = 1.2
	frameRate = 10 * time.Millisecond
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "GP RPG",
		Bounds:      pixel.R(0, 0, float64(global.gWindowWidth*2), float64(global.gWindowHeight*2)),
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
	castleRoom1Tmx, err := tmx.ReadFile("./castle-room-1.tmx")
	panicIfErr(err)

	CastleRoom1.Create(castleRoom1Tmx)
	CastleRoom1.CamToTile(5, 6) //pan camera

	pic, err := LoadPicture("./resources/walk_cycle.png")
	panicIfErr(err)

	gHero = Character{
		mEntity: CreateEntity(CharacterDefinition{
			texture: pic, width: 16, height: 24,
			startFrame: 1,
			tileX:      9,
			tileY:      2,
		}),
		mController: StateMachineCreate(
			map[string]func() State{
				"wait": func() State {
					return WaitStateCreate(gHero, *CastleRoom1)
				},
				"move": func() State {
					return MoveStateCreate(gHero, *CastleRoom1)
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
	// Camera
	camPos = pixel.V(CastleRoom1.mCamX, CastleRoom1.mCamY)
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
	global.gWin.SetMatrix(cam)

	tick := time.Tick(frameRate)
	for !global.gWin.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}

		global.gWin.Clear(global.gClearColor)

		select {
		case <-tick:
			CastleRoom1.Render()
			gHero.mEntity.TeleportAndDraw(*CastleRoom1)
			gHero.mController.Update(dt)
		}

		global.gWin.Update()
	}
}
