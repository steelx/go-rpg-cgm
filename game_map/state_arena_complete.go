package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
)

type ArenaCompleteState struct {
	Stack     *gui.StateStack
	prevState gui.StackInterface
	captions  gui.SimpleCaptionsScreen
}

func ArenaCompleteStateCreate(stack *gui.StateStack, prevState gui.StackInterface) *ArenaCompleteState {
	captions := []gui.CaptionStyle{
		{"YOU WON!", 3},
		{"Champion of the Arena", 1},
	}
	return &ArenaCompleteState{
		Stack:     stack,
		prevState: prevState,
		captions:  gui.SimpleCaptionsScreenCreate(captions, pixel.V(0, 0)),
	}
}

func (s *ArenaCompleteState) Enter() {

}

func (s *ArenaCompleteState) Exit() {

}

func (s *ArenaCompleteState) Update(dt float64) bool {
	return false
}

func (s *ArenaCompleteState) Render(win *pixelgl.Window) {
	s.captions.Render(win)
}

func (s *ArenaCompleteState) HandleInput(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeySpace) || win.JustPressed(pixelgl.KeyEscape) {
		s.Stack.Pop()
		s.Stack.Push(s.prevState)
	}
}
