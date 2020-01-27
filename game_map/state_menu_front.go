package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"golang.org/x/image/font/basicfont"
	"reflect"
)

type FrontMenuState struct {
	Parent                     *InGameMenuState
	Layout                     gui.Layout
	Stack                      *gui.StateStack
	StateMachine               *state_machine.StateMachine
	TopBarText, PrevTopBarText string
	Selections                 *gui.SelectionMenu
	PartyMenu                  *gui.SelectionMenu
	Panels                     []gui.Panel
	win                        *pixelgl.Window
	InPartyMenu                bool
}

func FrontMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *FrontMenuState {

	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 0, 0)
	layout.SplitHorz("screen", "top", "bottom", 0.12, 2)
	layout.SplitVert("bottom", "left", "party", 0.726, 2)
	layout.SplitHorz("left", "menu", "gold", 0.7, 2)

	fm := &FrontMenuState{
		win:          win,
		Parent:       parent,
		Stack:        parent.Stack,
		StateMachine: parent.StateMachine,
		Layout:       layout,
		TopBarText:   "Game Paused",
	}
	fm.PrevTopBarText = fm.TopBarText

	selectionsX, selectionsY := fm.Layout.MidX("menu")-60, fm.Layout.Top("menu")-24
	selectionMenu := gui.SelectionMenuCreate(32, 128,
		[]string{"Status", "Items"},
		false,
		pixel.V(selectionsX, selectionsY),
		fm.OnMenuClick,
		nil,
	)
	fm.Selections = &selectionMenu
	fm.Panels = []gui.Panel{
		layout.CreatePanel("gold"),
		layout.CreatePanel("top"),
		layout.CreatePanel("party"),
		layout.CreatePanel("menu"),
	}

	partyMembersMenu := gui.SelectionMenuCreate(100, 0,
		fm.CreatePartySummaries(),
		false,
		pixel.V(0, 0),
		fm.OnPartyMemberChosen,
		func(a ...interface{}) {
			//renderer pixel.Target, x, y float64, actorSummary ActorSummary
			rendererV := reflect.ValueOf(a[0])
			renderer := rendererV.Interface().(pixel.Target)
			xV := reflect.ValueOf(a[1])
			x := xV.Interface().(float64)
			yV := reflect.ValueOf(a[2])
			y := yV.Interface().(float64)
			actorSummaryV := reflect.ValueOf(a[3])
			actorSummary := actorSummaryV.Interface().(gui.ActorSummary)

			actorSummary.SetPosition(x, y+35)
			actorSummary.Render(renderer)
		},
	)
	partyMembersMenu.HideCursor()
	fm.PartyMenu = &partyMembersMenu

	return fm
}
func (fm *FrontMenuState) OnPartyMemberChosen(index int, actorSummaryI interface{}) {
	actorSummaryV := reflect.ValueOf(actorSummaryI)
	actorSummary := actorSummaryV.Interface().(gui.ActorSummary)

	fm.StateMachine.Change("status", actorSummary)
}

func (fm *FrontMenuState) OnMenuClick(index int, str interface{}) {
	if index == items {
		fm.StateMachine.Change("items", nil)
		return
	}

	if index == frontmenu {
		fm.InPartyMenu = true
		fm.Selections.HideCursor()
		fm.PartyMenu.ShowCursor()
		fm.PrevTopBarText = fm.TopBarText
		fm.TopBarText = "Choose a party member"
	}

}

func (fm FrontMenuState) CreatePartySummaries() []gui.ActorSummary {
	partyMembers := fm.Parent.World.Party.Members
	var summaryList []gui.ActorSummary
	for _, actor := range partyMembers {
		summaryList = append(summaryList, gui.ActorSummaryCreate(actor, true))
	}
	return summaryList
}

func (fm *FrontMenuState) goBackToFrontMenu() {
	fm.InPartyMenu = false
	fm.TopBarText = fm.PrevTopBarText
	fm.PartyMenu.HideCursor()
	fm.Selections.ShowCursor()
}

/*
   StateMachine :: State impl below
*/
func (fm FrontMenuState) Enter(data interface{}) {
}

func (fm FrontMenuState) Exit() {
}

func (fm *FrontMenuState) Update(dt float64) {

	if fm.InPartyMenu {
		fm.PartyMenu.HandleInput(fm.win)
		if fm.win.JustPressed(pixelgl.KeyEscape) {
			fm.goBackToFrontMenu()
		}
		return
	}

	fm.Selections.HandleInput(fm.win)
	if fm.win.JustPressed(pixelgl.KeyEscape) {
		fm.Stack.Pop()
	}
}

func (fm FrontMenuState) Render(renderer *pixelgl.Window) {
	for _, p := range fm.Panels {
		p.Draw(renderer)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//Selection List
	menuX := fm.Layout.Left("menu") - 16
	menuY := fm.Layout.Top("menu") - 24
	fm.Selections.SetPosition(menuX, menuY)
	fm.Selections.Render(renderer)

	//TOP Headline
	nameX := fm.Layout.MidX("top")
	nameY := fm.Layout.MidY("top")
	textBase := text.New(pixel.V(nameX, nameY), basicAtlas)
	textBase = text.New(pixel.V(nameX-getTextW(textBase, fm.TopBarText)/2, nameY), basicAtlas)
	fmt.Fprintln(textBase, fm.TopBarText)
	textBase.Draw(renderer, pixel.IM)

	//Bottom Left
	goldX := fm.Layout.Left("gold") + 16
	goldY := fm.Layout.Top("gold") - 24
	textBase = text.New(pixel.V(goldX, goldY), basicAtlas)
	fmt.Fprintln(textBase, "GP :")
	textBase.Draw(renderer, pixel.IM)

	textBase = text.New(pixel.V(goldX, goldY-25), basicAtlas)
	fmt.Fprintln(textBase, "TIME :")
	textBase.Draw(renderer, pixel.IM)

	//renderer:AlignText("left", "top")
	textBase = text.New(pixel.V(goldX+10, goldY), basicAtlas)
	textBase = text.New(pixel.V(goldX+10+getTextW(textBase, "GP :"), goldY), basicAtlas)
	fmt.Fprintln(textBase, "0")
	textBase.Draw(renderer, pixel.IM)

	textBase = text.New(pixel.V(goldX+10, goldY-25), basicAtlas)
	textBase = text.New(pixel.V(goldX+10+getTextW(textBase, "TIME :"), goldY-25), basicAtlas)
	fmt.Fprintln(textBase, "0")
	textBase.Draw(renderer, pixel.IM)

	// Party Members
	partyX := fm.Layout.Left("party") - 16
	partyY := fm.Layout.Top("party") - 45
	fm.PartyMenu.SetPosition(partyX, partyY)
	fm.PartyMenu.Render(renderer)
}

//get text Width
func getTextW(textBase *text.Text, txt string) float64 {
	return textBase.BoundsOf(txt).W()
}
