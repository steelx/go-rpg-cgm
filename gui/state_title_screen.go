package gui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/utilz"
	"golang.org/x/image/font/basicfont"
)

type TitleScreen struct {
	Stack       *StateStack
	titleImg    pixel.Picture
	titleSprite *pixel.Sprite
	titlePos    pixel.Vec
	menu        *SelectionMenu
	win         *pixelgl.Window
}

func TitleScreenCreate(stack *StateStack, win *pixelgl.Window) TitleScreen {
	titleImg, err := utilz.LoadPicture("../resources/title_screen.png")
	utilz.PanicIfErr(err)

	s := TitleScreen{
		Stack:    stack,
		titleImg: titleImg,
		win:      win,
	}

	position := win.Bounds().Center()
	s.titleSprite = pixel.NewSprite(s.titleImg, s.titleImg.Bounds())
	s.titlePos = pixel.V(position.X-s.titleImg.Bounds().W()/2, position.Y+s.titleImg.Bounds().H()/2)

	choices := []string{"Play", "Exit"}
	menu_ := SelectionMenuCreate(24, 128, 0,
		choices, false,
		s.titlePos.Add(pixel.V(0, -s.titleImg.Bounds().H()/2-50)),
		s.onSelection, nil)
	s.menu = &menu_
	return s
}

func (s *TitleScreen) onSelection(index int, str interface{}) {
	if index == 0 {
		s.Stack.Pop()
	}
}

/*
	StackInterface implemented below
*/
func (s TitleScreen) Enter() {
}

func (s TitleScreen) Exit() {
}

func (s TitleScreen) Update(dt float64) bool {
	s.menu.HandleInput(s.win)
	return true
}

func (s TitleScreen) Render(win *pixelgl.Window) {
	s.titleSprite.Draw(win, pixel.IM.Moved(s.titlePos))
	s.menu.Render(win)

	toTheScreen := pixel.V(s.titlePos.X, s.titlePos.Y-s.titleImg.Bounds().H()/2+20)

	txtPos := s.titlePos.Add(pixel.V(0, -s.titleImg.Bounds().H()))
	txt := `Press "Space" to select`
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textBase := text.New(txtPos, basicAtlas)
	fmt.Fprintln(textBase, txt)
	textBase.Draw(win, pixel.IM.Moved(txtPos))

	//Camera
	camera := pixel.IM.Scaled(toTheScreen, 1.0).Moved(win.Bounds().Center().Sub(toTheScreen))
	win.SetMatrix(camera)
}

func (s TitleScreen) HandleInput(win *pixelgl.Window) {
}
