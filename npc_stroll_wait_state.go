package main

type NPCStrollWaitState struct {
	CharacterStateBase

	mFrameResetSpeed, mFrameCount float64
	mCountDown                    float64
}

func NPCStrollWaitStateCreate(character *Character, gMap *GameMap) State {
	s := &NPCStrollWaitState{}
	s.mCharacter = character
	s.mMap = gMap
	s.mEntity = character.mEntity
	s.mController = character.mController

	s.mFrameResetSpeed = 0.015
	s.mFrameCount = 0
	s.mCountDown = randFloat(0, 3)
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *NPCStrollWaitState) Enter(data Direction) {
	s.mFrameCount = 0
	s.mCountDown = randFloat(0, 3)
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
			s.mEntity.SetFrame(s.mEntity.startFrame)
		}
	}

	s.mCountDown = s.mCountDown - dt
	if s.mCountDown <= 0 {
		choice := randInt(0, 4)
		if choice == 1 {
			s.mController.Change("move", Direction{-1, 0})
		}
		if choice == 2 || choice == 0 {
			s.mController.Change("move", Direction{1, 0})
		}
		if choice == 3 {
			s.mController.Change("move", Direction{0, 1})
		}
		if choice == 4 {
			s.mController.Change("move", Direction{0, -1})
		}
	}
}
