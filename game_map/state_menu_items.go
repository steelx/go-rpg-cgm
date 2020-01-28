package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/world"
	"golang.org/x/image/font/basicfont"
	"reflect"
)

type ItemsMenuState struct {
	win            *pixelgl.Window
	parent         *InGameMenuState
	Layout         gui.Layout
	Stack          *gui.StateStack
	StateMachine   *state_machine.StateMachine
	Panels         []gui.Panel
	ItemMenus      []*gui.SelectionMenu
	CategoryMenu   *gui.SelectionMenu
	InCategoryMenu bool
}

func ItemsMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *ItemsMenuState {

	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 118, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.12, 2)
	layout.SplitVert("top", "title", "category", 0.6, 2)
	layout.SplitHorz("bottom", "mid", "inv", 0.14, 2)

	im := &ItemsMenuState{
		win:            win,
		parent:         parent,
		Stack:          parent.Stack,
		StateMachine:   parent.StateMachine,
		Layout:         layout,
		InCategoryMenu: true,
	}

	im.Panels = []gui.Panel{
		layout.CreatePanel("title"),
		layout.CreatePanel("category"),
		layout.CreatePanel("mid"),
		layout.CreatePanel("inv"),
	}

	renderFunction := func(a ...interface{}) {
		//DrawItem
		rendererV := reflect.ValueOf(a[0])
		renderer := rendererV.Interface().(pixel.Target)
		xV := reflect.ValueOf(a[1])
		x := xV.Interface().(float64)
		yV := reflect.ValueOf(a[2])
		y := yV.Interface().(float64)
		itemV := reflect.ValueOf(a[3])
		item := itemV.Interface().(world.ItemIndex)

		parent.World.DrawItem(renderer, x, y, item)
	}

	itemsMenu := gui.SelectionMenuCreate(24, 128,
		parent.World.Items,
		false,
		pixel.V(0, 0),
		func(index int, s interface{}) {
			//Items menu screen selection
		},
		renderFunction,
	)
	keyItemsMenu := gui.SelectionMenuCreate(24, 128,
		parent.World.KeyItems,
		false,
		pixel.V(0, 0),
		func(index int, s interface{}) {
			//Items menu screen selection
		},
		renderFunction,
	)
	im.ItemMenus = []*gui.SelectionMenu{&itemsMenu, &keyItemsMenu}

	categoryMenu := gui.SelectionMenuCreate(24, 128,
		[]string{"Use", "Key Items"},
		true,
		pixel.V(0, 0),
		func(index int, s interface{}) {
			im.OnCategorySelect(index, s)
		},
		nil,
	)
	im.CategoryMenu = &categoryMenu

	//initially since we are InCategoryMenu, we hide ItemMenus selection arrow
	for _, v := range im.ItemMenus {
		v.HideCursor()
	}

	return im
}

func (im *ItemsMenuState) OnCategorySelect(index int, value interface{}) {
	im.CategoryMenu.HideCursor()
	im.InCategoryMenu = false
	menu := im.ItemMenus[index]
	menu.ShowCursor()
}

/*
	state_machine.State implemented below
*/
func (im ItemsMenuState) Enter(data interface{}) {

}

func (im ItemsMenuState) Render(win *pixelgl.Window) {
	for _, v := range im.Panels {
		v.Draw(win)
	}

	titleX := im.Layout.Left("title") + 16
	titleY := im.Layout.MidY("title")

	pos := pixel.V(titleX, titleY)
	textBase := text.New(pos, text.NewAtlas(basicfont.Face7x13, text.ASCII))
	fmt.Fprintln(textBase, "Items")
	textBase.Draw(win, pixel.IM)

	categoryX := im.Layout.Left("category") + 5
	categoryY := im.Layout.MidY("category")
	im.CategoryMenu.SetPosition(categoryX, categoryY)
	im.CategoryMenu.Render(win)

	menu := im.ItemMenus[im.CategoryMenu.GetIndex()]
	if menu.IsDataSourceEmpty() {
		return
	}

	if !im.InCategoryMenu || !im.CategoryMenu.IsShowCursor {
		//convert interface to world.ItemIndex type
		selectedItemIdxV := reflect.ValueOf(menu.SelectedItem())
		selectedItemIdx := selectedItemIdxV.Interface().(world.ItemIndex)
		itemDef := world.ItemsDB[selectedItemIdx.Id]

		//render description
		descX := im.Layout.Left("mid") + 20
		descY := im.Layout.MidY("mid")
		pos = pixel.V(descX, descY)
		textBase = text.New(pos, text.NewAtlas(basicfont.Face7x13, text.ASCII))
		fmt.Fprintln(textBase, itemDef.Description)
		textBase.Draw(win, pixel.IM)
	}

	itemX := im.Layout.Left("inv") - 6
	itemY := im.Layout.Top("inv") - 24
	menu.SetPosition(itemX, itemY)
	menu.Render(win)

}

func (im ItemsMenuState) Exit() {

}

func (im *ItemsMenuState) Update(dt float64) {

	if im.InCategoryMenu && im.CategoryMenu.IsShowCursor {
		if im.win.JustReleased(pixelgl.KeyBackspace) || im.win.JustReleased(pixelgl.KeyEscape) {
			im.StateMachine.Change("frontmenu", nil)
		}
		im.CategoryMenu.HandleInput(im.win)
		return
	}
	menu := im.ItemMenus[im.CategoryMenu.GetIndex()]
	menu.HandleInput(im.win)
	if im.win.JustReleased(pixelgl.KeyBackspace) || im.win.JustReleased(pixelgl.KeyEscape) {
		im.FocusOnCategoryMenu()
	}
}

func (im *ItemsMenuState) FocusOnCategoryMenu() {
	im.InCategoryMenu = true
	menu := im.ItemMenus[im.CategoryMenu.GetIndex()]
	menu.HideCursor()
	im.CategoryMenu.ShowCursor()
}
