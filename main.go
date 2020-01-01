package main

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

var (
	CastleRoomMap = &GameMap{}
	gHero         Character
	camPos        = pixel.ZV
	//camSpeed    = 1000.0
	camZoom = 1.4
	//camZoomSpeed = 1.2
	//frameRate = 10 * time.Millisecond
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

	gHero = Character{
		mAnimUp:    []int{0, 9, 18, 27},
		mAnimRight: []int{36, 45, 54, 63},
		mAnimDown:  []int{72, 81, 90, 99},
		mAnimLeft:  []int{108, 117, 126, 135},
		mEntity: CreateEntity(CharacterDefinition{
			texture: pic, width: 16, height: 24,
			startFrame: 0,
			tileX:      7,
			tileY:      2,
		}),
		mController: StateMachineCreate(
			map[string]func() State{
				"wait": func() State {
					return WaitStateCreate(gHero, *CastleRoomMap)
				},
				"move": func() State {
					return MoveStateCreate(gHero, *CastleRoomMap)
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

	//tick := time.Tick(frameRate)
	for !global.gWin.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}

		global.gWin.Clear(global.gClearColor)

		CastleRoomMap.DrawAfter(1, func(canvas *pixelgl.Canvas) {
			gHero.mEntity.TeleportAndDraw(*CastleRoomMap, canvas)
		})
		gHero.mController.Update(dt)

		global.gWin.Update()
	}
}
