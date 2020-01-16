package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/utilz"
)

type WaitState struct {
	Character  *Character
	Map        *GameMap
	Entity     *Entity
	Controller *state_machine.StateMachine

	mFrameResetSpeed, FrameCount float64
}

func WaitStateCreate(character *Character, gMap *GameMap) state_machine.State {
	s := &WaitState{}
	s.Character = character
	s.Map = gMap
	s.Entity = character.Entity
	s.Controller = character.Controller

	s.mFrameResetSpeed = 0.015
	s.FrameCount = 0
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *WaitState) Enter(data interface{}) {
	// Reset to default frame
	s.FrameCount = 0
	s.Entity.SetFrame(s.Entity.StartFrame)

	//check if an EXIT Trigger exists on given tile coords
	//tileX, tileY := s.Map.GetTileIndex(s.Entity.TileX, s.Entity.TileY)
	tileX, tileY := s.Entity.TileX, s.Entity.TileY
	if trigger := s.Map.GetTrigger(tileX, tileY); trigger.OnExit != nil {
		trigger.OnExit(s.Map, s.Entity, tileX, tileY)
	}
}

func (s *WaitState) Render(win *pixelgl.Window) {
	//pixelgl renderer
	//s.Entity.TeleportAndDraw(s.Map)
}

func (s *WaitState) Exit() {}

func (s *WaitState) Update(dt float64) {
	// If we're in the wait state for a few frames, reset the frame to
	// the starting frame.
	if s.FrameCount == 0 {
		s.FrameCount = s.FrameCount + dt
		if s.FrameCount >= s.mFrameResetSpeed {
			s.FrameCount = 0
			s.Entity.SetFrame(s.Entity.StartFrame)
		}
	}

	if globals.Global.Win.Pressed(pixelgl.KeyLeft) {
		s.Controller.Change("move", utilz.Direction{-1, 0})
	}
	if globals.Global.Win.Pressed(pixelgl.KeyRight) {
		s.Controller.Change("move", utilz.Direction{1, 0})
	}
	if globals.Global.Win.Pressed(pixelgl.KeyDown) {
		s.Controller.Change("move", utilz.Direction{0, 1})
	}
	if globals.Global.Win.Pressed(pixelgl.KeyUp) {
		s.Controller.Change("move", utilz.Direction{0, -1})
	}
}
