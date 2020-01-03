package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//=============================================================
// Global variables
//=============================================================
type Global struct {
	gWindowHeight     float64
	gWindowWidth      float64
	gVsync            bool
	gUndecorated      bool
	gClearColor       pixel.RGBA
	gWin              *pixelgl.Window
	collisionLayer    string
	collisionLayerPos int
}

var global = &Global{
	gWindowHeight:     480,
	gWindowWidth:      800,
	gVsync:            true,
	gUndecorated:      false,
	gClearColor:       pixel.RGBA{0.2, 0.2, 0.2, 1.0},
	gWin:              &pixelgl.Window{},
	collisionLayer:    "collision",
	collisionLayerPos: 3,
}
