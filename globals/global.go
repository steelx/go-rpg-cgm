package globals

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/utilz"
	"image/color"
)

//=============================================================
// Global variables
//=============================================================

type GlobalVars struct {
	PrimaryMonitor    *pixelgl.Monitor
	WindowHeight      float64
	WindowWidth       float64
	Vsync             bool
	Undecorated       bool
	ClearColor        color.Color
	Win               *pixelgl.Window
	DeltaTime         float64
	CollisionLayer    string
	CollisionLayerPos int
}

var Global = &GlobalVars{
	WindowHeight:      480,
	WindowWidth:       800,
	Vsync:             true,
	Undecorated:       false,
	ClearColor:        utilz.HexToColor("#12161A"),
	Win:               &pixelgl.Window{},
	CollisionLayer:    "collision",
	CollisionLayerPos: 3,
}
