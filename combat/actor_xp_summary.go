package combat

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
)

/* example ->
layout := gui.LayoutCreate(0, 0, win)
pos := pixel.V(400, 200)
layout.Panels["test"] = gui.PanelDef{pos, 400, 400}
actorHero := gWorld.Party.Members["hero"]
summary := combat.ActorXPSummaryCreate(actorHero, layout, "test")
summary.AddPopUp("Level Up!", "#34df6b")
*/

type ActorXPSummary struct {
	Actor                     *Actor
	X, Y                      float64
	Avatar                    *pixel.Sprite
	AvatarWidth, AvatarHeight float64
	Layout                    gui.Layout
	LayoutId                  string
	XPBar                     gui.ProgressBarIMD
	PopUpList                 []*gui.XPPopUp
	PopUpDisplayTime          float64
}

func ActorXPSummaryCreate(actor *Actor, layout gui.Layout, layoutId string) *ActorXPSummary {

	return &ActorXPSummary{
		Actor:        actor,
		Avatar:       actor.Portrait,
		AvatarWidth:  actor.PortraitTexture.Bounds().W(),
		AvatarHeight: actor.PortraitTexture.Bounds().H(),
		Layout:       layout,
		LayoutId:     layoutId,
		XPBar: gui.ProgressBarIMDCreate(
			0, 0,
			actor.XP,
			actor.XP+actor.NextLevelXP,
			"#A48B2C",
			"#00E7DA",
			3, 100,
			nil,
		),
	}
}

func (a *ActorXPSummary) SetPosition(x, y float64) {
	a.X = x
	a.Y = y
}

func (a ActorXPSummary) PopUpCount() int {
	return len(a.PopUpList)
}

func (a *ActorXPSummary) CancelPopUp() {
	if a.PopUpCount() == 0 {
		return
	}
	if popup := a.PopUpList[0]; popup != nil {
		popup.TurnOff()
	}
}

func (a *ActorXPSummary) removePopUpAtIndex(arr []*gui.XPPopUp, i int) []*gui.XPPopUp {
	return append(arr[:i], arr[i+1:]...)
}

func (a *ActorXPSummary) AddPopUp(text, hexColor string) {
	x := a.Layout.MidX(a.LayoutId)
	y := a.Layout.Top(a.LayoutId)
	popup := gui.XPPopUpCreate(x, y, text, hexColor)
	popup.TurnOn()
	a.PopUpList = append(a.PopUpList, popup)
}

func (a *ActorXPSummary) Update(dt float64) {
	if a.PopUpCount() == 0 {
		return
	}
	popup := a.PopUpList[0]
	if popup == nil {
		return
	}
	if popup.IsFinished() {
		a.PopUpList = a.removePopUpAtIndex(a.PopUpList, 0)
		return
	}
	popup.Update(dt)
	if popup.DisplayTime > a.PopUpDisplayTime && a.PopUpCount() > 1 {
		popup.TurnOff()
	}
}

func (a *ActorXPSummary) Render(renderer pixel.Target) {
	// portrait
	left := a.Layout.Left(a.LayoutId) + 25
	topY := a.Layout.Top(a.LayoutId) - a.AvatarHeight/2 + 25

	avatarLeft := left + a.AvatarWidth/2
	avatarY := topY - a.AvatarHeight/2
	pos := pixel.V(avatarLeft, avatarY)
	a.Avatar.Draw(renderer, pixel.IM.Moved(pos))

	// Name
	nameX := left + a.AvatarWidth + 20
	nameY := topY - 12
	pos = pixel.V(nameX, nameY)
	textBase := text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, a.Actor.Name)
	textBase.Draw(renderer, pixel.IM)

	// Level
	strLevel := fmt.Sprintf("Level: %+6v", a.Actor.Level)
	levelY := nameY - 42
	pos = pixel.V(nameX, levelY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, strLevel)
	textBase.Draw(renderer, pixel.IM)

	// XP
	strXPValue := fmt.Sprintf("EXP: %+2v", a.Actor.XP)
	right := a.Layout.Right(a.LayoutId) - a.AvatarWidth
	rightLabel := right - 96
	pos = pixel.V(rightLabel, nameY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, strXPValue)
	textBase.Draw(renderer, pixel.IM)

	//XPBar
	barX := right - a.XPBar.HalfWidth
	a.XPBar.SetPosition(barX, nameY-24)
	a.XPBar.SetValue(a.Actor.XP)
	a.XPBar.SetMax(a.Actor.XP + a.Actor.NextLevelXP)
	a.XPBar.Render(renderer)

	strNextLevel := fmt.Sprintf("Next Level XP: %+6v", a.Actor.NextLevelXP)
	pos = pixel.V(rightLabel, levelY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, strNextLevel)
	textBase.Draw(renderer, pixel.IM)

	if a.PopUpCount() == 0 {
		return
	}
	if popup := a.PopUpList[0]; popup != nil {
		popup.Render(renderer)
	}
}
