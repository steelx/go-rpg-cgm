package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type CSHurt struct {
	Name        string
	Character   *Character
	CombatState *CombatState
	Entity      *Entity
	Anim        animation.Animation
	AnimId      string
	PrevState   state_machine.State
}

//char *Character, cs *CombatState
func CSHurtCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	char := charV.Interface().(*Character)
	csV := reflect.ValueOf(args[1])
	cs := csV.Interface().(*CombatState)

	return &CSHurt{
		Name:        csHurt,
		Character:   char,
		CombatState: cs,
		Entity:      char.Entity,
	}
}

func (s CSHurt) IsFinished() bool {
	return true
}

func (s *CSHurt) Enter(data ...interface{}) {
	if len(data) == 0 {
		panic("Please pass 'state_machine.State' with Character.Controller.Change(csHurt, '_Char.Controller.Current_')")
	}
	state := reflect.ValueOf(data[0]).Interface().(state_machine.State)
	s.PrevState = state
	s.AnimId = s.Name
	frames := s.Character.GetCombatAnim(s.AnimId)
	s.Anim = animation.Create(frames, false, 0.4)
}

func (s *CSHurt) Exit() {
}

func (s *CSHurt) Update(dt float64) {
	if s.Anim.IsFinished() {
		s.Character.Controller.Current = s.PrevState
	}
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame())
}

func (s *CSHurt) Render(renderer *pixelgl.Window) {
}
