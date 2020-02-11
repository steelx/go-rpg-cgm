package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
	"reflect"
)

type BrowseListState struct {
	Stack         *gui.StateStack
	X, Y          float64
	Width, Height float64
	Title         string
	OnFocus       func(item interface{})
	OnExit        func()
	Selection     *gui.SelectionMenu
	Box           gui.Panel

	UpArrow, DownArrow                 *pixel.Sprite
	UpArrowPosition, DownArrowPosition pixel.Vec
	hide                               bool
}

func BrowseListStateCreate(
	stack *gui.StateStack, x, y, width, height float64, title string, onFocus func(item interface{}), onExit func(), data interface{}, args ...interface{}) *BrowseListState {
	//args sequence -> onSelection func(int, interface{}), renderFunc(a ...interface{})

	s := &BrowseListState{
		Stack:     stack,
		X:         x,
		Y:         y,
		Width:     width,
		Height:    height,
		Title:     title,
		OnExit:    onExit,
		OnFocus:   onFocus,
		UpArrow:   world.IconsDB.Get(12),
		DownArrow: world.IconsDB.Get(13),
		hide:      false,
	}

	var selectCallback func(*BrowseListState, int, interface{}) = nil
	if len(args) >= 1 {
		selectCallback = reflect.ValueOf(args[0]).Interface().(func(*BrowseListState, int, interface{}))
	}

	var renderFunc func(a ...interface{}) = nil
	if len(args) >= 2 {
		renderFunc = reflect.ValueOf(args[1]).Interface().(func(a ...interface{}))
	}

	menu := gui.SelectionMenuCreate(19, 132, width,
		data,
		false,
		pixel.V(x, y),
		func(index int, itemIdx interface{}) {
			selectCallback(s, index, itemIdx)
		},
		renderFunc,
	)

	s.Selection = &menu
	s.Selection.SetPosition(x-80, y+20)
	s.Box = gui.PanelCreate(pixel.V(x, y), s.Selection.GetWidth(), height)

	s.SetArrowPosition()

	return s
}

func (s *BrowseListState) Enter() {
	s.OnFocus(s.Selection.SelectedItem())
}

func (s *BrowseListState) Exit() {
	s.OnExit()
}

func (s *BrowseListState) Update(dt float64) bool {
	return false
}

func (s *BrowseListState) Render(renderer *pixelgl.Window) {
	if s.hide {
		return
	}
	s.Box.Draw(renderer)
	if s.Selection.CanScrollUp() {
		s.UpArrow.Draw(renderer, pixel.IM.Moved(s.UpArrowPosition))
	}
	if s.Selection.CanScrollDown() {
		s.DownArrow.Draw(renderer, pixel.IM.Moved(s.DownArrowPosition))
	}

	pos := pixel.V(s.X, s.Y)
	shadow := utilz.HexToColor("#000000")
	textBase := text.New(pos, gui.BasicAtlasAscii)
	textBase.Color = shadow
	fmt.Fprintln(textBase, s.Title)
	textBase.Draw(renderer, pixel.IM.Moved(pos))

	pos = pixel.V(s.X+0.5, s.Y-0.5)
	textBase = text.New(pos, gui.BasicAtlas12)
	fmt.Fprintln(textBase, s.Title)
	textBase.Draw(renderer, pixel.IM.Moved(pos))
	s.Selection.Render(renderer)
}

func (s *BrowseListState) HandleInput(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEscape) {
		s.Stack.Pop()
		return
	}

	s.Selection.HandleInput(win)
	s.OnFocus(s.Selection.SelectedItem())
}

func (s *BrowseListState) Hide() {
	s.hide = true
}

func (s *BrowseListState) Show() {
	s.hide = false
	s.OnFocus(s.Selection.SelectedItem())
}

func (s *BrowseListState) SetArrowPosition() {
	_, topRight, _, bottomRight := s.Box.GetCorners()
	s.UpArrowPosition = topRight.Sub(pixel.V(7, 7))
	s.DownArrowPosition = bottomRight.Sub(pixel.V(7, -7))
}
