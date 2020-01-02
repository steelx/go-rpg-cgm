package main

import (
	"github.com/faiface/pixel/pixelgl"
)

type WaitState struct {
	mCharacter  *Character
	mMap        *GameMap
	mEntity     *Entity
	mController *StateMachine

	mFrameResetSpeed, mFrameCount float64
}

func WaitStateCreate(character *Character, gMap *GameMap) State {
	s := &WaitState{}
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

func (s *WaitState) Enter(data Direction) {
	// Reset to default frame
	s.mFrameCount = 0
	s.mEntity.SetFrame(s.mEntity.startFrame)

	//check if an EXIT Trigger exists on given tile coords
	tileX, tileY := s.mMap.GetTileIndex(s.mEntity.mTileX, s.mEntity.mTileY)
	if trigger := s.mMap.GetTrigger(tileX, tileY); trigger.OnExit != nil {
		trigger.OnExit()
	}
}

func (s *WaitState) Render() {
	//pixelgl renderer
	//s.mEntity.TeleportAndDraw(s.mMap)
}

func (s *WaitState) Exit() {}

func (s *WaitState) Update(dt float64) {
	// If we're in the wait state for a few frames, reset the frame to
	// the starting frame.
	if s.mFrameCount == 0 {
		s.mFrameCount = s.mFrameCount + dt
		if s.mFrameCount >= s.mFrameResetSpeed {
			s.mFrameCount = 0
			s.mEntity.SetFrame(s.mEntity.startFrame)
		}
	}

	if global.gWin.Pressed(pixelgl.KeyLeft) {
		s.mController.Change("move", Direction{-1, 0})
	}
	if global.gWin.Pressed(pixelgl.KeyRight) {
		s.mController.Change("move", Direction{1, 0})
	}
	if global.gWin.Pressed(pixelgl.KeyDown) {
		s.mController.Change("move", Direction{0, 1})
	}
	if global.gWin.Pressed(pixelgl.KeyUp) {
		s.mController.Change("move", Direction{0, -1})
	}
}
