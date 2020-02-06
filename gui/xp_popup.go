package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/utilz"
	"image/color"
	"math"
)

type XPPopUp struct {
	X, Y      float64
	Text      string
	TextColor color.RGBA
	Tween     animation.Tween
	FadeTime,
	DisplayTime float64
	Pane     Panel
	textBase *text.Text
}

func XPPopUpCreate(x, y float64, txt, color string) *XPPopUp {
	pos := pixel.V(x, y)
	textBase := text.New(pos, BasicAtlas12)
	txtSize := textBase.BoundsOf(txt)

	return &XPPopUp{
		X:           pos.X,
		Y:           pos.Y,
		Text:        txt,
		textBase:    textBase,
		TextColor:   utilz.HexToColor(color),
		FadeTime:    0.25,
		DisplayTime: 0,
		Pane:        PanelCreate(pixel.V(pos.X+txtSize.W()/2, pos.Y+txtSize.H()/2), txtSize.W()+20, txtSize.H()+20),
	}
}

func (p *XPPopUp) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

func (p *XPPopUp) TurnOn() {
	p.Tween = animation.TweenCreate(0, 1, p.FadeTime)
}
func (p *XPPopUp) TurnOff() {
	current := p.Tween.Value()
	p.Tween = animation.TweenCreate(current, 0, current*p.FadeTime)
}

func (p XPPopUp) IsTurningOff() bool {
	return p.Tween.FinishValue() == 0
}

func (p XPPopUp) IsFinished() bool {
	return p.Tween.FinishValue() == 0 && p.Tween.Value() == 0
}

func (p *XPPopUp) Update(dt float64) {
	p.Tween.Update(dt)
	if p.Tween.IsFinished() {
		p.DisplayTime = math.Min(5, p.DisplayTime+dt)
	}
}

func (p *XPPopUp) Render(renderer pixel.Target) {
	alpha := p.Tween.Value()
	p.TextColor.A = utilz.GetAlpha(alpha)

	paneColor := utilz.HexToColor("#333333")
	paneColor.A = utilz.GetAlpha(alpha)
	p.Pane.BGColor = paneColor

	p.Pane.Draw(renderer)

	p.textBase.Clear()
	p.textBase.Color = p.TextColor
	fmt.Fprintln(p.textBase, p.Text)
	p.textBase.Draw(renderer, pixel.IM)
}
