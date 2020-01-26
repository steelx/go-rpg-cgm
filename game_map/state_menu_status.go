package game_map

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type StatusMenuState struct {
	parent                     *InGameMenuState
	win                        *pixelgl.Window
	Layout                     gui.Layout
	Stack                      *gui.StateStack
	StateMachine               *state_machine.StateMachine
	TopBarText, PrevTopBarText string
	Selections                 *gui.SelectionMenu
	Panels                     []gui.Panel
}

func StatusMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *StatusMenuState {
	return &StatusMenuState{
		win:          win,
		parent:       parent,
		Stack:        parent.Stack,
		StateMachine: parent.StateMachine,
	}
}

/////////////////////////////
// StateMachine impl below //
func (s StatusMenuState) Enter(actorI interface{}) {
	actorV := reflect.ValueOf(actorI)
	actor := actorV.Interface().(combat.Actor)
	fmt.Println(actor.Name)
}

func (s StatusMenuState) Render(win *pixelgl.Window) {

}

func (s StatusMenuState) Exit() {

}

func (s StatusMenuState) Update(dt float64) {
	//s.Selections.HandleInput(s.win)
	if s.win.JustPressed(pixelgl.KeyEscape) {
		s.Stack.Pop()
	}
}

//////////////////////////////////////////////
// StatusMenuState additional methods below //
//////////////////////////////////////////////
