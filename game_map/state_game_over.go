package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
)

type GameOverState struct {
	Stack    *gui.StateStack
	captions gui.SimpleCaptionsScreen
}

func GameOverStateCreate(stack *gui.StateStack, captions []gui.CaptionStyle) *GameOverState {
	return &GameOverState{
		Stack:    stack,
		captions: gui.SimpleCaptionsScreenCreate(captions, pixel.V(0, 0)),
	}
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *GameOverState) Enter() {
}

func (s GameOverState) Render(win *pixelgl.Window) {
	s.captions.Render(win)
}

func (s GameOverState) Exit() {
}

func (s *GameOverState) Update(dt float64) bool {
	return true
}

func (s GameOverState) HandleInput(win *pixelgl.Window) {
}
