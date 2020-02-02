package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type CSStandBy struct {
	Name        string
	Character   *Character
	CombatState *CombatState
	Entity      *Entity
	Anim        animation.Animation
}

//char *Character, cs *CombatState
func CSStandByCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	char := charV.Interface().(*Character)
	csV := reflect.ValueOf(args[1])
	cs := csV.Interface().(*CombatState)

	return &CSStandBy{
		Name:        csStandby,
		Character:   char,
		CombatState: cs,
		Entity:      char.Entity,
		Anim:        animation.Create([]int{char.Entity.StartFrame}, true, 0.16),
	}
}

func (s *CSStandBy) Enter(data ...interface{}) {
	animID := reflect.ValueOf(data[0])
	frames := s.Character.GetCombatAnim(animID.Interface().(string))
	s.Anim.SetFrames(frames)
}

func (s *CSStandBy) Render(win *pixelgl.Window) {
	//The *CombatState will do the render for us
}

func (s *CSStandBy) Exit() {
}

func (s *CSStandBy) Update(dt float64) {
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame())
}
