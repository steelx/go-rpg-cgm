package globals

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

//=============================================================
// Global variables
//=============================================================
var (
	PanelPng         pixel.Picture
	CursorPng        pixel.Picture
	ContinueCaretPng pixel.Picture
	AvatarPng        pixel.Picture
	ProgressBarBgPng pixel.Picture
	ProgressBarFbPng pixel.Picture
	BasicAtlas12     *text.Atlas
	BasicAtlas14     *text.Atlas
	CastleMapDef     *tilepix.Map
)

type GlobalVars struct {
	WindowHeight      float64
	WindowWidth       float64
	Vsync             bool
	Undecorated       bool
	ClearColor        pixel.RGBA
	Win               *pixelgl.Window
	CollisionLayer    string
	CollisionLayerPos int
}

var Global = &GlobalVars{
	WindowHeight:      480,
	WindowWidth:       800,
	Vsync:             true,
	Undecorated:       false,
	ClearColor:        pixel.RGBA{0.2, 0.2, 0.2, 1.0},
	Win:               &pixelgl.Window{},
	CollisionLayer:    "collision",
	CollisionLayerPos: 3,
}

func init() {
	fontFace14, err := LoadTTF("../resources/font/joystix.ttf", 14)
	PanicIfErr(err)
	fontFace12, err := LoadTTF("../resources/font/joystix.ttf", 12)
	PanicIfErr(err)
	BasicAtlas14 = text.NewAtlas(fontFace14, text.ASCII)
	BasicAtlas12 = text.NewAtlas(fontFace12, text.ASCII)

	//Game Map TMX
	CastleMapDef, err = tilepix.ReadFile("small_room.tmx")
	PanicIfErr(err)

	//images for Textbox & Panel
	AvatarPng, err = LoadPicture("../resources/avatar.png")
	PanicIfErr(err)
	ContinueCaretPng, err = LoadPicture("../resources/continue_caret.png")
	PanicIfErr(err)
	CursorPng, err = LoadPicture("../resources/cursor.png")
	PanicIfErr(err)
	PanelPng, err = LoadPicture("../resources/simple_panel.png")
	PanicIfErr(err)
	ProgressBarBgPng, err = LoadPicture("../resources/progressbar_bg.png")
	PanicIfErr(err)
	ProgressBarFbPng, err = LoadPicture("../resources/progressbar_fg.png")
	PanicIfErr(err)
}
