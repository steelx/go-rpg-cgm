package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/world"
	"math"
	"reflect"
)

type LootSummaryState struct {
	win        *pixelgl.Window
	Stack      *gui.StateStack
	CombatData CombatData
	World      *combat.WorldExtended
	Layout     gui.Layout
	Panels     []gui.Panel
	Loot       []world.ItemIndex
	Gold,
	GoldPerSec,
	GoldCounter float64
	IsCountingGold bool
	LootView       *gui.SelectionMenu
}

func LootSummaryStateCreate(stack *gui.StateStack, win *pixelgl.Window, world *combat.WorldExtended, combatData CombatData) *LootSummaryState {
	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.25, 2)
	layout.SplitHorz("top", "title", "detail", 0.55, 2)
	layout.SplitVert("detail", "left", "right", 0.5, 1)

	s := &LootSummaryState{
		win:            win,
		Stack:          stack,
		CombatData:     combatData,
		World:          world,
		Layout:         layout,
		Loot:           combatData.Loot,
		Gold:           combatData.Gold,
		GoldPerSec:     5.0,
		GoldCounter:    0,
		IsCountingGold: true,
	}

	s.Panels = []gui.Panel{
		s.Layout.CreatePanel("title"),
		s.Layout.CreatePanel("left"),
		s.Layout.CreatePanel("right"),
		s.Layout.CreatePanel("bottom"),
	}

	lootMenu := gui.SelectionMenuCreate(37, 175, 100,
		combatData.Loot,
		true,
		pixel.ZV,
		func(i int, itemIdx interface{}) {},
		s.RenderItem,
	)
	//lootMenu.MaxRows = 9
	//lootMenu.Columns = 3
	s.LootView = &lootMenu

	lootX := s.Layout.Left("bottom") + 10
	lootY := s.Layout.Top("bottom") - 35
	s.LootView.SetPosition(lootX, lootY)
	s.LootView.HideCursor()

	return s
}

func (s *LootSummaryState) RenderItem(a ...interface{}) {
	//renderer pixel.Target, x, y float64, itemIdx ItemIndex
	rendererV := reflect.ValueOf(a[0])
	renderer := rendererV.Interface().(pixel.Target)
	xV := reflect.ValueOf(a[1])
	x := xV.Interface().(float64)
	yV := reflect.ValueOf(a[2])
	y := yV.Interface().(float64)
	itemIdxV := reflect.ValueOf(a[3])
	itemIdx := itemIdxV.Interface().(world.ItemIndex)

	item := world.ItemsDB[itemIdx.Id]
	textStr := item.Name
	if itemIdx.Count > 1 {
		textStr = fmt.Sprintf("%s x %d", item.Name, itemIdx.Count)
	}

	textBase := text.New(pixel.V(x, y), gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, textStr)
	textBase.Draw(renderer, pixel.IM)
}

func (s *LootSummaryState) Enter() {
	s.IsCountingGold = true
	s.GoldCounter = 0

	// Add the items to the inventory.
	for _, v := range s.Loot {
		s.World.AddItem(v.Id, v.Count)
	}
}

func (s *LootSummaryState) Exit() {

}

func (s *LootSummaryState) Update(dt float64) bool {

	if s.IsCountingGold {
		s.GoldCounter = s.GoldCounter + (s.GoldPerSec * dt)
		goldToGive := math.Floor(s.GoldCounter)
		s.GoldCounter = s.GoldCounter - goldToGive
		s.Gold = s.Gold - goldToGive

		s.World.Gold = s.World.Gold + goldToGive

		if s.Gold == 0 {
			s.IsCountingGold = false
		}

	}

	return false //we dont want to update other states
}

func (s *LootSummaryState) Render(renderer *pixelgl.Window) {
	for _, v := range s.Panels {
		v.Draw(renderer)
	}

	titleX := s.Layout.MidX("title")
	titleY := s.Layout.MidY("title")
	textBase := text.New(pixel.V(titleX, titleY), gui.BasicAtlas12)
	fmt.Fprintln(textBase, "Found Loot!")
	textBase.Draw(renderer, pixel.IM)

	leftX := s.Layout.Left("left") + 12
	leftY := s.Layout.MidY("left")
	goldStr := fmt.Sprintf("Gold Found: %+6v Gold", s.Gold)
	textBase = text.New(pixel.V(leftX, leftY), gui.BasicAtlas12)
	fmt.Fprintln(textBase, goldStr)
	textBase.Draw(renderer, pixel.IM)

	rightX := s.Layout.Left("right") + 12
	rightY := leftY
	partyGPStr := fmt.Sprintf("Party Gold: %v gp", s.World.Gold)
	textBase = text.New(pixel.V(rightX, rightY), gui.BasicAtlas12)
	fmt.Fprintln(textBase, partyGPStr)
	textBase.Draw(renderer, pixel.IM)

	s.LootView.Render(renderer)

	//camera
	camera := pixel.IM.Scaled(pixel.ZV, 1.0).Moved(renderer.Bounds().Center().Sub(pixel.ZV))
	renderer.SetMatrix(camera)
}

func (s *LootSummaryState) HandleInput(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeySpace) {
		if s.Gold > 0 {
			s.SkipCountingGold()
			return
		}
		s.Stack.Pop()
		//pending new map after winning
		storyboardEvents := []interface{}{
			Wait(0),
			BlackScreen("blackscreen"),
			Wait(1),
			KillState("blackscreen"),
			//ReplaceState(s, combatState),//then dont Pop()
			ReplaceScene("handin", "map_sewer", 3, 5, false, win),
			PlayBGSound("../sound/reveal.mp3"),
			HandOffToMainStack("map_sewer"),
		}
		storyboard := StoryboardCreate(s.Stack, win, storyboardEvents, false)
		s.Stack.Push(storyboard)
	}
}

func (s *LootSummaryState) SkipCountingGold() {
	s.IsCountingGold = false
	s.GoldCounter = 0
	goldToGive := s.Gold
	s.Gold = 0
	s.World.Gold = s.World.Gold + goldToGive
}
