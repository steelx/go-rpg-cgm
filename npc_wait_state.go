package main

type NPCWaitState struct {
	mCharacter  *Character
	mMap        *GameMap
	mEntity     *Entity
	mController *StateMachine

	mFrameResetSpeed, mFrameCount float64
}

func NPCWaitStateCreate(character *Character, gMap *GameMap) State {
	s := &NPCWaitState{}
	s.mCharacter = character
	s.mMap = gMap
	s.mEntity = character.mEntity
	s.mController = character.mController

	s.mFrameResetSpeed = 0.015
	s.mFrameCount = 0
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *NPCWaitState) Enter(data Direction) {}

func (s *NPCWaitState) Render() {}

func (s *NPCWaitState) Exit() {}

func (s *NPCWaitState) Update(dt float64) {}
