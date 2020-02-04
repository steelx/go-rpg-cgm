package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"image/color"
	"reflect"
)

//csEnemyDie
type CSEnemyDie struct {
	Name        string
	Character   *Character
	CombatState *CombatState
	Entity      *Entity
	Tween       animation.Tween
}

//char *Character, cs *CombatState
func CSEnemyDieCreate(args ...interface{}) *CSEnemyDie {
	charV := reflect.ValueOf(args[0])
	char := charV.Interface().(*Character)
	csV := reflect.ValueOf(args[1])
	cs := csV.Interface().(*CombatState)

	s := &CSEnemyDie{
		Name:        csEnemyDie,
		Character:   char,
		CombatState: cs,
		Entity:      char.Entity,
	}

	return s
}

func (s *CSEnemyDie) Enter(data ...interface{}) {
	s.Tween = animation.TweenCreate(1, 0, 1)
}

func (s *CSEnemyDie) Render(win *pixelgl.Window) {
	alpha := s.Tween.Value()
	color_ := color.RGBA{255, 255, 255, 255 - uint8(alpha)}
	s.Entity.Sprite.DrawColorMask(win, pixel.IM, color_)
}

func (s *CSEnemyDie) Exit() {

}

func (s *CSEnemyDie) Update(dt float64) {
	s.Tween.Update(dt)
}

func (s *CSEnemyDie) IsFinished() bool {
	return s.Tween.IsFinished()
}
