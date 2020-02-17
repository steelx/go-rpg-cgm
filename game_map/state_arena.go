package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"math"
	"reflect"
)

var rounds = []*ArenaRound{
	{Name: "Round 1", Locked: false, Enemies: []combat.ActorDef{
		combat.GoblinDef,
	}},
	{Name: "Round 2", Locked: true, Enemies: []combat.ActorDef{
		combat.GoblinDef,
		combat.GoblinDef,
	}},
	{Name: "Round 3", Locked: true, Enemies: []combat.ActorDef{
		combat.GoblinDef,
		combat.GoblinDef,
		combat.GoblinDef,
	}},
	{Name: "Round 4", Locked: true, Enemies: []combat.ActorDef{
		combat.OgreDef,
		combat.OgreDef,
	}},
	{Name: "Round 5", Locked: true, Enemies: []combat.ActorDef{
		combat.DragonDef,
	}},
}

type ArenaRound struct {
	Name    string
	Locked  bool
	Enemies []combat.ActorDef
}

type ArenaState struct {
	prevState gui.StackInterface
	Stack     *gui.StateStack
	World     *combat.WorldExtended
	Layout    gui.Layout
	Panels    []gui.Panel
	Selection *gui.SelectionMenu

	Rounds []*ArenaRound
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
		Rounds:    rounds,
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
		s.OnRoundSelected,
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
	round := reflect.ValueOf(a[3]).Interface().(*ArenaRound)

	lockLabel := "Open"
	if round.Locked {
		lockLabel = "Locked"
	}

	label := fmt.Sprintf("%s: %s", round.Name, lockLabel)
	textBase := text.New(pixel.V(x, y), gui.BasicAtlas12)
	fmt.Fprintf(textBase, label)
	textBase.Draw(renderer, pixel.IM)
}

func (s *ArenaState) OnRoundSelected(index int, itemI interface{}) {
	item := reflect.ValueOf(itemI).Interface().(*ArenaRound)
	if item.Locked {
		return
	}

	enemyDefs := []combat.ActorDef{combat.GoblinDef}
	if len(item.Enemies) > 0 {
		enemyDefs = item.Enemies
	}

	var enemyList []*combat.Actor
	for k, v := range enemyDefs {
		enemy_ := combat.ActorCreate(v, fmt.Sprintf("%v", k))
		enemyList = append(enemyList, &enemy_)
	}
	combatDef := CombatDef{
		Background: "../resources/arena_background.png",
		Actors: Actors{
			Party:   s.World.Party.ToArray(),
			Enemies: enemyList,
		},
		CanFlee: false,
		OnWin: func() {
			s.WinRound(index, item)
		},
		OnDie: func() {
			s.LoseRound(index, item)
		},
	}
	state := CombatStateCreate(s.Stack, s.Stack.Win, combatDef)
	s.Stack.Push(state)
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

func (s *ArenaState) WinRound(index int, round *ArenaRound) {
	fmt.Println("WinRound", round.Name)
	//Check for win - is is last round
	if index == len(s.Rounds)-1 {
		s.Stack.Pop()
		state := ArenaCompleteStateCreate(s.Stack)
		s.Stack.Push(state)
		return
	}

	//Move the cursor to the next round if there is one
	s.Selection.MoveDown()

	//Unlock the newly selected round
	nextRoundI := s.Selection.SelectedItem()
	nextRound := reflect.ValueOf(nextRoundI).Interface().(*ArenaRound)
	nextRound.Locked = false
}
func (s *ArenaState) LoseRound(index int, round *ArenaRound) {
	party := s.World.Party.Members
	for _, v := range party {
		hpNow := v.Stats.Get("HpNow")
		hpNow = math.Max(hpNow, 1)
		v.Stats.Set("HpNow", hpNow)
	}
}
