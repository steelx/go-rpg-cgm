package character_states

import (
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

type NPCStrollWaitState struct {
	game_map.CharacterStateBase

	mFrameResetSpeed, mFrameCount float64
	mCountDown                    float64
}

func NPCStrollWaitStateCreate(character *game_map.Character, gMap *game_map.GameMap) state_machine.State {
	s := &NPCStrollWaitState{}
	s.Character = character
	s.Map = gMap
	s.Entity = character.Entity
	s.Controller = character.Controller

	s.mFrameResetSpeed = 0.015
	s.mFrameCount = 0
	s.mCountDown = globals.RandFloat(0, 3)
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *NPCStrollWaitState) Enter(data globals.Direction) {
	s.mFrameCount = 0
	s.mCountDown = globals.RandFloat(0, 3)
}

func (s *NPCStrollWaitState) Render() {}

func (s *NPCStrollWaitState) Exit() {}

func (s *NPCStrollWaitState) Update(dt float64) {
	// If we're in the wait state for a few frames, reset the frame to
	// the starting frame.
	if s.mFrameCount == 0 {
		s.mFrameCount = s.mFrameCount + dt
		if s.mFrameCount >= s.mFrameResetSpeed {
			s.mFrameCount = 0
			s.Entity.SetFrame(s.Entity.StartFrame)
		}
	}

	s.mCountDown = s.mCountDown - dt
	if s.mCountDown <= 0 {
		choice := globals.RandInt(0, 4)
		if choice == 1 {
			s.Controller.Change("move", globals.Direction{-1, 0})
		}
		if choice == 2 || choice == 0 {
			s.Controller.Change("move", globals.Direction{1, 0})
		}
		if choice == 3 {
			s.Controller.Change("move", globals.Direction{0, 1})
		}
		if choice == 4 {
			s.Controller.Change("move", globals.Direction{0, -1})
		}
	}
}
