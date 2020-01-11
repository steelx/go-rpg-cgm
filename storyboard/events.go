package storyboard

import (
	"github.com/steelx/go-rpg-cgm/game_states"
	"github.com/steelx/go-rpg-cgm/gui"
	"image/color"
)

func Wait(seconds float64) *WaitEvent {
	return WaitEventCreate(seconds)
}

func BlackScreen(id string) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		screen := game_states.ScreenStateCreate(storyboard.Stack, color.RGBA{R: 255, G: 0, B: 0, A: 1})
		storyboard.PushState(id, screen)
		return WaitEventCreate(0)
	}
}

func FadeScreen(id string, start, finish, duration float64) func(storyboard *Storyboard, dt float64) TweenEvent {
	var dtTime float64
	return func(storyboard *Storyboard, dt float64) TweenEvent {
		dtTime += dt
		screen := gui.FadeScreenCreate(storyboard.Stack, uint8(start), uint8(finish), duration)
		storyboard.PushState(id, &screen)

		return TweenEventCreate(
			start, finish, duration,
			&screen,
			func(e *TweenEvent) {
				e.Tween.Update(dtTime)
				screen.Update(e.Tween.Value())
			},
		)
	}
}
