package storyboard

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/actions"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/game_states"
	"github.com/steelx/go-rpg-cgm/globals"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/maps_db"
	"image/color"
	"reflect"
)

func Wait(seconds float64) *WaitEvent {
	return WaitEventCreate(seconds)
}

func BlackScreen(id string) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		screen := game_states.ScreenStateCreate(storyboard.Stack, color.RGBA{R: 255, G: 0, B: 0, A: 1})
		storyboard.PushState(id, screen)
		return WaitEventCreate(0)
	}
}

//pending KillState and FadeOutState 323
//FadeScreen not working properly
func FadeScreen(id string, start, finish, duration float64) func(storyboard *Storyboard, dt float64) TweenEvent {
	var dtTime float64
	return func(storyboard *Storyboard, dt float64) TweenEvent {
		dtTime += dt
		screen := gui.FadeScreenCreate(storyboard.Stack, uint8(start), uint8(finish), duration)
		storyboard.PushState(id, &screen)

		return TweenEventCreate(
			start, finish, duration,
			&screen,
			func(e *TweenEvent) {
				e.Tween.Update(dtTime)
				screen.Update(e.Tween.Value())
			},
		)
	}
}

func TitleCaptionScreen(id string, txt string, duration float64) func(storyboard *Storyboard, dt float64) TweenEvent {
	var dtTime float64
	return func(storyboard *Storyboard, dt float64) TweenEvent {
		dtTime += dt
		captions := gui.CaptionScreenCreate(txt, pixel.V(0, 100), 3)
		storyboard.PushState(id, &captions)

		return TweenEventCreate(
			1, 0, duration,
			&captions,
			func(e *TweenEvent) {
				e.Tween.Update(dtTime)
				captions.Update(e.Tween.Value())
			},
		)
	}
}

func SubTitleCaptionScreen(id string, txt string, duration float64) func(storyboard *Storyboard, dt float64) TweenEvent {
	var dtTime float64
	return func(storyboard *Storyboard, dt float64) TweenEvent {
		dtTime += dt
		captions := gui.CaptionScreenCreate(txt, pixel.V(0, 50), 1)
		storyboard.PushState(id, &captions)

		return TweenEventCreate(
			1, 0, duration,
			&captions,
			func(e *TweenEvent) {
				e.Tween.Update(dtTime)
				captions.Update(e.Tween.Value())
			},
		)
	}
}

func Scene(mapName string, hideHero bool, win *pixelgl.Window) func(storyboard *Storyboard) *NonBlockEvent {

	return func(storyboard *Storyboard) *NonBlockEvent {
		gMap, collision, collisionLayerName := maps_db.MapsDB[mapName]()
		exploreState := game_states.ExploreStateCreate(nil,
			gMap, collision, collisionLayerName, win,
		)
		if hideHero {
			exploreState.HideHero()
		}

		storyboard.PushState(mapName, &exploreState)

		return NonBlockEventCreate(0.1)
	}
}

//player_house, def = "sleeper", x = 14, y = 19
func RunActionAddNPC(mapName, entityDef string, x, y float64) func(storyboard *Storyboard) {
	return func(storyboard *Storyboard) {
		gMap := GetMapRef(storyboard, mapName)
		runFunc := actions.ActionAddNPC(gMap, globals.Direction{x, y})
		char := character_states.Characters[entityDef](gMap)
		runFunc(char)
	}
}

func GetMapRef(storyboard *Storyboard, stateId string) *game_map.GameMap {
	exploreStateI := storyboard.States[stateId]
	exploreStateV := reflect.ValueOf(exploreStateI)
	exploreState := exploreStateV.Interface().(*game_states.ExploreState)
	return exploreState.Map
}
