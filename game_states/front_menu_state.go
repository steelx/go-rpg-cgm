package game_states

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"golang.org/x/image/font/basicfont"
)

type FrontMenuState struct {
	Parent       *InGameMenuState
	Layout       gui.Layout
	Stack        *gui.StateStack
	StateMachine *state_machine.StateMachine
	TopBarText   string
	Selections   *gui.SelectionMenu
	Panels       []gui.Panel
	win          *pixelgl.Window
}

func FrontMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) FrontMenuState {

	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 0, 0)
	layout.SplitHorz("screen", "top", "bottom", 0.12, 2)
	layout.SplitVert("bottom", "left", "party", 0.726, 2)
	layout.SplitHorz("left", "menu", "gold", 0.7, 2)

	fm := FrontMenuState{
		win:          win,
		Parent:       parent,
		Stack:        parent.Stack,
		StateMachine: parent.StateMachine,
		Layout:       layout,
		TopBarText:   "Current Map Name",
	}

	selectionsX, selectionsY := fm.Layout.MidX("menu")-60, fm.Layout.Top("menu")-24
	selectionMenu := gui.SelectionMenuCreate(
		[]string{"Items", "Magic", "Equipment", "Status", "Save"},
		false,
		pixel.V(selectionsX, selectionsY),
		func(i int, str string) {
			fmt.Println("Menu", i, str)
			fm.OnMenuClick(i, str)
		},
	)
	fm.Selections = &selectionMenu
	fm.Panels = []gui.Panel{
		layout.CreatePanel("gold"),
		layout.CreatePanel("top"),
		layout.CreatePanel("party"),
		layout.CreatePanel("menu"),
	}

	return fm
}
func (fm *FrontMenuState) OnMenuClick(index int, str string) {
	ITEMS := 0
	if index == ITEMS {
		fm.StateMachine.Change("items", nil)
		return
	}
}

/*
   StateMachine :: State impl below
*/
func (fm FrontMenuState) Enter(data interface{}) {
}

func (fm FrontMenuState) Exit() {
}

func (fm FrontMenuState) Update(dt float64) {
	fm.Selections.HandleInput(fm.win)

	if fm.win.JustPressed(pixelgl.KeyBackspace) || fm.win.JustPressed(pixelgl.KeyEscape) {
		fm.Stack.Pop()
	}
}

func (fm FrontMenuState) Render(renderer *pixelgl.Window) {
	for _, p := range fm.Panels {
		p.Draw(renderer)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//renderer:ScaleText(1.5, 1.5)
	//renderer:AlignText("left", "center")
	menuX := fm.Layout.Left("menu") - 16
	menuY := fm.Layout.Top("menu") - 24
	fm.Selections.SetPosition(menuX, menuY)
	fm.Selections.Render(renderer)

	nameX := fm.Layout.MidX("top")
	nameY := fm.Layout.MidY("top")
	topBarText := text.New(pixel.V(nameX, nameY), basicAtlas)
	topBarText = text.New(pixel.V(nameX-topBarText.BoundsOf(fm.TopBarText).W()/2, nameY), basicAtlas)
	fmt.Fprintln(topBarText, fm.TopBarText)
	topBarText.Draw(renderer, pixel.IM)

	goldX := fm.Layout.MidX("gold") - 22
	goldY := fm.Layout.MidY("gold") + 22

	//renderer:ScaleText(1.22, 1.22)
	//renderer:AlignText("right", "top")
	topBarText = text.New(pixel.V(goldX, goldY), basicAtlas)
	fmt.Fprintln(topBarText, "GP :")
	topBarText.Draw(renderer, pixel.IM)

	topBarText = text.New(pixel.V(goldX, goldY+25), basicAtlas)
	fmt.Fprintln(topBarText, "TIME :")
	topBarText.Draw(renderer, pixel.IM)

	//renderer:AlignText("left", "top")
	topBarText = text.New(pixel.V(goldX+10, goldY), basicAtlas)
	fmt.Fprintln(topBarText, "0")
	topBarText.Draw(renderer, pixel.IM)

	topBarText = text.New(pixel.V(goldX+10, goldY+25), basicAtlas)
	fmt.Fprintln(topBarText, "0")
	topBarText.Draw(renderer, pixel.IM)
}
