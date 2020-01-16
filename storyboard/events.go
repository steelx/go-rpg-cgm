package storyboard

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/actions"
	"github.com/steelx/go-rpg-cgm/game_map"
	"github.com/steelx/go-rpg-cgm/gui"
	"image/color"
	"reflect"
)

func Wait(seconds float64) *WaitEvent {
	return WaitEventCreate(seconds)
}

func BlackScreen(id string) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		screen := game_map.ScreenStateCreate(storyboard.Stack, color.RGBA{R: 255, G: 255, B: 255, A: 1})
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
		mapInfo := game_map.MapsDB[mapName](storyboard.Stack)
		exploreState := game_map.ExploreStateCreate(nil, mapInfo, win)
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
		char := game_map.Characters[entityDef](exploreState.Map)
		runFunc(char)
		return WaitEventCreate(seconds)
	}
}

func getExploreState(storyboard *Storyboard, stateId string) *game_map.ExploreState {
	exploreStateI := storyboard.States[stateId]
	exploreStateV := reflect.ValueOf(exploreStateI)
	exploreState := exploreStateV.Interface().(*game_map.ExploreState)
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

func Say(mapName, npcId, textMessage string, time float64) func(storyboard *Storyboard) *TimedTextboxEvent {
	return func(storyboard *Storyboard) *TimedTextboxEvent {
		exploreState := getExploreState(storyboard, mapName)
		npc := exploreState.Map.NPCbyId[npcId]
		tileX, tileY := npc.GetFacedTileCoords()
		posX, posY := exploreState.Map.GetTileIndex(tileX, tileY)
		tBox := storyboard.InternalStack.PushFitted(posX, posY+32, textMessage)
		return TimedTextboxEventCreate(tBox, time)
	}
}

//ReplaceScene will remove mapName and add newMapName with a Hero at given Tile X, Y
func ReplaceScene(mapName string, newMapName string, tileX, tileY float64, hideHero bool, win *pixelgl.Window) func(storyboard *Storyboard) *NonBlockEvent {
	return func(storyboard *Storyboard) *NonBlockEvent {
		storyboard.RemoveState(mapName) //remove previous map (exploreState)

		mapInfo := game_map.MapsDB[newMapName](storyboard.Stack)
		newExploreState := game_map.ExploreStateCreate(nil, mapInfo, win)

		if hideHero {
			newExploreState.HideHero()
		} else {
			newExploreState.ShowHero(tileX, tileY)
		}

		storyboard.PushState(newMapName, &newExploreState) //ADD new map (exploreState)

		return NonBlockEventCreate(0)
	}
}

//HandOffToMainStack will remove the exploreState from Storyboard and push it to main stack
func HandOffToMainStack(mapName string) func(storyboard *Storyboard) *WaitEvent {
	return func(storyboard *Storyboard) *WaitEvent {
		exploreState := getExploreState(storyboard, mapName)
		storyboard.Stack.Pop()
		exploreState.Stack = storyboard.Stack

		storyboard.Stack.Push(exploreState)

		return WaitEventCreate(1)
	}
}
