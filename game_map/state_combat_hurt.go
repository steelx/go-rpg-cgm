package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type CSHurt struct {
	Name        string
	Character   *Character
	CombatState *CombatState
	Entity      *Entity
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
}

func (s *CSHurt) Exit() {
}

func (s *CSHurt) Update(dt float64) {
}

func (s *CSHurt) Render(renderer *pixelgl.Window) {
}
