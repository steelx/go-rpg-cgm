package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
	"math"
	"reflect"
)

type CombatChoiceState struct {
	Stack       *gui.StateStack //The internal stack of states from the CombatState object.
	CombatState *CombatState
	Actor       *combat.Actor
	Character   *Character
	UpArrow,
	DownArrow,
	Marker *pixel.Sprite
	UpArrowPosition,
	DownArrowPosition,
	MarkerPosition pixel.Vec
	time      float64
	Selection *gui.SelectionMenu
	textbox   *gui.Textbox
}

func CombatChoiceStateCreate(combatState *CombatState, owner *combat.Actor) *CombatChoiceState {
	c := &CombatChoiceState{
		CombatState: combatState,
		Stack:       combatState.InternalStack,
		Actor:       owner,
		Character:   combatState.ActorCharMap[owner],
		UpArrow:     world.IconsDB.Get(11),
		DownArrow:   world.IconsDB.Get(12),
		Marker:      pixel.NewSprite(gui.ContinueCaretPng, gui.ContinueCaretPng.Bounds()),
	}
	c.MarkerPosition = c.Character.Entity.GetSelectPosition()
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
	c.bounceMarker(dt)
	return true
}

func (c CombatChoiceState) Render(renderer *pixelgl.Window) {
	c.textbox.Render(renderer)

	c.Marker.Draw(renderer, pixel.IM.Moved(c.MarkerPosition))
}

func (c CombatChoiceState) HandleInput(win *pixelgl.Window) {
	c.Selection.HandleInput(win)
}

func (c *CombatChoiceState) OnSelect(index int, str interface{}) {
	actionItem := reflect.ValueOf(str).Interface().(string)
	if actionItem == combat.ActionAttack {
		c.Selection.HideCursor()

		state := CombatTargetStateCreate(c.CombatState, CombatChoiceParams{
			OnSelect: func(targets []*combat.Actor) {
				c.TakeAction(actionItem, targets)
			},
			OnExit: func() {
				c.Selection.ShowCursor()
			},
			SwitchSides:     true,
			DefaultSelector: nil,
			TargetType:      CombatTargetTypeONE,
		})
		c.Stack.Push(state)
	}

	if actionItem == combat.ActionFlee {
		c.Stack.Pop() // choice state
		queue := c.CombatState.EventQueue
		event := CEFleeCreate(c.CombatState, c.Actor, CSMoveParams{Dir: 8, Distance: 180, Time: 0.6})
		tp := event.TimePoints(queue)
		queue.Add(event, tp)
	}
}

//TakeAction function pops the CombatTargetState and CombatChoiceState off the
//stack. This leaves the CombatState internal stack empty and causes the EventQueue
//to start updating again.
func (c *CombatChoiceState) TakeAction(id string, targets []*combat.Actor) {
	c.Stack.Pop() //select state
	c.Stack.Pop() //action state

	if id == combat.ActionAttack {
		logrus.Info("Entered TakeAction 'attack'")
		attack := CEAttackCreate(c.CombatState, c.Actor, targets, AttackOptions{})
		tp := attack.TimePoints(*c.CombatState.EventQueue)
		c.CombatState.EventQueue.Add(attack, tp)
		return
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

	x := c.Stack.Win.Bounds().W() / 2
	y := c.Stack.Win.Bounds().H() / 2

	height := c.Selection.GetHeight() + 18
	//width := c.Selection.GetWidth() + 16

	y = y - height
	x = x - 90

	c.textbox = gui.TextboxFITPassedMenuCreate(
		c.Stack,
		x, y, "",
		c.Selection,
	)
	c.textbox.Panel.BGColor = utilz.HexToColor("#3c2f2f")
}

func (c *CombatChoiceState) bounceMarker(dt float64) {
	c.time = c.time + dt
	bounce := pixel.V(c.MarkerPosition.X, c.MarkerPosition.Y+math.Sin(c.time*5))
	c.MarkerPosition = bounce
}
