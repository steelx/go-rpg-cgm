package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/globals"
	"math"
)

/* e.g.
menu2 := gui.SelectionMenuCreate([]string{"Menu 1", "", "Menu 2", "Menu 03", "Menu 04", "Menu 05"}, false, pixel.V(200, 250), func(i int, item string) {
		fmt.Println(i, item)
	})
*/

type SelectionMenu struct {
	X, Y           float64
	width, height  float64
	dataSource     []string //The list of items to be displayed. It can’t be empty
	columns        int      //The number of columns the menu has. This defaults to 1
	focusX, focusY int      //Indicates which item in the list is currently selected.
	//focusX tells us which column is selected and focusY which element in that column
	spacingY, spacingX        float64 //space btw each items
	scale                     float64 //menu scale in size
	cursor                    *pixel.Sprite
	cursorWidth, cursorHeight float64
	showCursor                bool
	maxRows, displayRows      int //rows might be 30 but only 5 maxRows are displayed at once
	displayStart              int //index at which we start displaying menu, e.g. out of 30 max 5 are visible from index 6
	renderer                  pixel.Target
	textBase                  *text.Text
	OnSelection               func(int, string) //to be called after selection
}

func SelectionMenuCreate(data []string, showColumns bool, position pixel.Vec, onSelection func(int, string)) SelectionMenu {
	m := SelectionMenu{
		X:            position.X,
		Y:            position.Y,
		dataSource:   data,
		columns:      1,
		focusX:       0,
		focusY:       0,
		spacingY:     24,
		spacingX:     128,
		showCursor:   true,
		maxRows:      len(data) - 1,
		displayStart: 0,
		scale:        1,
		OnSelection:  onSelection,
	}
	m.textBase = text.New(position, globals.BasicAtlas12)
	m.renderer = globals.Global.Win
	m.displayRows = 4
	m.cursor = pixel.NewSprite(globals.CursorPng, globals.CursorPng.Bounds())
	m.cursorWidth = globals.CursorPng.Bounds().W()
	m.cursorHeight = globals.CursorPng.Bounds().H()

	if showColumns {
		m.columns += m.maxRows / m.displayRows
	}

	m.width = m.calcTotalWidth()
	m.height = m.calcTotalHeight()
	return m
}

func (m *SelectionMenu) SetPosition(x, y float64) {
	m.X = x
	m.Y = y
}

func (m SelectionMenu) calcTotalHeight() float64 {
	height := float64(m.displayRows) * m.spacingY
	return height - m.spacingY/2
}
func (m SelectionMenu) calcTotalWidth() float64 {
	if m.columns == 1 {
		maxEntryWidth := 0.0
		for _, txt := range m.dataSource {
			width := m.textBase.BoundsOf(txt).W()
			maxEntryWidth = math.Max(width, maxEntryWidth)
		}
		return maxEntryWidth + m.cursorWidth
	}
	return m.spacingX * float64(m.columns)
}

func (m SelectionMenu) renderItem(pos pixel.Vec, item string, renderer pixel.Target) {
	textBase := text.New(pos, globals.BasicAtlas12)
	if item == "" {
		fmt.Fprintf(textBase, "--")
	} else {
		fmt.Fprintf(textBase, item)
	}
	textBase.Draw(m.renderer, pixel.IM.Scaled(pixel.V(0, 0), m.scale))
}

func (m SelectionMenu) Render(renderer *pixelgl.Window) {
	displayStart := m.displayStart
	displayEnd := displayStart + m.displayRows

	cursorWidth := m.cursorWidth * m.scale
	cursorHalfWidth := cursorWidth / 2
	cursorHalfHeight := m.cursorHeight / 2
	spacingX := m.spacingX * m.scale
	rowHeight := m.spacingY * m.scale

	var x, y = m.X, m.Y
	var mat = pixel.IM.Scaled(pixel.V(x, y), m.scale)

	if m.columns == 1 {
		for i := 0; i < len(m.dataSource); i++ {
			if i == 0 && m.showCursor {
				m.cursor.Draw(renderer, mat.Moved(pixel.V(x+cursorHalfWidth, y+cursorHalfHeight/2)))
			}
			m.renderItem(pixel.V(x+cursorWidth, y), m.dataSource[i], renderer)
			x = x + spacingX
		}
		return
	}

	//itemIndex := ((displayStart - 1) * m.columns) + 1
	itemIndex := displayStart * m.columns
	for i := displayStart; i < displayEnd; i++ {
		for j := 0; j < m.columns; j++ {
			if i == m.focusY && j == m.focusX && m.showCursor {
				m.cursor.Draw(renderer, mat.Moved(pixel.V(x+cursorHalfWidth, y+cursorHalfHeight/2)))
			}
			item := m.dataSource[itemIndex]
			m.renderItem(pixel.V(x+cursorWidth, y), item, renderer)

			x = x + spacingX
			itemIndex = itemIndex + 1
		}
		y = y - rowHeight
		x = m.X
	}
}

func (m *SelectionMenu) MoveUp() {
	m.focusY = globals.MaxInt(m.focusY-1, 0)
	if m.focusY < m.displayStart {
		m.MoveDisplayUp()
	}
}

func (m *SelectionMenu) MoveDown() {
	if m.columns == 1 {
		m.focusY = globals.MinInt(m.focusY+1, m.maxRows)
	} else {
		m.focusY = globals.MinInt(m.focusY+1, m.displayRows-1)
	}

	if m.focusY >= m.displayStart+m.displayRows {
		m.MoveDisplayDown()
	}
}

func (m *SelectionMenu) MoveLeft() {
	m.focusX = globals.MaxInt(m.focusX-1, 0)
}

func (m *SelectionMenu) MoveRight() {
	m.focusX = globals.MinInt(m.focusX+1, m.columns-1)
}

func (m *SelectionMenu) MoveDisplayUp() {
	m.displayStart = m.displayStart - 1
}

func (m *SelectionMenu) MoveDisplayDown() {
	m.displayStart = m.displayStart + 1
}

func (m SelectionMenu) GetIndex() int {
	//return m.focusX + ((m.focusY - 1) * m.columns)
	return m.focusX + (m.focusY * m.columns)
}

func (m SelectionMenu) OnClick() {
	index := m.GetIndex()
	m.OnSelection(index, m.dataSource[index])
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
