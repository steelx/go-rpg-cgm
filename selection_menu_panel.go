package main

import (
	"github.com/faiface/pixel"
)

type SelectionMenuPanel struct {
	textbox *Textbox
	menu    *SelectionMenu
}

func SelectionMenuPanelCreate(
	textBoxText string,
	panelPos pixel.Vec,
	panelWidth, panelHeight float64,
	data []string, columns int, onSelection func(int, string)) SelectionMenuPanel {

	textbox := TextboxCreate(
		textBoxText,
		basicAtlas12,
		PanelCreate(panelPng, panelPos, panelWidth, panelHeight),
		continueCaretPng,
		"",
		nil,
		true,
	)

	//i think pos should be panel MinX and MinY
	menu := SelectionMenuCreate(data, columns, textbox.Position.Add(pixel.V(5, -textbox.Height-10)), onSelection)

	return SelectionMenuPanel{
		textbox: &textbox,
		menu:    &menu,
	}
}

func (sm SelectionMenuPanel) Render() {
	sm.textbox.DrawTextWithPanel()
	sm.menu.Render()
	sm.textbox.HandleInput()
	sm.menu.HandleInput()
}
