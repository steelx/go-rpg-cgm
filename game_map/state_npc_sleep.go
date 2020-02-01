package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

type SleepState struct {
	Character           *Character
	Map                 *GameMap
	Entity, SleepEntity *Entity
	Controller          *state_machine.StateMachine
	Anim                animation.Animation
}

//character *Character, gMap *GameMap
func SleepStateCreate(args ...interface{}) *SleepState {
	charV := reflect.ValueOf(args[0])
	character := charV.Interface().(*Character)
	gMapV := reflect.ValueOf(args[1])
	gMap := gMapV.Interface().(*GameMap)

	s := &SleepState{
		Character:   character,
		Map:         gMap,
		Entity:      character.Entity,
		Controller:  character.Controller,
		Anim:        animation.Create([]int{12, 13, 14, 15}, true, 0.3),
		SleepEntity: CreateEntity(Entities["sleeper"]),
	}

	s.Entity.SetFrame(character.Anims[character.Facing][0]) //13
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *SleepState) Enter(data ...interface{}) {
	s.Entity.AddChild("snore", s.SleepEntity)
}

func (s *SleepState) Render(win *pixelgl.Window) {}

func (s *SleepState) Exit() {
	s.Entity.RemoveChild("snore")
}

func (s *SleepState) Update(dt float64) {
	s.Anim.Update(dt)
	s.SleepEntity.SetFrame(s.Anim.Frame())
}
