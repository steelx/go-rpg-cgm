package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"reflect"
)

type ArenaRound struct {
	Name   string
	Locked bool
}

type ArenaState struct {
	prevState gui.StackInterface
	Stack     *gui.StateStack
	World     *combat.WorldExtended
	Layout    gui.Layout
	Panels    []gui.Panel
	Selection *gui.SelectionMenu

	Rounds []ArenaRound
}

func ArenaStateCreate(stack *gui.StateStack, prevState gui.StackInterface) gui.StackInterface {
	layout := gui.LayoutCreate(0, 0, stack.Win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.15, 0)
	layout.SplitVert("bottom", "", "bottom", 0.75, 0)
	layout.SplitVert("bottom", "bottom", "", 0.5, 0)
	layout.Contract("bottom", -20, 40)
	layout.SplitHorz("bottom", "header", "bottom", 0.18, 2)

	gWorld := reflect.ValueOf(stack.Globals["world"]).Interface().(*combat.WorldExtended)
	s := &ArenaState{
		prevState: prevState,
		Stack:     stack,
		World:     gWorld,
		Layout:    layout,
		Rounds: []ArenaRound{
			{Name: "Round 1", Locked: false},
			{Name: "Round 2", Locked: true},
			{Name: "Round 3", Locked: true},
			{Name: "Round 4", Locked: true},
			{Name: "Round 5", Locked: true},
		},
	}

	s.Panels = []gui.Panel{
		layout.CreatePanel("top"),
		layout.CreatePanel("bottom"),
		layout.CreatePanel("header"),
	}

	roundsSelectionMenu := gui.SelectionMenuCreate(25, 25, 70,
		s.Rounds,
		false,
		pixel.V(0, 0),
		func(index int, round interface{}) {

		},
		s.RenderRoundItem,
	)
	txtSize := 100.0
	xPos := -roundsSelectionMenu.GetWidth() / 2
	xPos += roundsSelectionMenu.CursorWidth / 2
	xPos -= txtSize / 2
	roundsSelectionMenu.SetPosition(xPos, 18)
	s.Selection = &roundsSelectionMenu

	return s
}

//renderer pixel.Target, x, y float64, item ArenaRound
func (s *ArenaState) RenderRoundItem(a ...interface{}) {
	renderer := reflect.ValueOf(a[0]).Interface().(pixel.Target)
	x := reflect.ValueOf(a[1]).Interface().(float64)
	y := reflect.ValueOf(a[2]).Interface().(float64)
	round := reflect.ValueOf(a[3]).Interface().(ArenaRound)

	lockLabel := "Open"
	if round.Locked {
		lockLabel = "Locked"
	}

	label := fmt.Sprintf("%s: %s", round.Name, lockLabel)
	textBase := text.New(pixel.V(x, y), gui.BasicAtlas12)
	fmt.Fprintf(textBase, label)
	textBase.Draw(renderer, pixel.IM)
}

func (s *ArenaState) Enter() {

}

func (s *ArenaState) Exit() {

}

func (s *ArenaState) HandleInput(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEscape) {
		s.Stack.Pop() //remove self
		s.Stack.Push(s.prevState)
		return
	}

	s.Selection.HandleInput(win)
}

func (s *ArenaState) Update(dt float64) bool {

	return false
}

func (s *ArenaState) Render(renderer *pixelgl.Window) {
	//for _, v := range s.Panels {
	//	v.Draw(renderer)
	//}

	titleX := s.Layout.MidX("top")
	titleY := s.Layout.MidY("top")
	pos := pixel.V(titleX, titleY)
	textBase := text.New(pos, gui.BasicAtlasAscii)
	titleTxt := "Welcome to the Arena"
	pos = pixel.V(titleX-textBase.BoundsOf(titleTxt).W()/2, titleY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintf(textBase, titleTxt)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, 2))

	headerX := s.Layout.MidX("header")
	headerY := s.Layout.MidY("header")
	pos = pixel.V(headerX, headerY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintf(textBase, "Choose Round")
	textBase.Draw(renderer, pixel.IM)

	s.Selection.Render(renderer)

	//camera
	camera := pixel.IM.Scaled(pixel.ZV, 1.0).Moved(renderer.Bounds().Center().Sub(pixel.ZV))
	renderer.SetMatrix(camera)
}
