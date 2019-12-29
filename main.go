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
	hero        = &Entity{}
	heroFrames  []pixel.Rect
	camPos      = pixel.ZV
	camSpeed    = 1000.0
	camZoom     = 2.0
	//camZoomSpeed = 1.2
	frameRate = 33 * time.Millisecond
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
	//tilemap, err := tmx.ReadFile("./larger_map.tmx")
	//panicIfErr(err)
	//Map.Create(tilemap)
	//Map.CamToTile(5, 50)

	castleRoom1Tmx, err := tmx.ReadFile("./castle-room-1.tmx")
	panicIfErr(err)

	CastleRoom1.Create(castleRoom1Tmx)
	CastleRoom1.CamToTile(5, 6) //pan camera

	heroDef := CharacterDefinition{
		texture:    "./resources/walk_cycle.png",
		width:      16,
		height:     24,
		startFrame: 8,
		tileX:      9,
		tileY:      2,
	}
	hero.Create(heroDef)
}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {

	// Camera
	camPos = pixel.V(CastleRoom1.mCamX, CastleRoom1.mCamY)
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
	global.gWin.SetMatrix(cam)

	tick := time.Tick(frameRate)
	for !global.gWin.Closed() {

		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}

		global.gWin.Clear(global.gClearColor)

		select {
		case <-tick:
			if global.gWin.JustPressed(pixelgl.KeyLeft) {
				hero.mTileX -= 1
			}
			if global.gWin.JustPressed(pixelgl.KeyRight) {
				hero.mTileX += 1
			}
			if global.gWin.JustPressed(pixelgl.KeyDown) {
				hero.mTileY += 1
			}
			if global.gWin.JustPressed(pixelgl.KeyUp) {
				hero.mTileY -= 1
			}

			CastleRoom1.Render()
			hero.TeleportAndDraw(*CastleRoom1)
		}

		global.gWin.Update()
	}
}
