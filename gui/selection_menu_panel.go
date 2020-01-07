package gui

import (
	"github.com/faiface/pixel"
)

type SelectionMenuPanel struct {
	textbox *Textbox
	menu    *SelectionMenu
}

func SelectionMenuPanelCreate(
	textBoxText string,
	panelPos pixel.Vec, panelWidth, panelHeight float64,
	choices []string, onSelection func(int, string)) SelectionMenuPanel {

	textbox := TextboxCreate(
		textBoxText,
		panelPos, panelWidth, panelHeight,
		"",
		nil,
		true,
	)

	textBounds := textbox.getTextBound()

	menu := SelectionMenuCreate(choices, true, pixel.V(
		textbox.Position.X-10, textbox.Position.Y-textBounds.H()-10), onSelection)

	return SelectionMenuPanel{
		textbox: &textbox,
		menu:    &menu,
	}
}

func (sm SelectionMenuPanel) Render(renderer pixel.Target) {
	sm.textbox.RenderWithPanel(renderer)
	sm.menu.Render()
	sm.menu.HandleInput()
}
