package game_map

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
)

type EquipMenuState struct {
	parent *InGameMenuState
	win    *pixelgl.Window

	TopBarText,
	PrevTopBarText string
	Selections *gui.SelectionMenu
	PartyMenu  *gui.SelectionMenu
	Panels     []gui.Panel
	Layout     gui.Layout
}

func EquipMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *EquipMenuState {
	e := &EquipMenuState{
		win:    win,
		parent: parent,
	}

	return e
}

func (e *EquipMenuState) Enter(data interface{}) {
}

func (e EquipMenuState) Render(win *pixelgl.Window) {
}

func (e EquipMenuState) Exit() {
}

func (e *EquipMenuState) Update(dt float64) {
}
