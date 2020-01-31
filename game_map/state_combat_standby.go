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
		Name:        "CSStandBy",
		Character:   char,
		CombatState: cs,
		Entity:      char.Entity,
	}
}

func (s *CSStandBy) Enter(data interface{}) {
}

func (s *CSStandBy) Render(win *pixelgl.Window) {
}

func (s *CSStandBy) Exit() {
}

func (s *CSStandBy) Update(dt float64) {
}
