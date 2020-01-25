package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/animation"
	"github.com/steelx/go-rpg-cgm/utilz"
	"golang.org/x/image/font/basicfont"
	"math"
	"strings"
)

/* e.g.
tBox := TextboxCreateFixed(
		"A nation can survive its fools, and even the ambitious. But it cannot survive treason from within. An enemy at the gates is less formidable, for he is known and carries his banner openly. But the traitor moves amongst those within the gate freely, his sly whispers rustling through all the alleys, heard in the very halls of government itself. For the traitor appears not a traitor; he speaks in accents familiar to his victims, and he wears their face and their arguments, he appeals to the baseness that lies deep in the hearts of all men. He rots the soul of a nation, he works secretly and unknown in the night to undermine the pillars of the city, he infects the body politic so that it can no longer resist. A murderer is less to fear. Jai Hind I Love India <3 ",
		pixel.V(-150, 200), 300, 100,
		"Ajinkya",
		avatarPng,
	)
*/
var (
	continueCaretPng pixel.Picture
	cursorPng        pixel.Picture
	basicAtlas12     *text.Atlas
	basicAtlasAscii  = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

func init() {
	var err error
	continueCaretPng, err = utilz.LoadPicture("../resources/continue_caret.png")
	utilz.PanicIfErr(err)
	cursorPng, err = utilz.LoadPicture("../resources/cursor.png")
	utilz.PanicIfErr(err)

	fontFace12, err := utilz.LoadTTF("../resources/font/joystix.ttf", 12)
	utilz.PanicIfErr(err)
	basicAtlas12 = text.NewAtlas(fontFace12, text.ASCII)
}

type Textbox struct {
	Stack                       *StateStack
	text                        string
	textScale, size, topPadding float64
	Position                    pixel.Vec
	textBounds                  pixel.Rect
	textBase                    *text.Text
	textAtlas                   *text.Atlas
	mPanel                      Panel
	continueMark                pixel.Picture
	Width, Height               float64
	textBlocks                  []string
	textBlockLimitIndex         int
	textRowLimit                int
	avatarName                  string
	avatarImg                   pixel.Picture
	AppearTween                 animation.Tween
	time                        float64
	isFixed, isDead, hasMenu    bool
	menu                        SelectionMenu
}

func TextboxNew(stack *StateStack, txt string, size float64, atlas *text.Atlas, avatarName string, avatarImg pixel.Picture) Textbox {
	return Textbox{
		Stack:        stack,
		text:         txt,
		textScale:    1,
		size:         size,
		continueMark: continueCaretPng,
		avatarName:   avatarName,
		avatarImg:    avatarImg,
		textAtlas:    atlas,
		time:         0,
	}
}

func TextboxWithMenuCreate(stack *StateStack, textBoxText string, panelPos pixel.Vec, panelWidth, panelHeight float64,
	choices []string, onSelection func(int, string), showColumns bool) *Textbox {

	textbox := TextboxCreateFixed(
		stack,
		textBoxText,
		panelPos, panelWidth, panelHeight,
		"",
		nil,
		true,
	)

	textBounds := textbox.getTextBound()

	textbox.menu = SelectionMenuCreate(24, 128, choices, showColumns,
		pixel.V(textbox.Position.X-10, textbox.Position.Y-textBounds.H()-10), func(i int, s string) {
			onSelection(i, s)
			textbox.isDead = true
		}, nil)

	return &textbox
}

func TextboxFITMenuCreate(stack *StateStack, x, y float64, textBoxText string, choices []string, onSelection func(int, string)) *Textbox {
	panelPos := pixel.V(x, y)
	t := TextboxNew(stack, textBoxText, 14, basicAtlas12, "", nil)
	t.AppearTween = animation.TweenCreate(1, 1, 1)
	t.isFixed = false

	t.hasMenu = true
	if t.hasMenu {
		t.topPadding = 10
	}
	fmt.Println("choices", choices)
	textBounds := t.getTextBound()
	menu := SelectionMenuCreate(24, 128, choices, true,
		pixel.V(t.Position.X, t.Position.Y-textBounds.H()-10), func(i int, s string) {
			onSelection(i, s)
			t.isDead = true
		}, nil)

	panel := PanelCreate(panelPos, menu.GetWidth(), menu.GetHeight())

	t.menu = menu
	t.mPanel = panel
	t.textBounds = panel.mBounds

	t.makeTextColumns()
	t.buildTextBlocks()

	return &t
}

func TextboxCreateFixed(stack *StateStack, txt string, panelPos pixel.Vec, panelWidth, panelHeight float64, avatarName string, avatarImg pixel.Picture, hasMenu bool) Textbox {
	panel := PanelCreate(panelPos, panelWidth, panelHeight)
	t := TextboxNew(stack, txt, 14, basicAtlas12, avatarName, avatarImg)
	t.AppearTween = animation.TweenCreate(1, 0, 1)
	t.isFixed = true
	t.mPanel = panel
	t.textBounds = panel.mBounds
	t.hasMenu = hasMenu
	if hasMenu {
		t.topPadding = 10
	}

	t.makeTextColumns()
	t.buildTextBlocks()

	return t
}

//TextboxCreateFitted are good for small chats
// height and width gets set automatically
func TextboxCreateFitted(stack *StateStack, txt string, panelPos pixel.Vec, hasMenu bool) Textbox {
	const padding = 20.0

	tBox := TextboxNew(stack, txt, 13, basicAtlasAscii, "", nil)
	tBox.AppearTween = animation.TweenCreate(0.9, 1, 0.3)
	tBox.textBase = text.New(panelPos, tBox.textAtlas)
	tBox.textBase.LineHeight = padding
	textBounds := tBox.getTextBound()
	tBox.hasMenu = hasMenu
	panelHeight := textBounds.H() + padding
	if hasMenu {
		panelHeight = panelHeight + 100 //+menu height
	}
	panel := PanelCreate(panelPos, textBounds.W()+(padding*2), panelHeight)
	topLeft, _, _, _ := panel.GetCorners()
	textPos := pixel.V(topLeft.X+padding, topLeft.Y-padding)
	tBox.textBase = text.New(textPos, tBox.textAtlas) //reset text position to bounds
	tBox.mPanel = panel

	tBox.makeTextColumns()

	return tBox
}

func (t *Textbox) makeTextColumns() {
	var makeColumns bool
	if len(t.avatarName) != 0 {
		makeColumns = true
	}
	var textColumnWidth = t.mPanel.mBounds.W() - (t.size * 2)
	var textColumnHeight = t.mPanel.mBounds.H() - (t.size * 2)
	var topLeft, _, _, _ = t.mPanel.GetCorners()
	var textPos = pixel.V(topLeft.X+t.size, topLeft.Y-t.size-t.topPadding)
	if makeColumns {
		textColumnWidth -= t.avatarImg.Bounds().W()
		textColumnHeight -= t.size
		textPos.X += t.avatarImg.Bounds().W() - t.size/2
		textPos.Y -= t.size / 2
	}
	if t.hasMenu {
		textColumnHeight = textColumnHeight / 2
	}

	t.Position = textPos
	t.Width = textColumnWidth
	t.Height = textColumnHeight
}

func (t Textbox) getTextBound() pixel.Rect {
	if t.textBase == nil {
		t.textBase = text.New(t.Position, t.textAtlas)
	}
	return t.textBase.BoundsOf(t.text)
}

func (t *Textbox) buildTextBlocks() {
	t.textBase = text.New(t.Position, t.textAtlas)
	t.textBase.LineHeight = t.size
	t.textBlocks = make([]string, 0) //imp to avoid null reference
	blocks := math.Abs(t.getTextBound().W() / t.Width)
	eachBlockWidth := math.Ceil(t.textBase.BoundsOf(t.text).W()) / blocks
	splitTextAt := math.Ceil(eachBlockWidth / (t.size))

	var tempTxtLine = ""
	ss := strings.Fields(t.text)
	for i, word := range ss {
		tempTxtLine += word + " "
		if len(ss)-1 == i || (len(tempTxtLine) > int(splitTextAt)) {
			t.textBlocks = append(t.textBlocks, tempTxtLine)
			tempTxtLine = ""
		}
	}
}

func (t *Textbox) IsDead() bool {
	return t.isDead
}

func (t Textbox) HasReachedLimit() bool {
	return t.textBlockLimitIndex >= len(t.textBlocks)
}

func (t Textbox) drawAvatar(renderer pixel.Target) {
	if len(t.avatarName) == 0 {
		return
	}

	avatarSprite := pixel.NewSprite(t.avatarImg, t.avatarImg.Bounds())
	topLeft := pixel.V(
		t.mPanel.mBounds.Min.X+(t.avatarImg.Bounds().W()/2)+t.size/2,
		t.mPanel.mBounds.Max.Y-(t.avatarImg.Bounds().H()/2)-5)

	titlePos := pixel.V(t.mPanel.mBounds.Min.X+t.size, t.mPanel.mBounds.Min.Y+t.avatarImg.Bounds().H()-(t.size/2)-2)

	title := text.New(titlePos, basicAtlasAscii)
	fmt.Fprintln(title, t.avatarName)

	title.Draw(renderer, pixel.IM.Scaled(titlePos, 1))
	avatarSprite.Draw(renderer, pixel.IM.Moved(topLeft).Scaled(topLeft, 0.9))
}
func (t Textbox) drawContinueArrow(renderer pixel.Target) {
	if t.textBlockLimitIndex+t.textRowLimit < len(t.textBlocks) {
		mat := pixel.IM
		bottomRight := pixel.V(t.mPanel.mBounds.Max.X-t.size, t.mPanel.mBounds.Min.Y+t.size)
		sprite := pixel.NewSprite(t.continueMark, t.continueMark.Bounds())
		sprite.Draw(renderer, mat.Moved(bottomRight))

		title := text.New(bottomRight, basicAtlasAscii)
		keyHintTxt, padding := "spacebar", 20.0
		textPos := bottomRight.Sub(
			pixel.V(title.BoundsOf(keyHintTxt).W(), -padding),
		)
		fmt.Fprintln(title, keyHintTxt)
		title.Draw(renderer, pixel.IM.Moved(textPos))
	}
}

func (t *Textbox) Next() {
	t.textBlockLimitIndex += t.textRowLimit
}

func (t *Textbox) renderFitted(renderer pixel.Target) {
	scale := t.AppearTween.Value()
	t.textBase.Clear()
	t.mPanel.Draw(renderer)
	fmt.Fprintln(t.textBase, t.text)
	t.textBase.Draw(renderer, pixel.IM.Scaled(t.Position, scale))
}

//RenderWithPanel will render text based on Panel height and divide rows
//based on available height, it will destroy at the end of Next last user input
func (t *Textbox) renderFixed(renderer pixel.Target) {
	t.textBase.Clear()
	//limit prints
	eachBlockHeight := math.Abs(t.textBase.BoundsOf(t.text).H())
	t.textRowLimit = int(math.Ceil(t.Height / eachBlockHeight))
	lastIndex := utilz.MinInt(t.textBlockLimitIndex+t.textRowLimit, len(t.textBlocks))
	firstIndex := t.textBlockLimitIndex

	if t.HasReachedLimit() {
		//reached limit
		t.isDead = true
		return
	}
	readFrom := t.textBlocks[firstIndex:lastIndex]

	for _, line := range readFrom {
		_, err := fmt.Fprintln(t.textBase, line)
		utilz.PanicIfErr(err)
	}

	t.mPanel.Draw(renderer)
	t.textBase.Draw(renderer, pixel.IM)

	t.drawAvatar(renderer)
	t.drawContinueArrow(renderer)
}

func (t *Textbox) Render(renderer *pixelgl.Window) {

	if t.hasMenu {
		t.renderFitted(renderer)
		t.menu.Render(renderer)
		return
	}

	if t.isFixed {
		t.renderFixed(renderer)
	} else {
		t.renderFitted(renderer)
	}
}

func (t *Textbox) Update(dt float64) bool {
	t.time = t.time + dt
	t.AppearTween.Update(dt)
	if t.IsDead() {
		t.Stack.Pop()
	}
	return t.IsDead()
}

func (t *Textbox) Enter() {

}

func (t *Textbox) Exit() {

}

//HandleInput takes care of 3 types of textbox's
//1 textbox with menu
//2 Fixed then we let the blocks render
//3 Fitted users marks as read and goes to next text popup
func (t *Textbox) HandleInput(window *pixelgl.Window) {
	if t.hasMenu {
		t.menu.HandleInput(window)
	}
	if window.JustPressed(pixelgl.KeySpace) {
		if t.isFixed {
			t.Next()
			return
		}

		t.OnClick()
	}
}

func (t *Textbox) OnClick() {
	t.AppearTween = animation.TweenCreate(1, 0, 0.2)
	t.isDead = true
}
