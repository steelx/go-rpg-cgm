package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"golang.org/x/image/font/basicfont"
)

type ActorSummary struct {
	X, Y, Width         float64
	Actor               combat.Actor
	HPBar, MPBar, XPBar gui.ProgressBar
	ShowXP              bool
	AvatarTextPadding,
	LabelRightPadding,
	LabelValuePadding,
	VerticalPadding float64
}

func ActorSummaryCreate(actor combat.Actor, showXP bool) ActorSummary {

	s := ActorSummary{
		X: 0, Y: 0, Width: 340, Actor: actor, ShowXP: showXP,
		HPBar:             gui.ProgressBarCreate(0, 0, actor.Stats.Get("HpNow"), actor.Stats.Get("HpMax")),
		MPBar:             gui.ProgressBarCreate(0, 0, actor.Stats.Get("MpNow"), actor.Stats.Get("MpMax")),
		AvatarTextPadding: 14,
		LabelRightPadding: 15,
		LabelValuePadding: 8,
		VerticalPadding:   18,
	}

	if s.ShowXP {
		s.XPBar = gui.ProgressBarCreate(0, 0, actor.XP, actor.NextLevelXP)
	}

	s.SetPosition(s.X, s.Y)

	return s
}

//SetPosition also updates party members Panel positions
func (s *ActorSummary) SetPosition(x, y float64) {
	s.X = x
	s.Y = y

	if s.ShowXP {
		boxRight := s.X + s.Width
		barX := boxRight - s.XPBar.HalfWidth
		barY := s.Y - 44
		s.XPBar.SetPosition(barX, barY)
	}

	// HP & MP
	avatarW := s.Actor.PortraitTexture.Bounds().W()
	barX := s.X + avatarW + s.AvatarTextPadding
	barX = barX + s.LabelRightPadding + s.LabelValuePadding
	barX = barX + s.MPBar.HalfWidth

	s.HPBar.SetPosition(barX, s.Y-72)
	s.MPBar.SetPosition(barX, s.Y-54)
}

func (s ActorSummary) GetCursorPosition() pixel.Vec {
	return pixel.V(s.X, s.Y-40)
}

func (s *ActorSummary) Render(renderer pixel.Target) {
	actor := s.Actor

	// Position avatar image from top left
	avatar := actor.Portrait
	avatarW := actor.PortraitTexture.Bounds().W()
	avatarH := actor.PortraitTexture.Bounds().H()
	avatarX := s.X + avatarW/2
	avatarY := s.Y - avatarH/2
	avatar.Draw(renderer, pixel.IM.Moved(pixel.V(avatarX, avatarY)))

	basicAtlasAscii := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textBase := text.New(pixel.V(0, 0), basicAtlasAscii)

	// Position basic stats to the left of the avatar
	textPadY := 2.0
	textX := avatarX + avatarW/2 + s.AvatarTextPadding
	textY := s.Y - textPadY
	pos := pixel.V(textX, textY)
	fmt.Fprintln(textBase, actor.Name)
	textBase.Draw(renderer, pixel.IM.Scaled(pos, 1.6).Moved(pos))

	//Draw LVL, HP and MP labels
	textX = textX + s.LabelRightPadding
	textY = textY - 20
	statsStartY := textY
	pos = pixel.V(textX, textY)
	fmt.Fprintln(textBase, "LV")
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	fmt.Fprintln(textBase, "HP")
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	fmt.Fprintln(textBase, "MP")
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	// Fill in the values
	textY = statsStartY
	textX = textX + s.LabelValuePadding
	level := actor.Level
	hp := actor.Stats.Get("HpNow")
	maxHP := actor.Stats.Get("HpMax")
	mp := actor.Stats.Get("MpNow")
	maxMP := actor.Stats.Get("MpMax")

	hpTxt := fmt.Sprintf("%v/%v", hp, maxHP)
	mpTxt := fmt.Sprintf("%v/%v", mp, maxMP)

	pos = pixel.V(textX, textY)
	fmt.Fprintln(textBase, level)
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	fmt.Fprintln(textBase, hpTxt)
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	textY = textY - s.VerticalPadding
	pos = pixel.V(textX, textY)
	fmt.Fprintln(textBase, mpTxt)
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	// Next Level area
	if s.ShowXP {
		s.XPBar.Render(renderer)

		boxRight := s.X + s.Width
		textY := statsStartY
		left := boxRight - s.XPBar.HalfWidth*2

		pos = pixel.V(left, textY)
		fmt.Fprintln(textBase, "Next Level")
		textBase.Draw(renderer, pixel.IM.Moved(pos))
	}

	// MP & HP bars
	s.HPBar.Render(renderer)
	s.MPBar.Render(renderer)
}
