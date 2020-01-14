package storyboard

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/actions"
	"github.com/steelx/go-rpg-cgm/game_map/character_states"
	"github.com/steelx/go-rpg-cgm/game_states"
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
		screen := game_states.ScreenStateCreate(storyboard.Stack, color.RGBA{R: 255, G: 255, B: 255, A: 1})
		storyboard.PushState(id, screen)
		return WaitEventCreate(0)
	}
}

//FadeScreen not working properly
func FadeScreen(id string, start, finish, duration float64) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		screen := gui.FadeScreenCreate(storyboard.Stack, uint8(start), uint8(finish), duration)
		storyboard.PushState(id, &screen)

		return WaitEventCreate(0)
	}
}

func TitleCaptionScreen(id string, txt string, duration float64) func(storyboard *Storyboard) *TweenEvent {
	return func(storyboard *Storyboard) *TweenEvent {
		captions := gui.CaptionScreenCreate(txt, pixel.V(0, 100), 3)
		storyboard.PushState(id, &captions)

		return TweenEventCreate(
			1, 0, duration,
			&captions,
			func(e *TweenEvent) {
				captions.Update(e.Tween.Value())
			},
		)
	}
}

func SubTitleCaptionScreen(id string, txt string, duration float64) func(storyboard *Storyboard) *TweenEvent {

	return func(storyboard *Storyboard) *TweenEvent {
		captions := gui.CaptionScreenCreate(txt, pixel.V(0, 50), 1)
		storyboard.PushState(id, &captions)

		return TweenEventCreate(
			1, 0, duration,
			&captions,
			func(e *TweenEvent) {
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

		return NonBlockEventCreate(0)
	}
}

//player_house, def = "sleeper", x = 14, y = 19
func RunActionAddNPC(mapName, entityDef string, x, y, seconds float64) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		exploreState := getExploreState(storyboard, mapName)
		exploreState.Hero.Entity.SetTilePos(x, y)
		runFunc := actions.ActionAddNPC(exploreState.Map, x, y)
		char := character_states.Characters[entityDef](exploreState.Map)
		runFunc(char)
		return WaitEventCreate(seconds)
	}
}

func getExploreState(storyboard *Storyboard, stateId string) *game_states.ExploreState {
	exploreStateI := storyboard.States[stateId]
	exploreStateV := reflect.ValueOf(exploreStateI)
	exploreState := exploreStateV.Interface().(*game_states.ExploreState)
	return exploreState
}

func KillState(id string) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		storyboard.RemoveState(id)
		return WaitEventCreate(0)
	}
}

func MoveNPC(npcId, mapName string, path []string) func(storyboard *Storyboard) *BlockUntilEvent {

	return func(storyboard *Storyboard) *BlockUntilEvent {
		exploreState := getExploreState(storyboard, mapName)
		npc := exploreState.Map.NPCbyId[npcId]
		npc.FollowPath(path)

		return BlockUntilEventCreate(func() bool {
			return npc.PathIndex > len(path)
		})
	}
}
