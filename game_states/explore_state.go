package game_states

import "github.com/steelx/go-rpg-cgm/globals"

type ExploreState struct {
}

func (es ExploreState) Enter(data globals.Direction) {}
func (es ExploreState) Exit()                        {}
func (es ExploreState) Update(dt float64)            {}
func (es ExploreState) Render()                      {}
func (es ExploreState) HandleInput()                 {}
