package game_map

import "github.com/faiface/pixel"

/*
IsFinished :Reports when the effect is finished so it can be removed from the list.
Update: Updates the effect according to the elapsed frame time.
Render :Renders the effect to the screen.
Priority: Controls the render order. For instance, the jumping numbers should appear
on top of everything else. Lower priority numbers are rendered later. 0 is
considered the highest priority.
*/
type EffectState interface {
	IsFinished() bool
	Update(dt float64)
	Render(renderer pixel.Target)
	Priority() int //0 is highest
}
