package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/utilz"
	"image/color"
	"reflect"
)

type JumpingNumbersFX struct {
	X, Y, CurrentY, VelocityY,
	Number, // to display
	Gravity, //pixels per second
	FadeDistance, // pixels
	Scale float64
	Color    color.RGBA
	priority int
}

func JumpingNumbersFXCreate(x, y, number float64, hexColorI ...interface{}) *JumpingNumbersFX {
	//hexColorI must be HEX string
	color_ := utilz.HexToColor("#FFFFFF")
	if len(hexColorI) > 0 {
		hexColor := reflect.ValueOf(hexColorI).Interface().(string)
		color_ = utilz.HexToColor(hexColor)
	}

	return &JumpingNumbersFX{
		X:            x,
		Y:            y,
		CurrentY:     y,
		VelocityY:    230,
		Gravity:      400,
		FadeDistance: 33,
		Scale:        1.5,
		Number:       number,
		Color:        color_,
		priority:     0,
	}
}

func (f *JumpingNumbersFX) IsFinished() bool {
	return f.CurrentY <= (f.Y - f.FadeDistance)
}

func (f *JumpingNumbersFX) Update(dt float64) {
	f.CurrentY = f.CurrentY + (f.VelocityY * dt)
	f.VelocityY = f.VelocityY - (f.Gravity * dt)
	if f.CurrentY <= f.Y {
		fade01 := (f.Y - f.CurrentY) / f.FadeDistance
		f.Color.A = 255 - uint8(fade01)
	}
}

func (f *JumpingNumbersFX) Render(renderer pixel.Target) {
	pos := pixel.V(f.X, f.CurrentY)
	textBase := text.New(pos, gui.BasicAtlas12)
	textBase.Color = f.Color

	fmt.Fprintln(textBase, f.Number)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, f.Scale).Moved(pos))
}

func (f *JumpingNumbersFX) Priority() int {
	return f.priority
}
