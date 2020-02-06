package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"reflect"
)

type GameOverState struct {
	Stack    *gui.StateStack
	captions gui.SimpleCaptionsScreen
	World    *combat.WorldExtended
	Menu     *gui.SelectionMenu
}

func GameOverStateCreate(stack *gui.StateStack) *GameOverState {
	world := reflect.ValueOf(stack.Globals["world"]).Interface().(*combat.WorldExtended)
	captions := []gui.CaptionStyle{
		{"Game Over", 3},
	}
	s := &GameOverState{
		Stack:    stack,
		World:    world,
		captions: gui.SimpleCaptionsScreenCreate(captions, pixel.V(0, 0)),
	}

	menu := gui.SelectionMenuCreate(36, 0, 0,
		[]string{"Continue", "New Game"},
		false,
		pixel.V(0, 0),
		s.OnSelection,
		nil,
	)
	menu.SetPosition(-menu.GetWidth()/2, -50)
	s.Menu = &menu

	return s
}

//The StateMachine requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *GameOverState) Enter() {
}

func (s GameOverState) Render(renderer *pixelgl.Window) {
	s.captions.Render(renderer)
	s.Menu.Render(renderer)
}

func (s GameOverState) Exit() {
}

func (s *GameOverState) Update(dt float64) bool {
	return true
}

func (s *GameOverState) HandleInput(win *pixelgl.Window) {
	s.Menu.HandleInput(win)
}

func (s *GameOverState) OnSelection(index int, str interface{}) {
	const (
		continueGame int = iota
		newgame
	)

	if index == continueGame {
		logrus.Info("No Save system yet.")
		return
	}

	if index == newgame {
		s.Stack.Clear()
		newWorld := combat.WorldExtendedCreate()
		newWorld.Party.Add(combat.ActorCreate(combat.HeroDef))
		s.Stack.Globals["world"] = newWorld
		storyboard := StoryboardCreate(s.Stack, s.Stack.Win, IntroScene, false)
		s.Stack.Push(storyboard)
		return
	}
}
