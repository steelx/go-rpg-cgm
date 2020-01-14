package character_states

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

type NPCWaitState struct {
	game_map.CharacterStateBase

	mFrameResetSpeed, mFrameCount float64
}

func NPCStandStateCreate(character *game_map.Character, gMap *game_map.GameMap) state_machine.State {
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
