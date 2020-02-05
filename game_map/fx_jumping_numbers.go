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

func JumpingNumbersFXCreate(x, y, number float64, params ...interface{}) *JumpingNumbersFX {
	//hexColorI must be HEX string
	color_ := utilz.HexToColor("#FFFFFF")
	priority := 0
	scale := 1.5
	if len(params) > 0 {
		hexColor := reflect.ValueOf(params[0]).Interface().(string)
		color_ = utilz.HexToColor(hexColor)
	}
	if len(params) >= 2 {
		priority = reflect.ValueOf(params[1]).Interface().(int)
	}
	if len(params) >= 3 {
		scale = reflect.ValueOf(params[2]).Interface().(float64)
	}

	return &JumpingNumbersFX{
		X:            x,
		Y:            y,
		CurrentY:     y,
		VelocityY:    230,
		Gravity:      400,
		FadeDistance: 33,
		Scale:        scale,
		Number:       number,
		Color:        color_,
		priority:     priority,
	}
}

func (f *JumpingNumbersFX) IsFinished() bool {
	return f.CurrentY <= (f.Y - f.FadeDistance)
}

func (f *JumpingNumbersFX) Update(dt float64) {
	f.CurrentY = f.CurrentY + (f.VelocityY * dt)
	f.VelocityY = f.VelocityY - (f.Gravity * dt)
	if f.CurrentY <= f.Y {
		alpha := (f.Y - f.CurrentY) / f.FadeDistance
		f.Color.A = utilz.GetAlpha(alpha)
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
