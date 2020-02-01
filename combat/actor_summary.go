package combat

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
)

type ActorSummary struct {
	X, Y, Width  float64
	Actor        Actor
	HPBar, MPBar gui.ProgressBarIMD
	XPBar        gui.ProgressBar
	ShowXP       bool
	AvatarTextPadding,
	LabelRightPadding,
	LabelValuePadding,
	VerticalPadding,
	TextPaddingY float64
}

func ActorSummaryCreate(actor Actor, showXP bool) ActorSummary {

	s := ActorSummary{
		X: 0, Y: 0, Width: 380, Actor: actor, ShowXP: showXP,
		HPBar: gui.ProgressBarIMDCreate(
			0, 0,
			actor.Stats.Get("HpNow"),
			actor.Stats.Get("HpMax"),
			"#FF001E",
			"#15FF00",
			3, 100,
			nil,
		),
		MPBar: gui.ProgressBarIMDCreate(
			0, 0,
			actor.Stats.Get("MpNow"),
			actor.Stats.Get("MpMax"),
			"#A48B2C",
			"#00E7DA",
			3, 100,
			nil,
		),
		AvatarTextPadding: 14,
		LabelRightPadding: 12,
		LabelValuePadding: 18,
		VerticalPadding:   20,
		TextPaddingY:      10,
	}

	if s.ShowXP {
		s.XPBar = gui.ProgressBarCreate(
			0, 0,
			actor.XP,
			actor.XP+actor.NextLevelXP,
			"../resources/progressbar_bg.png",
			"../resources/progressbar_fg.png",
		)
	}

	s.SetPosition(s.X, s.Y)

	return s
}

//SetPosition
func (s *ActorSummary) SetPosition(x, y float64) {
	s.X = x
	s.Y = y - s.AvatarTextPadding

	if s.ShowXP {
		boxRight := s.X + s.Width
		barX := boxRight + s.XPBar.HalfWidth
		barY := s.Y - s.TextPaddingY - 15
		s.XPBar.SetPosition(barX, barY)
	}

	// HP & MP
	avatarW := s.Actor.PortraitTexture.Bounds().W()
	barX := s.X + avatarW
	barX = barX + s.LabelRightPadding + s.LabelValuePadding
	barX = barX + s.MPBar.HalfWidth

	s.HPBar.SetPosition(barX, s.Y-55)
	s.MPBar.SetPosition(barX, s.Y-75)
}

func (s ActorSummary) GetCursorPosition() pixel.Vec {
	return pixel.V(s.X, s.Y-40)
}

func (s *ActorSummary) HideXP() {
	s.ShowXP = false
}

func (s *ActorSummary) Render(renderer pixel.Target) {
	actor := s.Actor

	// Position avatar image from top left
	avatar := actor.Portrait
	avatarW := actor.PortraitTexture.Bounds().W()
	avatarH := actor.PortraitTexture.Bounds().H()
	avatarX := s.X + avatarW
	avatarY := s.Y - avatarH/2
	avatar.Draw(renderer, pixel.IM.Moved(pixel.V(avatarX, avatarY)))

	// Position basic stats to the left of the avatar
	textX := avatarX + avatarW/2 + s.AvatarTextPadding
	textY := s.Y - s.TextPaddingY
	pos := pixel.V(textX, textY)
	textBase := text.New(pos, gui.BasicAtlas12)
	fmt.Fprintln(textBase, actor.Name)
	textBase.Draw(renderer, pixel.IM)

	//Draw LVL, HP and MP labels
	textX = textX + s.LabelRightPadding
	textY = textY - 20
	statsStartY := textY
	pos = pixel.V(textX, textY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, "LV")
	textBase.Draw(renderer, pixel.IM)

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, "HP")
	textBase.Draw(renderer, pixel.IM)

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, "MP")
	textBase.Draw(renderer, pixel.IM)

	// Fill in the values
	textY = statsStartY
	textX = textX + s.LabelValuePadding
	level := actor.Level
	pos = pixel.V(textX, textY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, level)
	textBase.Draw(renderer, pixel.IM)

	hp := actor.Stats.Get("HpNow")
	maxHP := actor.Stats.Get("HpMax")
	mp := actor.Stats.Get("MpNow")
	maxMP := actor.Stats.Get("MpMax")

	hpTxt := fmt.Sprintf("%v/%v", hp, maxHP)
	mpTxt := fmt.Sprintf("%v/%v", mp, maxMP)

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, hpTxt)
	textBase.Draw(renderer, pixel.IM)

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	textBase = text.New(pos, gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, mpTxt)
	textBase.Draw(renderer, pixel.IM)

	// Next Level area
	if s.ShowXP {
		s.XPBar.Render(renderer)

		boxRight := s.X + s.Width
		textY := s.Y - s.TextPaddingY
		left := boxRight + s.XPBar.HalfWidth

		pos = pixel.V(left, textY)
		textBase = text.New(pos, gui.BasicAtlasAscii)
		fmt.Fprintln(textBase, "Next Level")
		textBase.Draw(renderer, pixel.IM)
	}

	// MP & HP bars
	s.HPBar.Render(renderer)
	s.MPBar.Render(renderer)
}
