package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/salviati/go-tmx/tmx"
	"math"
	"time"
)

const mapPath = "./larger_map.tmx"

var (
	Map          = &GameMap{}
	camPos       = pixel.ZV
	camSpeed     = 1000.0
	camZoom      = 2.0
	camZoomSpeed = 1.2
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
	// Initialize art assets (i.e. the tilemap)
	tilemap, err := tmx.ReadFile(mapPath)
	panicIfErr(err)
	Map.Create(tilemap)
	Map.GotoTile(5, 50)
}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()

	for !global.gWin.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}

		global.gWin.Clear(global.gClearColor)

		camPos = pixel.V(Map.mCamX, Map.mCamY)
		// Camera movement
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
		global.gWin.SetMatrix(cam)
		if global.gWin.Pressed(pixelgl.KeyLeft) {
			Map.mCamX -= camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyRight) {
			Map.mCamX += camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyDown) {
			Map.mCamY -= camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyUp) {
			Map.mCamY += camSpeed * dt
		}

		camZoom *= math.Pow(camZoomSpeed, global.gWin.MouseScroll().Y)

		Map.Render()
		global.gWin.Update()
	}
}
