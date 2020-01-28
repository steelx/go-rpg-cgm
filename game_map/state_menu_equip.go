package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/world"
)

type EquipMenuState struct {
	parent *InGameMenuState
	win    *pixelgl.Window

	TopBarText,
	PrevTopBarText string
	Selections                    *gui.SelectionMenu
	PartyMenu                     *gui.SelectionMenu
	Panels                        []gui.Panel
	Layout                        gui.Layout
	betterStatsIcon, badStatsIcon *pixel.Sprite
	inList                        bool
}

func EquipMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *EquipMenuState {
	// Create panel layout
	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.12, 2)
	layout.SplitVert("top", "title", "category", 0.75, 2)
	titlePanel := layout.Panels["title"]

	layout = gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.42, 2)
	layout.SplitHorz("bottom", "desc", "bottom", 0.2, 2)
	layout.SplitVert("bottom", "stats", "list", 0.6, 2)
	layout.Panels["title"] = titlePanel

	e := &EquipMenuState{
		win:             win,
		parent:          parent,
		betterStatsIcon: world.IconsDB.Get(10),
		badStatsIcon:    world.IconsDB.Get(11),
		Layout:          layout,
		Panels: []gui.Panel{
			layout.CreatePanel("top"),
			layout.CreatePanel("desc"),
			layout.CreatePanel("stats"),
			layout.CreatePanel("list"),
			layout.CreatePanel("title"),
		},
	}

	return e
}

func (e *EquipMenuState) Enter(data interface{}) {
}

func (e EquipMenuState) Render(win *pixelgl.Window) {
	for _, v := range e.Panels {
		v.Draw(win)
	}
}

func (e EquipMenuState) Exit() {
}

func (e *EquipMenuState) Update(dt float64) {
	if e.inList {

	} else {

		if e.win.JustPressed(pixelgl.KeyEscape) {
			e.parent.StateMachine.Change("frontmenu", nil)
			return
		}
	}
}
