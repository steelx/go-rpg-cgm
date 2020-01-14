package character_states

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/state_machine"
)

type SleepState struct {
	Character           *game_map.Character
	Map                 *game_map.GameMap
	Entity, SleepEntity *game_map.Entity
	Controller          *state_machine.StateMachine
	Anim                animation.Animation
}

func SleepStateCreate(character *game_map.Character, gMap *game_map.GameMap) state_machine.State {
	s := &SleepState{
		Character:   character,
		Map:         gMap,
		Entity:      character.Entity,
		Controller:  character.Controller,
		Anim:        animation.AnimationCreate([]int{12, 13, 14, 15}, true, 0.12),
		SleepEntity: game_map.CreateEntity(Entities["sleeper"]),
	}

	s.Entity.SetFrame(character.Anims[character.Facing][0]) //13
	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *SleepState) Enter(data interface{}) {
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
