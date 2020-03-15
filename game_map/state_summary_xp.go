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

type XPSummaryState struct {
	win        *pixelgl.Window
	Stack      *gui.StateStack
	CombatData CombatData
	Layout     gui.Layout
	TitlePanels,
	ActorPanels []gui.Panel
	XP, //The experience points to give out.
	XPcopy,
	XPPerSec,
	XPCounter float64
	IsCountingXP  bool
	Party         []*combat.Actor
	PartySummary  []*combat.ActorXPSummary
	OnWinCallback func()
}

type CombatData struct {
	XP, Gold float64
	Loot     []world.ItemIndex
}

func XPSummaryStateCreate(stack *gui.StateStack, win *pixelgl.Window, party combat.Party, combatData CombatData, onWinCallback func()) *XPSummaryState {
	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.5, 2)
	layout.SplitHorz("top", "top", "one", 0.5, 2)
	layout.SplitHorz("bottom", "two", "three", 0.5, 2)

	layout.SplitHorz("top", "title", "detail", 0.5, 2)

	s := &XPSummaryState{
		win:           win,
		Stack:         stack,
		CombatData:    combatData,
		Layout:        layout,
		XP:            combatData.XP,
		XPcopy:        combatData.XP,
		XPPerSec:      5.0,
		XPCounter:     0,
		IsCountingXP:  true,
		Party:         party.ToArray(),
		OnWinCallback: onWinCallback,
	}

	digitNumber := math.Log10(s.XP + 1)
	s.XPPerSec = s.XPPerSec * digitNumber

	s.TitlePanels = []gui.Panel{
		s.Layout.CreatePanel("title"),
		s.Layout.CreatePanel("detail"),
	}
	s.ActorPanels = []gui.Panel{
		s.Layout.CreatePanel("one"),
		s.Layout.CreatePanel("two"),
		s.Layout.CreatePanel("three"),
	}

	s.PartySummary = make([]*combat.ActorXPSummary, 0)
	//summaryLeft := s.Layout.Left("detail") + 16
	index := 0
	panelIds := []string{"one", "two", "three"}

	for _, v := range s.Party {
		panelId := panelIds[index]
		summary := combat.ActorXPSummaryCreate(v, s.Layout, panelId)
		// summaryTop := s.Layout.Top(panelId)
		// summary.SetPosition(summaryLeft, summaryTop)
		s.PartySummary = append(s.PartySummary, summary)
		index++
	}

	return s
}

func (s *XPSummaryState) Enter() {
	s.IsCountingXP = true
	s.XPCounter = 0
}

func (s *XPSummaryState) Exit() {

}

func (s *XPSummaryState) Update(dt float64) bool {
	for _, v := range s.PartySummary {
		v.Update(dt)
	}

	if s.IsCountingXP {

		s.XPCounter = s.XPCounter + s.XPPerSec*dt
		xpToApply := math.Floor(s.XPCounter)
		s.XPCounter = s.XPCounter - xpToApply
		s.XP = s.XP - xpToApply

		s.ApplyXPToParty(xpToApply)

		if s.XP == 0 {
			s.IsCountingXP = false
		}

	}

	return false //we dont want to update other states
}

func (s *XPSummaryState) Render(renderer *pixelgl.Window) {
	for _, v := range s.TitlePanels {
		v.Draw(renderer)
	}

	titleX := s.Layout.MidX("title")
	titleY := s.Layout.MidY("title")
	pos := pixel.V(titleX, titleY)
	textBase := text.New(pos, gui.BasicAtlas12)
	titleStr := "Experience Increased!"
	titleStrW := textBase.BoundsOf(titleStr).W()
	pos = pixel.V(titleX+titleStrW, titleY)
	fmt.Fprintln(textBase, titleStr)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, 1.5))

	//XP
	detailX := s.Layout.Left("detail") + 16
	detailY := s.Layout.MidY("detail")
	pos = pixel.V(detailX, detailY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	detailStr := fmt.Sprintf("XP increased by %v.", s.XPcopy)
	fmt.Fprintln(textBase, detailStr)
	textBase.Draw(renderer, pixel.IM)

	for i := 0; i < len(s.PartySummary); i++ {
		s.ActorPanels[i].Draw(renderer)
		s.PartySummary[i].Render(renderer)
	}

	//camera
	camera := pixel.IM.Scaled(pixel.ZV, 1.0).Moved(renderer.Bounds().Center().Sub(pixel.ZV))
	renderer.SetMatrix(camera)
}

func (s *XPSummaryState) HandleInput(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeySpace) {

		if s.XP > 0 {
			s.SkipCountingXP()
			return
		}

		if s.ArePopUpsRemaining() {
			s.CloseNextPopUp()
			return
		}

		s.GotoLootSummary()
	}
}

func (s *XPSummaryState) UnlockPopUps(summary *combat.ActorXPSummary, levelUpActions map[string][]string) {
	for k, v := range levelUpActions {
		hexColor := "#6dff25"
		db := world.SpecialsDB
		if k == combat.ActionMagic {
			db = world.SpellsDB
			hexColor = "#b725ff"
		}

		for _, id := range v {
			msg := fmt.Sprintf("+ %s", db[id].Name)
			summary.AddPopUp(msg, hexColor)
		}
	}
}

func (s *XPSummaryState) ApplyXPToParty(xp float64) {
	for k, actor := range s.Party {
		if actor.Stats.Get("HpNow") > 0 {
			summary := s.PartySummary[k]
			actor.AddXP(xp)

			for actor.ReadyToLevelUp() {
				levelUp := actor.CreateLevelUp()
				levelNumber := actor.Level + levelUp.Level
				summary.AddPopUp(fmt.Sprintf("Level Up! %d", levelNumber), "#e9d79b")

				s.UnlockPopUps(summary, levelUp.Actions)

				actor.ApplyLevel(levelUp)
			}
		}
	}
}

func (s *XPSummaryState) SkipCountingXP() {
	s.IsCountingXP = false
	s.XPCounter = 0
	xpToApply := s.XP
	s.XP = 0
	s.ApplyXPToParty(xpToApply)
}

func (s XPSummaryState) ArePopUpsRemaining() bool {
	for _, v := range s.PartySummary {
		if v.PopUpCount() > 0 {
			return true
		}
	}
	return false
}

func (s *XPSummaryState) CloseNextPopUp() {
	for _, v := range s.PartySummary {
		if v.PopUpCount() > 0 {
			v.CancelPopUp()
		}
	}
}

func (s *XPSummaryState) GotoLootSummary() {
	world_ := reflect.ValueOf(s.Stack.Globals["world"]).Interface().(*combat.WorldExtended)
	lootSummaryState := LootSummaryStateCreate(s.Stack, s.win, world_, s.CombatData, s.OnWinCallback)

	storyboardEvents := []interface{}{
		Wait(0),
		BlackScreen("blackscreen"),
		Wait(0.3),
		KillState("blackscreen"),
		ReplaceState(s, lootSummaryState),
	}
	storyboard := StoryboardCreate(s.Stack, s.win, storyboardEvents, false)
	s.Stack.Push(storyboard)
}
