package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type NPCWaitState struct {
	CharacterStateBase

	mFrameResetSpeed, mFrameCount float64
}

func NPCStandCombatStateCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	character := charV.Interface().(*Character)

	s := &NPCWaitState{}
	s.Character = character
	s.Entity = character.Entity
	s.Controller = character.Controller

	s.mFrameResetSpeed = 0.015
	s.mFrameCount = 0
	return s
}

func NPCStandStateCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	character := charV.Interface().(*Character)
	gMapV := reflect.ValueOf(args[1])
	gMap := gMapV.Interface().(*GameMap)

	s := &NPCWaitState{}
	s.Character = character
	s.Map = gMap
	s.Entity = character.Entity
	s.Controller = character.Controller

	s.mFrameResetSpeed = 0.015
	s.mFrameCount = 0
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *NPCWaitState) Enter(data interface{}) {}

func (s *NPCWaitState) Render(win *pixelgl.Window) {}

func (s *NPCWaitState) Exit() {}

func (s *NPCWaitState) Update(dt float64) {}
