package game_map

import "github.com/faiface/pixel/pixelgl"

type NonBlockingTimer struct {
	Seconds   float64
	ApplyFunc func(e *NonBlockingTimer)
	HasPopOut bool
}

func NonBlockingTimerCreate(seconds float64, applyFunc func(e *NonBlockingTimer)) *NonBlockingTimer {
	return &NonBlockingTimer{
		Seconds:   seconds,
		ApplyFunc: applyFunc,
	}
}

func (e NonBlockingTimer) TimeUp() bool {
	return e.Seconds <= 0
}

func (e *NonBlockingTimer) Update(dt float64) {
	e.Seconds = e.Seconds - dt
	e.ApplyFunc(e)
}
func (e NonBlockingTimer) IsBlocking() bool {
	return false
}
func (e NonBlockingTimer) IsFinished() bool {
	return e.TimeUp()
}
func (e NonBlockingTimer) Render(win *pixelgl.Window) {

}
