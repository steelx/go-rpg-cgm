package state_machine

import "github.com/faiface/pixel/pixelgl"

/*
Controller :
	Create({
		"wait" : function() return WaitStateCreate(Entity, gMap),
		"move" : function() return MoveStateCreate(gHero, gMap),
	})
*/

//
// StateMachine - a state machine
//
// Usage:
//
// gStateMachine = StateMachine:Create
// {
// 		['MainMenu'] = function()
// 			return MainMenu:Create()
// 		end,
// 		['InnerGame'] = function()
// 			return InnerGame:Create()
// 		end,
// 		['GameOver'] = function()
// 			return GameOver:Create()
// 		end,
// }
// gStateMachine:Change("MainGame")
//
type State interface {
	Enter(data interface{})
	Render(win *pixelgl.Window)
	Exit()
	Update(dt float64)
}

//StateMachine Controller
type StateMachine struct {
	states  map[string]func() State
	current State
}

func Create(states map[string]func() State) *StateMachine {
	return &StateMachine{
		states:  states,
		current: nil,
	}
}

//Change state
// e.g. Controller.Change("move", {x = -1, y = 0})
func (m *StateMachine) Change(stateName string, enterParams interface{}) {
	if m.current != nil {
		m.current.Exit()
	}
	m.current = m.states[stateName]()
	m.current.Enter(enterParams)
}

func (m *StateMachine) Update(dt float64) {
	m.current.Update(dt)
}

func (m *StateMachine) Render(win *pixelgl.Window) {
	m.current.Render(win)
}

func (m *StateMachine) Enter(data interface{}) {
}
func (m *StateMachine) Exit() {
}
