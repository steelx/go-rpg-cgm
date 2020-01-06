package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"math"
	"strings"
)

/* e.g.
tBox := TextboxCreate(
		"A nation can survive its fools, and even the ambitious. But it cannot survive treason from within. An enemy at the gates is less formidable, for he is known and carries his banner openly. But the traitor moves amongst those within the gate freely, his sly whispers rustling through all the alleys, heard in the very halls of government itself. For the traitor appears not a traitor; he speaks in accents familiar to his victims, and he wears their face and their arguments, he appeals to the baseness that lies deep in the hearts of all men. He rots the soul of a nation, he works secretly and unknown in the night to undermine the pillars of the city, he infects the body politic so that it can no longer resist. A murderer is less to fear. Jai Hind I Love India <3 ",
		basicAtlas12,
		PanelCreate(pic, pixel.V(-150, 200), 300, 100),
		continueCaretPng,
		"Ajinkya",
		avatarPng,
	)
*/

type Textbox struct {
	text                string
	textScale, size     float64
	Position            pixel.Vec
	textBounds          pixel.Rect
	textBase            *text.Text
	textAtlas           *text.Atlas
	mPanel              Panel
	continueMark        pixel.Picture
	Width, Height       float64
	textBlocks          []string
	textBlockLimitIndex int
	textRowLimit        int
	avatarName          string
	avatarImg           pixel.Picture
}

func TextboxCreate(txt string, textAtlas *text.Atlas, panel Panel, continueImage pixel.Picture, avatarName string, avatarImg pixel.Picture, withMenu bool) Textbox {
	t := Textbox{
		text:         txt,
		textScale:    1,
		size:         14,
		mPanel:       panel,
		textBounds:   panel.mBounds,
		continueMark: continueImage,
	}

	t.avatarName = avatarName
	t.avatarImg = avatarImg
	t.textAtlas = textAtlas

	t.makeTextColumns(withMenu)
	t.buildTextBlocks()

	return t
}

func (t *Textbox) makeTextColumns(withMenu bool) {
	var makeColumns bool
	if len(t.avatarName) != 0 {
		makeColumns = true
	}
	var textColumnWidth = t.mPanel.mBounds.W() - (t.size * 2)
	var textColumnHeight = t.mPanel.mBounds.H() - (t.size * 2)
	var topLeft, _, _, _ = t.mPanel.GetCorners()
	var textPos = pixel.V(topLeft.X+t.size, topLeft.Y-t.size)
	if makeColumns {
		textColumnWidth -= t.avatarImg.Bounds().W()
		textColumnHeight -= t.size
		textPos.X += t.avatarImg.Bounds().W() - t.size/2
		textPos.Y -= t.size / 2
	}
	if withMenu {
		textColumnHeight = textColumnHeight / 2
	}

	t.Position = textPos
	t.Width = textColumnWidth
	t.Height = textColumnHeight
}

func (t *Textbox) buildTextBlocks() {
	t.textBase = text.New(t.Position, t.textAtlas)
	t.textBase.LineHeight = t.size
	t.textBlocks = make([]string, 0) //imp to avoid null reference
	blocks := math.Abs(t.textBase.BoundsOf(t.text).W() / t.Width)
	eachBlockWidth := math.Ceil(t.textBase.BoundsOf(t.text).W()) / blocks
	splitTextAt := math.Ceil(eachBlockWidth / (t.size))

	var tempTxtLine = ""
	ss := strings.Fields(t.text)
	for _, word := range ss {
		tempTxtLine += word + " "
		if len(tempTxtLine) >= int(splitTextAt) {
			t.textBlocks = append(t.textBlocks, tempTxtLine)
			tempTxtLine = ""
		}
	}
}

func (t *Textbox) DrawTextWithPanel() {
	t.textBase.Clear()
	//limit prints
	eachBlockHeight := math.Abs(t.textBase.BoundsOf(t.text).H())
	t.textRowLimit = int(math.Ceil(t.Height / eachBlockHeight))
	lastIndex := minInt(t.textBlockLimitIndex+t.textRowLimit, len(t.textBlocks))
	firstIndex := t.textBlockLimitIndex

	if t.textBlockLimitIndex >= len(t.textBlocks) {
		//reached limit
		t.textBase.Clear()
		return
	}
	readFrom := t.textBlocks[firstIndex:lastIndex]

	for _, line := range readFrom {
		_, err := fmt.Fprintln(t.textBase, line)
		panicIfErr(err)
	}

	t.mPanel.Draw()
	t.textBase.Draw(global.gWin, pixel.IM)

	t.drawAvatar()
	t.drawContinueArrow()
}
func (t Textbox) drawAvatar() {
	if len(t.avatarName) == 0 {
		return
	}

	avatarSprite := pixel.NewSprite(t.avatarImg, t.avatarImg.Bounds())
	topLeft := pixel.V(
		t.mPanel.mBounds.Min.X+(t.avatarImg.Bounds().W()/2)+t.size/2,
		t.mPanel.mBounds.Max.Y-(t.avatarImg.Bounds().H()/2))

	titlePos := pixel.V(t.mPanel.mBounds.Min.X+t.size, t.mPanel.mBounds.Min.Y+t.avatarImg.Bounds().H()-t.size/2)
	title := text.New(titlePos, t.textAtlas)
	fmt.Fprintln(title, t.avatarName)

	title.Draw(global.gWin, pixel.IM.Scaled(titlePos, 0.8))
	avatarSprite.Draw(global.gWin, pixel.IM.Moved(topLeft).Scaled(topLeft, 0.9))
}
func (t Textbox) drawContinueArrow() {
	if t.textBlockLimitIndex+t.textRowLimit < len(t.textBlocks) {
		bottomRight := pixel.V(t.mPanel.mBounds.Max.X-t.size, t.mPanel.mBounds.Min.Y+t.size)
		sprite := pixel.NewSprite(t.continueMark, t.continueMark.Bounds())
		sprite.Draw(global.gWin, pixel.IM.Moved(bottomRight))
	}
}

func (t *Textbox) Next() {
	t.textBlockLimitIndex += t.textRowLimit
}

func (t *Textbox) HandleInput() {
	if global.gWin.JustPressed(pixelgl.KeySpace) {
		t.Next()
	}
}
