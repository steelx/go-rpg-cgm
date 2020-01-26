package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"golang.org/x/image/font/basicfont"
	"reflect"
)

type StatusMenuState struct {
	parent                     *InGameMenuState
	win                        *pixelgl.Window
	Layout                     gui.Layout
	Stack                      *gui.StateStack
	StateMachine               *state_machine.StateMachine
	TopBarText, PrevTopBarText string
	EquipMenu, Actions         *gui.SelectionMenu
	Panels                     []gui.Panel
	ActorSummary               gui.ActorSummary
}

func StatusMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *StatusMenuState {
	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "title", "bottom", 0.12, 2)

	return &StatusMenuState{
		win:          win,
		parent:       parent,
		Stack:        parent.Stack,
		StateMachine: parent.StateMachine,
		Layout:       layout,
		Panels: []gui.Panel{
			layout.CreatePanel("title"),
			layout.CreatePanel("bottom"),
		},
	}
}

/////////////////////////////
// StateMachine impl below //
func (s *StatusMenuState) Enter(actorSumI interface{}) {
	actorSumV := reflect.ValueOf(actorSumI)
	s.ActorSummary = actorSumV.Interface().(gui.ActorSummary)

	equipmentMenu := gui.SelectionMenuCreate(26, 40,
		s.ActorSummary.Actor.ActiveEquipSlots,
		false,
		pixel.V(0, 0),
		func(i int, equipId interface{}) {
			fmt.Println(i, equipId)
		},
		func(a ...interface{}) {
			rendererV := reflect.ValueOf(a[0])
			renderer := rendererV.Interface().(pixel.Target)
			xV := reflect.ValueOf(a[1])
			x := xV.Interface().(float64)
			yV := reflect.ValueOf(a[2])
			y := yV.Interface().(float64)
			itemV := reflect.ValueOf(a[3])
			equipId := itemV.Interface().(int)
			s.ActorSummary.Actor.RenderEquipment(renderer, x, y, equipId)
		},
	)
	s.EquipMenu = &equipmentMenu
	s.EquipMenu.HideCursor()

	actionsMenu := gui.SelectionMenuCreate(18, 40,
		s.ActorSummary.Actor.Actions,
		false,
		pixel.V(0, 0),
		func(i int, equipId interface{}) {
			fmt.Println(i, equipId)
		},
		nil,
	)
	s.Actions = &actionsMenu
	s.Actions.HideCursor()
}

func (s StatusMenuState) Render(renderer *pixelgl.Window) {
	for _, v := range s.Panels {
		v.Draw(renderer)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	titleX := s.Layout.MidX("title")
	titleY := s.Layout.MidY("title")
	pos := pixel.V(titleX, titleY)
	textBase := text.New(pos, basicAtlas)
	fmt.Fprintln(textBase, "Status")
	textBase.Draw(renderer, pixel.IM)

	left := s.Layout.Left("bottom") + 10
	top := s.Layout.Top("bottom") - 10
	s.ActorSummary.SetPosition(left, top)
	s.ActorSummary.Render(renderer)

	xp := fmt.Sprintf("XP: %v/%v \n", s.ActorSummary.Actor.XP, s.ActorSummary.Actor.NextLevelXP)
	pos = pixel.V(left+240, top-58)
	textBase = text.New(pos, basicAtlas)
	fmt.Fprintln(textBase, xp)
	textBase.Draw(renderer, pixel.IM)

	s.EquipMenu.SetPosition(-10, -64)
	s.EquipMenu.Render(renderer)

	//stats
	stats := s.ActorSummary.Actor.Stats
	x := left + 106
	y := 0.0

	for k, v := range combat.ActorLabels.ActorStats {
		label := combat.ActorLabels.ActorStatLabels[k]
		s.DrawStat(renderer, x, y, label, stats.Get(v))
		y -= 16
	}

	y -= 16
	for k, v := range combat.ActorLabels.ItemStats {
		label := combat.ActorLabels.ItemStatLabels[k]
		s.DrawStat(renderer, x, y, label, stats.Get(v))
		y -= 16
	}

	// this should be a panel
	var x1, y1, w, h float64 = 75, 25, 100, 56
	box := gui.TextboxCreateFixed(
		s.Stack,
		"",
		pixel.V(x1, y1),
		w, h,
		"", nil,
		false,
	)
	box.AppearTween = animation.TweenCreate(1, 1, 0)
	box.Render(renderer)

	s.Actions.SetPosition(x-14, y-10)
	s.Actions.Render(renderer)
}

func (s StatusMenuState) Exit() {

}

func (s StatusMenuState) Update(dt float64) {

	if s.win.JustPressed(pixelgl.KeyEscape) {
		s.StateMachine.Change("frontmenu", nil)
	}
}

//////////////////////////////////////////////
// StatusMenuState additional methods below //
//////////////////////////////////////////////

func (s StatusMenuState) DrawStat(renderer pixel.Target, x, y float64, label string, value float64) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	pos := pixel.V(x, y)
	textBase := text.New(pos, basicAtlas)
	fmt.Fprintln(textBase, fmt.Sprintf("%-6s (%v)", label, value))
	textBase.Draw(renderer, pixel.IM)
}
