package game_map

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/animation"
	"reflect"
)

type AnimEntityFx struct {
	X, Y     float64
	Entity   *Entity
	Anim     animation.Animation
	priority int
}

func AnimEntityFxCreate(x, y float64, entityDef EntityDefinition, frames []int, args ...interface{}) *AnimEntityFx {
	spf := 0.09
	if len(args) >= 1 {
		spf = reflect.ValueOf(args[0]).Interface().(float64)
	}

	return &AnimEntityFx{
		X:        x,
		Y:        y,
		Entity:   CreateEntity(entityDef),
		Anim:     animation.Create(frames, false, spf),
		priority: 1,
	}
}

func (f AnimEntityFx) IsFinished() bool {
	return f.Anim.IsFinished()
}

func (f *AnimEntityFx) Update(dt float64) {
	f.Anim.Update(dt)
	f.Entity.SetFrame(f.Anim.Frame())
}

func (f *AnimEntityFx) Render(renderer pixel.Target) {
	f.Entity.Render(nil, renderer, pixel.V(f.X, f.Y))
}

func (f AnimEntityFx) Priority() int {
	return f.priority
}
