package gui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"image/color"
)

//FadeState
type FadeScreen struct {
	Stack                   *StateStack
	AlphaStart, AlphaFinish uint8
	Duration                float64
	Color                   color.RGBA
	Tween                   animation.Tween
	imd                     *imdraw.IMDraw
	position                pixel.Vec
}

func FadeScreenCreate(stack *StateStack, alphaStart, alphaFinish uint8, duration float64) FadeScreen {

	fd := FadeScreen{
		Stack:       stack,
		AlphaStart:  (alphaStart >> 24) & 0xff,
		AlphaFinish: (alphaFinish >> 24) & 0xff,
		Duration:    duration,
		imd:         imdraw.New(nil),
		position:    pixel.V(0, 0),
	}
	fd.Color = color.RGBA{R: 255, G: 255, B: 255, A: fd.AlphaStart}
	fd.Tween = animation.TweenCreate(float64(fd.AlphaStart), float64(fd.AlphaFinish), fd.Duration)
	return fd
}

func (fd FadeScreen) Enter() {
}
func (fd FadeScreen) Exit() {
}
func (fd *FadeScreen) Update(dt float64) bool {
	fd.Tween.Update(dt)
	alpha := fd.Tween.Value()
	fd.Color = color.RGBA{R: 255, G: 255, B: 255, A: uint8(alpha)}
	if fd.Tween.IsFinished() {
		fd.Stack.Pop()
	}
	return true
}
func (fd FadeScreen) HandleInput(win *pixelgl.Window) {
}
func (fd FadeScreen) Render(win *pixelgl.Window) {
	fd.imd.Clear()
	toTheScreen := pixel.V(fd.position.X, fd.position.Y)
	// Draw the rectangle.
	fd.imd.Color = fd.Color
	fd.imd.Push(win.Bounds().Min.Sub(pixel.V(win.Bounds().W()/2, win.Bounds().H()/2)), win.Bounds().Max)
	fd.imd.Rectangle(0)
	fd.imd.SetMatrix(pixel.IM.Moved(toTheScreen))

	camera := pixel.IM.Scaled(toTheScreen, 1.0).Moved(win.Bounds().Center().Sub(toTheScreen))
	win.SetMatrix(camera)
	fd.imd.Draw(win)
}
