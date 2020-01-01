package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//=============================================================
// Global variables
//=============================================================
type Global struct {
	gWindowHeight int
	gWindowWidth  int
	gVsync        bool
	gUndecorated  bool
	gClearColor   pixel.RGBA
	gWin          *pixelgl.Window
}

var global = &Global{
	gWindowHeight: 480,
	gWindowWidth:  800,
	gVsync:        true,
	gUndecorated:  false,
	gClearColor:   pixel.RGBA{0.2, 0.2, 0.2, 1.0},
	gWin:          &pixelgl.Window{},
}
