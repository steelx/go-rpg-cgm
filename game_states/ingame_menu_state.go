package game_states

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type InGameMenuState struct {
	Stack        *gui.StateStack
	StateMachine *state_machine.StateMachine
}

func InGameMenuStateCreate(stack *gui.StateStack) InGameMenuState {
	igm := InGameMenuState{
		Stack: stack,
	}

	//igm.StateMachine = state_machine.Create(map[string]func() state_machine.State{
	//	"frontmenu": func() state_machine.State {
	//		//return FrontMenuStateCreate(this)
	//	},
	//	"items": func() state_machine.State {
	//		//return ItemMenuState:Create(this)
	//		return igm.StateMachine.Empty
	//	},
	//	"magic": func() state_machine.State {
	//		//return MagicMenuStateCreate(this)
	//		return igm.StateMachine.Empty
	//	},
	//	"equip": func() state_machine.State {
	//		//return EquipMenuStateCreate(this)
	//		return igm.StateMachine.Empty
	//	},
	//	"status": func() state_machine.State {
	//		//return StatusMenuStateCreate(this)
	//		return igm.StateMachine.Empty
	//	},
	//})

	igm.StateMachine.Change("frontmenu", nil)

	return igm
}

func (igm *InGameMenuState) Update(dt float64) bool {
	if reflect.DeepEqual(igm.Stack.Top(), igm) {
		igm.StateMachine.Update(dt)
	}
	return true
}
func (igm InGameMenuState) Render(win *pixelgl.Window) {
	igm.StateMachine.Render(win)
}

func (igm InGameMenuState) Enter()                          {}
func (igm InGameMenuState) Exit()                           {}
func (igm InGameMenuState) HandleInput(win *pixelgl.Window) {}
