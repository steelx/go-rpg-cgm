package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/gui"
)

type TweenEvent struct {
	Tween     animation.Tween
	Target    gui.StackInterface
	ApplyFunc func(e *TweenEvent)
}

func TweenEventCreate(start, finish, duration float64, target gui.StackInterface, applyFunc func(e *TweenEvent)) *TweenEvent {
	return &TweenEvent{
		Tween:     animation.TweenCreate(start, finish, duration),
		Target:    target,
		ApplyFunc: applyFunc,
	}
}

func (e *TweenEvent) Update(dt float64) {
	e.Tween.Update(dt)
	e.ApplyFunc(e)
}
func (e TweenEvent) IsBlocking() bool {
	return true
}
func (e TweenEvent) IsFinished() bool {
	return e.Tween.IsFinished()
}
func (e TweenEvent) Render(win *pixelgl.Window) {

}
