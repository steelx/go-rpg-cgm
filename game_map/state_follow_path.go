package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

type FollowPathState struct {
	Character  *Character
	Map        GameMap
	Entity     Entity
	Controller *state_machine.StateMachine
}

func FollowPathStateCreate(character *Character, gMap *GameMap) state_machine.State {
	s := FollowPathState{}
	s.Character = character
	s.Map = *gMap
	s.Entity = *character.Entity
	s.Controller = character.Controller

	return &s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *FollowPathState) Enter(data interface{}) {

	if s.Character.PathIndex >= len(s.Character.Path) || len(s.Character.Path) == 0 {
		s.Character.DefaultState = s.Character.PrevDefaultState //we set at Character.FollowPath
		s.Controller.Change(s.Character.DefaultState, Direction{0, 0})
		return
	}

	direction := s.Character.Path[s.Character.PathIndex]
	if direction == "left" {
		s.Controller.Change("move", Direction{-1, 0})
	} else if direction == "right" {
		s.Controller.Change("move", Direction{1, 0})
	} else if direction == "up" {
		s.Controller.Change("move", Direction{0, -1})
	} else if direction == "down" {
		s.Controller.Change("move", Direction{0, 1})
	}
}

func (s *FollowPathState) Render(win *pixelgl.Window) {}

func (s *FollowPathState) Exit() {
	s.Character.PathIndex = s.Character.PathIndex + 1
}

func (s *FollowPathState) Update(dt float64) {}
