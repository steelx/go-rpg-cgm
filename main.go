package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/salviati/go-tmx/tmx"
	"time"
)

var (
	Map         = &GameMap{}
	CastleRoom1 = &GameMap{}
	hero        pixel.Picture
	heroFrames  []pixel.Rect
	camPos      = pixel.ZV
	camSpeed    = 1000.0
	camZoom     = 2.0
	//camZoomSpeed = 1.2
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

	hero, err = LoadPicture("./resources/walk_cycle.png")
	panicIfErr(err)

	heroFrames = LoadAsFrames(hero, 16, 24)
}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	heroFrame := heroFrames[9]
	sprite := pixel.NewSprite(hero, heroFrame)

	for !global.gWin.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}

		global.gWin.Clear(global.gClearColor)

		camPos = pixel.V(CastleRoom1.mCamX, CastleRoom1.mCamY)
		// Camera movement
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
		global.gWin.SetMatrix(cam)
		if global.gWin.Pressed(pixelgl.KeyLeft) {
			CastleRoom1.mCamX -= camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyRight) {
			CastleRoom1.mCamX += camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyDown) {
			CastleRoom1.mCamY -= camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyUp) {
			CastleRoom1.mCamY += camSpeed * dt
		}
		if global.gWin.Pressed(pixelgl.KeyLeftControl) {
			fmt.Printf("x %v y %v \n", global.gWin.MousePosition().X, global.gWin.MousePosition().Y)
		}

		//camZoom *= math.Pow(camZoomSpeed, global.gWin.MouseScroll().Y)
		CastleRoom1.Render()
		tilePos := CastleRoom1.GetTilePositionAtFeet(9, 2, heroFrame.W(), heroFrame.H())
		sprite.Draw(global.gWin, pixel.IM.Moved(tilePos))

		global.gWin.Update()
	}
}
