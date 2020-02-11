package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/utilz"
	"image/color"
	"math"
	"reflect"
)

type CombatTextFX struct {
	X, Y                float64
	CurrentY, VelocityY float64
	Text                string
	Color               color.RGBA
	Gravity             float64
	Alpha               float64
	Scale               float64
	HoldTime            float64
	HoldCounter         float64
	FadeSpeed           float64
	priority            int
}

func CombatTextFXCreate(x, y float64, txt string, params ...interface{}) *CombatTextFX {
	//hexColorI must be HEX string
	color_ := utilz.HexToColor("#FFFFFF")
	if len(params) > 0 {
		hexColor := reflect.ValueOf(params[0]).Interface().(string)
		color_ = utilz.HexToColor(hexColor)
	}
	priority := 2
	if len(params) >= 2 {
		priority = reflect.ValueOf(params[1]).Interface().(int)
	}

	return &CombatTextFX{
		X:           x,
		Y:           y,
		CurrentY:    y,
		VelocityY:   125,
		Text:        txt,
		Color:       color_,
		Gravity:     700,
		Alpha:       1,
		Scale:       1.5,
		HoldTime:    0.5,
		HoldCounter: 0,
		FadeSpeed:   3,
		priority:    priority,
	}
}

func (f *CombatTextFX) IsFinished() bool {
	return f.Alpha == 0
}

func (f *CombatTextFX) Update(dt float64) {
	f.CurrentY = f.CurrentY + (f.VelocityY * dt)
	f.VelocityY = f.VelocityY - (f.Gravity * dt)

	if f.CurrentY <= f.Y {
		f.CurrentY = f.Y
		f.HoldCounter = f.HoldCounter + dt
		if f.HoldCounter > f.HoldTime {
			f.Alpha = math.Max(0, f.Alpha-(dt*f.FadeSpeed))
			f.Color.A = utilz.GetAlpha(f.Alpha)
		}
	}
}

func (f *CombatTextFX) Render(renderer pixel.Target) {
	x := f.X
	y := f.CurrentY

	pos := pixel.V(x, y)
	textBase := text.New(pos, gui.BasicAtlasAscii)
	textBase.Color = f.Color
	fmt.Fprintln(textBase, f.Text)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, f.Scale))
}

func (f *CombatTextFX) Priority() int {
	return f.priority
}
