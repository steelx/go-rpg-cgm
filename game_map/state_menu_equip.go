package game_map

import "github.com/faiface/pixel/pixelgl"

type EquipMenuState struct {
	parent *InGameMenuState
	win    *pixelgl.Window
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
