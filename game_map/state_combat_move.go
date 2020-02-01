package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type CSMove struct {
	Name        string
	Character   *Character
	CombatState *CombatState
	Entity      *Entity
	Anim        animation.Animation
}

//char *Character, cs *CombatState
func CSMoveCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	char := charV.Interface().(*Character)
	csV := reflect.ValueOf(args[1])
	cs := csV.Interface().(*CombatState)

	return &CSMove{
		Name:        CS_Move,
		Character:   char,
		CombatState: cs,
		Entity:      char.Entity,
	}
}

func (s *CSMove) Enter(data ...interface{}) {
}

func (s *CSMove) Exit() {
}

func (s *CSMove) Update(dt float64) {
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame())
}

func (s *CSMove) Render(renderer *pixelgl.Window) {
}
