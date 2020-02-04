package game_map

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"math"
	"reflect"
)

type CSMoveParams struct {
	Dir, Distance, Time float64
}

type CSMove struct {
	Name                   string
	Character              *Character
	CombatState            *CombatState
	Entity                 *Entity
	Tween                  animation.Tween
	Anim                   animation.Animation
	AnimId                 string
	MoveTime, MoveDistance float64
	PixelX, PixelY         float64
}

//char *Character, cs *CombatState
func CSMoveCreate(args ...interface{}) state_machine.State {
	charV := reflect.ValueOf(args[0])
	char := charV.Interface().(*Character)
	csV := reflect.ValueOf(args[1])
	cs := csV.Interface().(*CombatState)

	return &CSMove{
		Name:         csMove,
		Character:    char,
		CombatState:  cs,
		Entity:       char.Entity,
		Anim:         animation.Create([]int{char.Entity.StartFrame}, true, 0.12),
		MoveTime:     0.3,
		MoveDistance: 32,
	}
}

func (s CSMove) IsFinished() bool {
	return s.Tween.IsFinished()
}

//data = CSMoveParams
func (s *CSMove) Enter(data ...interface{}) {
	if len(data) == 0 || !reflect.ValueOf(data[0]).CanInterface() {
		panic(fmt.Sprintf("Please pass CSMoveParams while changing State"))
		return
	}
	backForth := reflect.ValueOf(data[0]).Interface().(CSMoveParams)

	if len(data) == 3 {
		s.MoveTime = backForth.Time
		s.MoveDistance = backForth.Distance
	}

	s.AnimId = s.Name
	frames := s.Character.GetCombatAnim(s.Name)
	var dir float64 = -1
	if s.Character.Facing == CharacterFacingDirection[1] {
		s.AnimId = csRetreat
		frames = s.Character.GetCombatAnim(s.AnimId)
		dir = 1
	}
	dir = dir * backForth.Dir
	s.Anim.SetFrames(frames)

	// Store current position
	s.PixelX = s.Entity.X
	s.PixelY = s.Entity.Y

	s.Tween = animation.TweenCreate(0, dir, s.MoveTime)
}

func (s *CSMove) Exit() {
}

func (s *CSMove) Update(dt float64) {
	s.Anim.Update(dt)
	s.Entity.SetFrame(s.Anim.Frame())

	s.Tween.Update(dt)
	value := s.Tween.Value()
	x := s.PixelX + (value * s.MoveDistance)
	y := s.PixelY
	s.Entity.X = math.Floor(x)
	s.Entity.Y = math.Floor(y)
}

func (s *CSMove) Render(renderer *pixelgl.Window) {
}
