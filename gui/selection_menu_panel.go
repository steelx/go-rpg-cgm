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

	menu := SelectionMenuCreate(choices, textbox.Position.Add(pixel.V(5, -textbox.Height-10)), onSelection)

	return SelectionMenuPanel{
		textbox: &textbox,
		menu:    &menu,
	}
}

func (sm SelectionMenuPanel) Render() {
	sm.textbox.RenderWithPanel()
	sm.menu.Render()
	sm.textbox.HandleInput()
	sm.menu.HandleInput()
}
