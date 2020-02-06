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

	return &CombatTextFX{
		X:           x,
		Y:           y,
		CurrentY:    y,
		VelocityY:   125,
		Text:        txt,
		Color:       color_,
		Gravity:     700,
		Alpha:       1,
		Scale:       1.2,
		HoldTime:    0.5,
		HoldCounter: 0,
		FadeSpeed:   3,
		priority:    2,
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
	y := math.Floor(f.CurrentY)

	pos := pixel.V(x, y)
	shadow := utilz.HexToColor("#000000")
	shadow.A = f.Color.A
	textBase := text.New(pos, gui.BasicAtlasAscii)
	textBase.Color = shadow
	fmt.Fprintln(textBase, f.Text)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, f.Scale).Moved(pos))

	pos = pixel.V(x+0.5, y-0.5)
	textBase = text.New(pos, gui.BasicAtlas12)
	textBase.Color = f.Color
	fmt.Fprintln(textBase, f.Text)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, f.Scale).Moved(pos))
}

func (f *CombatTextFX) Priority() int {
	return f.priority
}
