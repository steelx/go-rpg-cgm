package game_map

import (
	"github.com/faiface/pixel/pixelgl"
)

type SBEvent interface {
	Update(dt float64)
	IsBlocking() bool
	IsFinished() bool
	Render(win *pixelgl.Window)
}

type WaitEvent struct {
	Seconds float64
}

func WaitEventCreate(seconds float64) *WaitEvent {
	return &WaitEvent{
		Seconds: seconds,
	}
}

func (e *WaitEvent) Update(dt float64) {
	e.Seconds = e.Seconds - dt
}
func (e WaitEvent) IsBlocking() bool {
	return true
}
func (e WaitEvent) IsFinished() bool {
	return e.Seconds <= 0
}
func (e WaitEvent) Render(win *pixelgl.Window) {

}
