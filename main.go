package main

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
	"sort"
	"time"
)

const camZoom = 1.0

var (
	fontFace      font.Face
	basicAtlas    *text.Atlas
	CastleRoomMap = &GameMap{}
	camPos        = pixel.ZV
	//camSpeed    = 1000.0
	//camZoomSpeed = 1.2
	frameRate = 15 * time.Millisecond
	panel     = imdraw.New(nil)
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
func init() {
	var err error
	fontFace, err = loadTTF("./resources/font/joystix.ttf", 14)
	panicIfErr(err)
	basicAtlas = text.NewAtlas(fontFace, text.ASCII)
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
			fmt.Println("Pot is full of snakes!")
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

			// Camera
			CastleRoomMap.CamToTile(gHero.mEntity.mTileX, gHero.mEntity.mTileY)
			camPos = pixel.V(CastleRoomMap.mCamX, CastleRoomMap.mCamY)
			cam := pixel.IM.Scaled(camPos, camZoom).Moved(global.gWin.Bounds().Center().Sub(camPos))
			global.gWin.SetMatrix(cam)

			DrawFixedTopPanel(CastleRoomMap)
		}

		global.gWin.Update()
	}
}
