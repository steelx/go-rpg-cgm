package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/utilz"
	"golang.org/x/image/font/basicfont"
	"math"
	"reflect"
)

//var (
//	cursorPng pixel.Picture
//	BasicAtlasAscii *text.Atlas
//)
//
//func init() {
//	var err error
//	cursorPng, err = utilz.LoadPicture("../resources/cursor.png")
//	utilz.PanicIfErr(err)
//	BasicAtlasAscii = text.NewAtlas(basicfont.Face7x13, text.ASCII)
//}

/* e.g.
menu2 := gui.SelectionMenuCreate(24, 128,[]string{"Menu 1", "", "Menu 2", "Menu 03", "Menu 04", "Menu 05"}, false, pixel.V(200, 250), func(i int, item string) {
		fmt.Println(i, item)
	})
*/

type SelectionMenu struct {
	X, Y          float64
	width, height float64

	columns        int //The number of columns the menu has. This defaults to 1
	focusX, focusY int //Indicates which item in the list is currently selected.
	//focusX tells us which column is selected and focusY which element in that column
	SpacingY, SpacingX        float64 //space btw each items
	scale                     float64 //menu scale in size
	cursor                    *pixel.Sprite
	useCursorPos              bool
	cursorPosOffset           pixel.Vec
	cursorWidth, cursorHeight float64
	IsShowCursor              bool
	maxRows, displayRows      int //rows might be 30 but only 5 maxRows are displayed at once
	displayStart              int //index at which we start displaying menu, e.g. out of 30 max 5 are visible from index 6
	textBase                  *text.Text
	OnSelection               func(int, interface{}) //to be called after selection
	DataI                     []interface{}
	RenderFunction            func(a ...interface{})
}

func SelectionMenuCreate(spacingY, spacingX, xWidth float64, data interface{}, showColumns bool, position pixel.Vec, onSelection func(int, interface{}), renderFunc func(a ...interface{})) SelectionMenu {
	//xWidth should be passed if Data is not []string
	m := SelectionMenu{
		X:            position.X,
		Y:            position.Y,
		columns:      1,
		focusX:       0,
		focusY:       0,
		SpacingY:     spacingY,
		SpacingX:     spacingX,
		IsShowCursor: true,
		displayStart: 0,
		scale:        1,
		OnSelection:  onSelection,
	}
	m.textBase = text.New(position, BasicAtlas12)
	m.displayRows = 4
	m.cursor = pixel.NewSprite(CursorPng, CursorPng.Bounds())
	m.cursorWidth = CursorPng.Bounds().W()
	m.cursorHeight = CursorPng.Bounds().H()

	if renderFunc != nil {
		m.RenderFunction = renderFunc
	} else {
		m.RenderFunction = m.renderItem
	}

	dataI := reflect.ValueOf(data)
	m.maxRows = dataI.Len() - 1

	if dataI.Len() > 0 {
		m.DataI = make([]interface{}, dataI.Len())
		for i := 0; i < dataI.Len(); i++ {
			m.DataI[i] = dataI.Index(i).Interface()
		}
	}

	//temp implement correct columns pending
	if showColumns {
		m.columns += m.maxRows / m.displayRows
		if m.maxRows == 1 {
			m.columns = 2
			m.displayRows = 1
		}
	}

	m.width = m.calcTotalWidth(xWidth)
	m.height = m.calcTotalHeight()
	return m
}

func (m *SelectionMenu) SetPosition(x, y float64) {
	m.X = x
	m.Y = y
}
func (m *SelectionMenu) OffsetCursorPosition(x, y float64) {
	m.useCursorPos = true
	m.cursorPosOffset = pixel.V(x, y)
}
func (m *SelectionMenu) ShowCursor() {
	m.IsShowCursor = true
}
func (m *SelectionMenu) HideCursor() {
	m.IsShowCursor = false
}

func (m SelectionMenu) GetWidth() float64 {
	return m.width
}
func (m SelectionMenu) GetHeight() float64 {
	return m.height
}

func (m SelectionMenu) calcTotalHeight() float64 {
	height := float64(m.displayRows) * m.SpacingY
	return height - m.SpacingY/2
}
func (m SelectionMenu) calcTotalWidth(xWidth float64) float64 {
	if m.columns == 1 {
		maxEntryWidth := 0.0
		for _, v := range m.DataI {
			switch x := v.(type) {
			case string:
				width := m.textBase.BoundsOf(x).W()
				maxEntryWidth = math.Max(width, maxEntryWidth)

			default:
				//fmt.Println("SelectionMenu:calcTotalWidth :: type unknown")
				maxEntryWidth = xWidth
			}
		}
		return maxEntryWidth + m.cursorWidth
	}
	return m.SpacingX * float64(m.columns)
}

func (m SelectionMenu) IsDataSourceEmpty() bool {
	return len(m.DataI) == 0 || m.DataI == nil
}

