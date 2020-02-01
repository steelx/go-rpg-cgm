package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type CSRunAnim struct {
	Name        string
	Character   *Character
	CombatState *CombatState
	Entity      *Entity
	Anim        animation.Animation
	AnimId      string
}

//char *Character, cs *CombatState
func CSRunAnimCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	char := charV.Interface().(*Character)
	csV := reflect.ValueOf(args[1])
	cs := csV.Interface().(*CombatState)

	return &CSRunAnim{
		Name:        CS_RunAnim,
		Character:   char,
		CombatState: cs,
		Entity:      char.Entity,
	}
}

func (s *CSRunAnim) Enter(data ...interface{}) {
	animV := reflect.ValueOf(data[0])
	loop, spf := true, 0.12
	s.AnimId = animV.Interface().(string)

	frames := s.Character.GetCombatAnim(s.AnimId)
	s.Anim = animation.Create(frames, loop, spf)
}

func (s *CSRunAnim) Exit() {
}

func (s *CSRunAnim) Update(dt float64) {
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame())
}

func (s *CSRunAnim) Render(renderer *pixelgl.Window) {
}

func (s CSRunAnim) IsFinished() bool {
	return s.Anim.IsFinished()
}
