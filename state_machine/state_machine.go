package state_machine

import (
	"github.com/faiface/pixel/pixelgl"
)

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
	Enter(data ...interface{})
	Render(win *pixelgl.Window)
	Exit()
	Update(dt float64)
	IsFinished() bool
}

//StateMachine Controller
type StateMachine struct {
	states  map[string]func() State
	Current State
}

func Create(states map[string]func() State) *StateMachine {
	return &StateMachine{
		states:  states,
		Current: nil,
	}
}

//Change state
// e.g. Controller.Change("move", {x = -1, y = 0})
func (s *StateMachine) Change(stateName string, enterParams ...interface{}) {
	if s.Current != nil {
		s.Current.Exit()
	}
	s.Current = s.states[stateName]()
	s.Current.Enter(enterParams...)
}

func (s StateMachine) IsFinished() bool {
	return true
}

func (s *StateMachine) Update(dt float64) {
	s.Current.Update(dt)
}

func (s *StateMachine) Render(win *pixelgl.Window) {
	s.Current.Render(win)
}

func (s *StateMachine) Enter(data ...interface{}) {
}
func (s *StateMachine) Exit() {
}
