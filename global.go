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
	gWindowHeight: 224,
	gWindowWidth:  256,
	gVsync:        true,
	gUndecorated:  false,
	gClearColor:   pixel.RGBA{0.4, 0.4, 0.4, 1.0},
	gWin:          &pixelgl.Window{},
}
