package globals

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/tilepix"
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
	PrimaryMonitor    *pixelgl.Monitor
	WindowHeight      float64
	WindowWidth       float64
	Vsync             bool
	Undecorated       bool
	ClearColor        pixel.RGBA
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
	ClearColor:        pixel.RGBA{0.2, 0.2, 0.2, 1.0},
	Win:               &pixelgl.Window{},
	CollisionLayer:    "collision",
	CollisionLayerPos: 3,
}

func init() {
	fontFace14, err := utilz.LoadTTF("../resources/font/joystix.ttf", 14)
	utilz.PanicIfErr(err)
	fontFace12, err := utilz.LoadTTF("../resources/font/joystix.ttf", 12)
	utilz.PanicIfErr(err)
	BasicAtlas14 = text.NewAtlas(fontFace14, text.ASCII)
	BasicAtlas12 = text.NewAtlas(fontFace12, text.ASCII)

	//Game Map TMX
	CastleMapDef, err = tilepix.ReadFile("small_room.tmx")
	utilz.PanicIfErr(err)

	//images for Textbox & Panel
	AvatarPng, err = utilz.LoadPicture("../resources/avatar.png")
	utilz.PanicIfErr(err)
	ContinueCaretPng, err = utilz.LoadPicture("../resources/continue_caret.png")
	utilz.PanicIfErr(err)
	CursorPng, err = utilz.LoadPicture("../resources/cursor.png")
	utilz.PanicIfErr(err)
	PanelPng, err = utilz.LoadPicture("../resources/simple_panel.png")
	utilz.PanicIfErr(err)
	ProgressBarBgPng, err = utilz.LoadPicture("../resources/progressbar_bg.png")
	utilz.PanicIfErr(err)
	ProgressBarFbPng, err = utilz.LoadPicture("../resources/progressbar_fg.png")
	utilz.PanicIfErr(err)
}
