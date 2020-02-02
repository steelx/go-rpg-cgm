package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/world"
	"reflect"
)

type CombatChoiceState struct {
	Stack       *gui.StateStack
	CombatState *CombatState
	Actor       *combat.Actor
	Character   *Character
	UpArrow,
	DownArrow,
	Marker *pixel.Sprite
	UpArrowPosition,
	DownArrowPosition pixel.Vec
	Selection *gui.SelectionMenu
	textbox   *gui.Textbox
}

func CombatChoiceStateCreate(combatState *CombatState, owner *combat.Actor) *CombatChoiceState {
	c := &CombatChoiceState{
		CombatState: combatState,
		Stack:       combatState.GameState,
		Actor:       owner,
		Character:   combatState.ActorCharMap[owner.Id],
		UpArrow:     world.IconsDB.Get(11),
		DownArrow:   world.IconsDB.Get(12),
		Marker:      pixel.NewSprite(gui.ContinueCaretPng, gui.ContinueCaretPng.Bounds()),
	}

	c.CreateActionDialog(owner.Actions)
	return c
}

func (c *CombatChoiceState) Enter() {
	c.CombatState.SelectedActor = c.Actor
}

func (c *CombatChoiceState) Exit() {
	c.CombatState.SelectedActor = nil
}

func (c *CombatChoiceState) Update(dt float64) bool {
	c.textbox.Update(dt)

	return true
}

func (c CombatChoiceState) Render(renderer *pixelgl.Window) {
	c.textbox.Render(renderer)
}

func (c CombatChoiceState) HandleInput(win *pixelgl.Window) {
	c.Selection.HandleInput(win)
	//if win.JustPressed(pixelgl.KeyEscape) {
	//	c.CombatState.InternalStack.Pop()
	//}
}

func (c *CombatChoiceState) OnSelect(index int, str interface{}) {
	actionItem := reflect.ValueOf(str).Interface().(string)
	if actionItem == "attack" {
		fmt.Println("Character attacks")
		c.Selection.HideCursor()
		//pending
	}
}

func (c *CombatChoiceState) SetArrowPosition() {
	x, y := c.textbox.Position.X, c.textbox.Position.Y
	width, height := c.textbox.Width, c.textbox.Height

	arrowPad := 9.0
	arrowX := x + width - arrowPad
	c.UpArrowPosition = pixel.V(arrowX, y-arrowPad)
	c.DownArrowPosition = pixel.V(arrowX, y-height+arrowPad)
}
func (c *CombatChoiceState) CreateActionDialog(choices interface{}) {
	selectionMenu := gui.SelectionMenuCreate(20, 0, 0,
		choices,
		false,
		pixel.ZV,
		c.OnSelect,
		nil,
	)
	c.Selection = &selectionMenu

	x := -c.Stack.Win.Bounds().W() / 2
	y := -c.Stack.Win.Bounds().H() / 2

	height := c.Selection.GetHeight() + 18
	//width := c.Selection.GetWidth() + 16

	y = y + height + 16
	x = x + 200

	c.textbox = gui.TextboxFITPassedMenuCreate(
		c.Stack,
		x, y, "",
		c.Selection,
	)
}
