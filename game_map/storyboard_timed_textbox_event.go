package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
)

type TimedTextboxEvent struct {
	Textbox   *gui.Textbox
	Countdown float64
}

func TimedTextboxEventCreate(tbox *gui.Textbox, countDown float64) *TimedTextboxEvent {
	return &TimedTextboxEvent{
		Textbox:   tbox,
		Countdown: countDown,
	}
}

func (e *TimedTextboxEvent) Update(dt float64) {
	e.Countdown = e.Countdown - dt
	if e.Countdown <= 0 {
		e.Textbox.OnClick()
	}
}
func (e TimedTextboxEvent) IsBlocking() bool {
	return e.Countdown > 0
}

func (e TimedTextboxEvent) IsFinished() bool {
	return !e.IsBlocking()
}
func (e TimedTextboxEvent) Render(win *pixelgl.Window) {
}
