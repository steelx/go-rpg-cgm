package game_states

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
	"image/color"
)

type ScreenState struct {
	Stack *gui.StateStack
	Color color.Color
}

func ScreenStateCreate(stack *gui.StateStack, color color.Color) ScreenState {
	//Color: color.RGBA{R: 0, G: 0, B: 0, A: 1}, //Black
	//Color: color.RGBA{R: 255, G: 255, B: 255, A: 1}, //white
	//Color: color.RGBA{R: 255, G: 255, B: 255, A: 0}, //white opacity 0
	s := ScreenState{
		Stack: stack,
		Color: color,
	}

	return s
}

/*
	StackInterface implemented below
*/

func (s ScreenState) Enter() {
}

func (s ScreenState) Exit() {
}

func (s ScreenState) Update(dt float64) bool {
	return true
}

func (s ScreenState) Render(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Clear()
	toTheScreen := pixel.V(0, 0)
	// Draw the rectangle.
	imd.Color = s.Color
	imd.Push(win.Bounds().Min.Sub(pixel.V(win.Bounds().W()/2, win.Bounds().H()/2)), win.Bounds().Max)
	imd.Rectangle(0)
	imd.SetMatrix(pixel.IM.Moved(toTheScreen))

	camera := pixel.IM.Scaled(toTheScreen, 1.0).Moved(win.Bounds().Center().Sub(toTheScreen))
	win.SetMatrix(camera)
	imd.Draw(win)
}

func (s ScreenState) HandleInput(win *pixelgl.Window) {
}
