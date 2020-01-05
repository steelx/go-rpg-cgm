package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"math"
	"strings"
)

type Textbox struct {
	text                string
	textScale, size     float64
	position            pixel.Vec
	textBounds          pixel.Rect
	textBase            *text.Text
	mPanel              Panel
	mContinueMark       pixel.Picture
	mWidth, mHeight     float64
	textBlocks          []string
	textBlockLimitIndex int
	textRowLimit        int
}

func TextboxCreate(txt string, textAtlas *text.Atlas, panel Panel, continueImage pixel.Picture) Textbox {
	t := Textbox{
		text:          txt,
		textScale:     1,
		size:          14,
		mPanel:        panel,
		textBounds:    panel.mBounds,
		mContinueMark: continueImage,
	}

	topLeft, _, _, _ := t.mPanel.GetCorners()
	t.position = pixel.V(topLeft.X+t.size, topLeft.Y-t.size)
	t.mWidth = panel.mBounds.W() - (t.size * 2)
	t.mHeight = panel.mBounds.H() - (t.size * 2)
	t.textBase = text.New(t.position, textAtlas)
	t.textBlocks = make([]string, 0)
	t.buildTextBlocks()

	return t
}

func (t *Textbox) buildTextBlocks() {
	t.textBase.LineHeight = 14
	blocks := math.Abs(t.textBase.BoundsOf(t.text).W() / t.mWidth)
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

func (t *Textbox) Draw() {
	t.textBase.Clear()
	//limit prints
	eachBlockHeight := math.Abs(t.textBase.BoundsOf(t.text).H())
	t.textRowLimit = int(math.Ceil(t.mHeight / eachBlockHeight))
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

	if t.textBlockLimitIndex+t.textRowLimit < len(t.textBlocks) {
		bottomCenter := pixel.V(t.textBounds.Center().X, t.textBounds.Center().Y-(t.mHeight/2))
		sprite := pixel.NewSprite(t.mContinueMark, t.mContinueMark.Bounds())
		sprite.Draw(global.gWin, pixel.IM.Moved(bottomCenter))
	}
}

func (t *Textbox) Next() {
	t.textBlockLimitIndex += t.textRowLimit
}