func (m SelectionMenu) renderItem(a ...interface{}) {
	//renderer pixel.Target, x, y float64, item string
	rendererV := reflect.ValueOf(a[0])
	renderer := rendererV.Interface().(pixel.Target)
	xV := reflect.ValueOf(a[1])
	x := xV.Interface().(float64)
	yV := reflect.ValueOf(a[2])
	y := yV.Interface().(float64)
	itemV := reflect.ValueOf(a[3])
	item := itemV.Interface().(string)

	pos := pixel.V(x, y)
	//textBase := text.New(pos, BasicAtlas12)
	textBase := text.New(pos, text.NewAtlas(basicfont.Face7x13, text.ASCII))
	if item == "" {
		fmt.Fprintf(textBase, "--")
	} else {
		fmt.Fprintf(textBase, item)
	}
	textBase.Draw(renderer, pixel.IM.Scaled(pixel.V(0, 0), m.scale))
}

func (m SelectionMenu) Render(renderer *pixelgl.Window) {
	displayStart := m.displayStart
	displayEnd := displayStart + m.displayRows

	cursorWidth := m.cursorWidth * m.scale
	cursorHalfWidth := cursorWidth / 2
	cursorHalfHeight := m.cursorHeight / 2
	spacingX := m.SpacingX * m.scale
	rowHeight := m.SpacingY * m.scale

	var x, y = m.X, m.Y
	var mat = pixel.IM.Scaled(pixel.V(x, y), m.scale)

	//temp single columns not rendering hence
	if m.columns == 1 {
		for i, v := range m.DataI {
			cursorPos := pixel.V(x+cursorHalfWidth, y+(cursorHalfHeight/2))
			if m.useCursorPos {
				cursorPos = pixel.V(x+m.cursorPosOffset.X, y+(cursorHalfHeight/2))
			}
			if i == m.focusY && m.IsShowCursor {
				m.cursor.Draw(renderer, mat.Moved(cursorPos))
			}

			switch d := v.(type) {
			case string:
				m.RenderFunction(renderer, x+cursorWidth, y, d)

			default:
				m.RenderFunction(renderer, x, y, d)
			}
			y = y - rowHeight
		}

		return
	}

	//itemIndex := ((displayStart - 1) * m.columns) + 1
	itemIndex := displayStart * m.columns
	for i := displayStart; i < displayEnd; i++ {
		for j := 0; j < m.columns; j++ {
			cursorPos := pixel.V(x+cursorHalfWidth, y+(cursorHalfHeight/2))
			if m.useCursorPos {
				cursorPos = pixel.V(x+m.cursorPosOffset.X, y+(cursorHalfHeight/2))
			}
			if i == m.focusY && j == m.focusX && m.IsShowCursor {
				m.cursor.Draw(renderer, mat.Moved(cursorPos))
			}
			item := m.DataI[itemIndex]
			m.RenderFunction(renderer, x+cursorWidth, y, item)

			x = x + spacingX
			itemIndex = itemIndex + 1
		}
		y = y - rowHeight
		x = m.X
	}
}

func (m *SelectionMenu) MoveUp() {
	m.focusY = utilz.MaxInt(m.focusY-1, 0)
	if m.focusY < m.displayStart {
		m.MoveDisplayUp()
	}
}

func (m *SelectionMenu) MoveDown() {
	if m.columns == 1 {
		m.focusY = utilz.MinInt(m.focusY+1, m.maxRows)
	} else {
		m.focusY = utilz.MinInt(m.focusY+1, m.displayRows-1)
	}

	if m.focusY >= m.displayStart+m.displayRows {
		m.MoveDisplayDown()
	}
}

func (m *SelectionMenu) MoveLeft() {
	m.focusX = utilz.MaxInt(m.focusX-1, 0)
}

func (m *SelectionMenu) MoveRight() {
	m.focusX = utilz.MinInt(m.focusX+1, m.columns-1)
}

func (m *SelectionMenu) MoveDisplayUp() {
	m.displayStart = m.displayStart - 1
}

func (m *SelectionMenu) MoveDisplayDown() {
	m.displayStart = m.displayStart + 1
}

func (m SelectionMenu) SelectedItem() interface{} {
	//return m.GetIndex()
	return m.DataI[m.GetIndex()]
}

func (m SelectionMenu) GetIndex() int {
	//return m.focusX + ((m.focusY - 1) * m.columns)
	return m.focusX + (m.focusY * m.columns)
}

func (m SelectionMenu) OnClick() {
	index := m.GetIndex()
	m.OnSelection(index, m.DataI[index])
}

func (m *SelectionMenu) HandleInput(window *pixelgl.Window) {
	if window.JustPressed(pixelgl.KeyUp) {
		m.MoveUp()
	} else if window.JustPressed(pixelgl.KeyDown) {
		m.MoveDown()
	} else if m.columns > 1 && window.JustPressed(pixelgl.KeyLeft) {
		m.MoveLeft()
	} else if m.columns > 1 && window.JustPressed(pixelgl.KeyRight) {
		m.MoveRight()
	} else if window.JustPressed(pixelgl.KeySpace) {
		m.OnClick()
	}

}

func (m SelectionMenu) CanScrollUp() bool {
	return m.displayStart > 0
}

func (m SelectionMenu) CanScrollDown() bool {
	return m.displayStart <= (m.maxRows - m.displayRows)
}
