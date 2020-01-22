package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type CaptionScreen struct {
	text     string
	textBase *text.Text
	scale    float64
	Position pixel.Vec
}

func CaptionScreenCreate(txt string, position pixel.Vec, scale float64) CaptionScreen {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	return CaptionScreen{
		text:     txt,
		scale:    scale,
		Position: position,
		textBase: text.New(position, basicAtlas),
	}
}

/*
	StackInterface implemented below
*/
func (s CaptionScreen) Enter() {
}

func (s CaptionScreen) Exit() {
}

func (s CaptionScreen) Update(dt float64) bool {
	return true
}

func (s CaptionScreen) Render(win *pixelgl.Window) {
	centerPos := pixel.V(s.textBase.BoundsOf(s.text).W()/2, s.Position.Y)
	s.textBase.Clear()
	fmt.Fprintln(s.textBase, s.text)
	s.textBase.Draw(win, pixel.IM.Scaled(centerPos, s.scale))

	camera := pixel.IM.Scaled(centerPos, 1.0).Moved(win.Bounds().Center().Sub(centerPos))
	win.SetMatrix(camera)
}

func (s CaptionScreen) HandleInput(win *pixelgl.Window) {
}
