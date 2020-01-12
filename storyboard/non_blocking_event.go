package storyboard

import "github.com/faiface/pixel/pixelgl"

type NonBlockEvent struct {
	Seconds float64
}

func NonBlockEventCreate(seconds float64) *NonBlockEvent {
	return &NonBlockEvent{
		Seconds: seconds,
	}
}

func (e *NonBlockEvent) Update(dt float64) {
	e.Seconds = e.Seconds - dt
}
func (e NonBlockEvent) IsBlocking() bool {
	return false
}
func (e NonBlockEvent) IsFinished() bool {
	//return e.Seconds <= 0
	return false
}
func (e NonBlockEvent) Render(win *pixelgl.Window) {

}
