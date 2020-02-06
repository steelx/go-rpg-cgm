package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type SimpleCaptionsScreen struct {
	captions []CaptionStyle
	Position pixel.Vec
}

type CaptionStyle struct {
	Text  string
	Scale float64
}

func SimpleCaptionsScreenCreate(captions []CaptionStyle, position pixel.Vec) SimpleCaptionsScreen {
	s := SimpleCaptionsScreen{
		captions: captions,
		Position: position,
	}

	return s
}

func (s SimpleCaptionsScreen) Render(win *pixelgl.Window) {
	var centerPosMain pixel.Vec
	var margin float64
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	for i := 0; i < len(s.captions); i++ {
		textBase := text.New(s.Position, basicAtlas)
		txt := s.captions[i].Text
		var centerPos pixel.Vec
		if i == 0 {
			centerPosMain = pixel.V(-textBase.BoundsOf(txt).W()/2, s.Position.Y)
			centerPos = pixel.V(-textBase.BoundsOf(txt).W()/2, s.Position.Y)
		}
		centerPos = pixel.V(centerPosMain.X, s.Position.Y-textBase.BoundsOf(txt).H()-margin)
		textBase.Clear()
		fmt.Fprintln(textBase, txt)
		textBase.Draw(win, pixel.IM.Moved(centerPos).Scaled(centerPos, s.captions[i].Scale))
		margin += 30
	}

	camera := pixel.IM.Scaled(centerPosMain, 1.0).Moved(win.Bounds().Center().Sub(centerPosMain))
	win.SetMatrix(camera)
}
