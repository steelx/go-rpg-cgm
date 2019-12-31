package main

import "github.com/faiface/pixel/pixelgl"

type WaitState struct {
	mCharacter  FSMObject
	mMap        GameMap
	mEntity     *Entity
	mController *StateMachine
}

func WaitStateCreate(character FSMObject, gMap GameMap) State {
	s := &WaitState{}
	s.mCharacter = character
	s.mMap = gMap
	s.mEntity = character.mEntity
	s.mController = character.mController
	return s
}

//The StateMachine class requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *WaitState) Enter(data Direction) {
	// Reset to default frame
	s.mEntity.SetFrame(s.mEntity.startFrame)
}

func (s *WaitState) Render() {
	//pixelgl renderer
}

func (s *WaitState) Exit() {}

func (s *WaitState) Update(dt float64) {
	if global.gWin.JustPressed(pixelgl.KeyLeft) {
		s.mController.Change("move", Direction{-1, 0})
	}
	if global.gWin.JustPressed(pixelgl.KeyRight) {
		s.mController.Change("move", Direction{1, 0})
	}
	if global.gWin.JustPressed(pixelgl.KeyDown) {
		s.mController.Change("move", Direction{0, 1})
	}
	if global.gWin.JustPressed(pixelgl.KeyUp) {
		s.mController.Change("move", Direction{0, -1})
	}
}
