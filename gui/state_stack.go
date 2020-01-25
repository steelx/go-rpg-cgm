package gui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/animation"
)

type StackInterface interface {
	Enter()
	Exit()
	Update(dt float64) bool
	Render(win *pixelgl.Window)
	HandleInput(win *pixelgl.Window)
}

//A state stack is, simply, a stack of states. Every time you add a new state it goes on
//top of the stack and is rendered last.
//aka Last in First out
type StateStack struct {
	States  []StackInterface
	Win     *pixelgl.Window
	Globals map[string]interface{}
}

func StateStackCreate(win *pixelgl.Window) *StateStack {
	//call PushFixed, PushFitted after
	return &StateStack{
		Win:     win,
		Globals: make(map[string]interface{}),
	}
}

func (ss *StateStack) Push(state StackInterface) {
	ss.States = append(ss.States, state)
	state.Enter()
}
func (ss *StateStack) Pop() *StackInterface {
	top := ss.States[ss.GetLastIndex()]
	ss.States = ss.States[:len(ss.States)-1] //remove
	top.Exit()
	return &top
}

func (ss *StateStack) Update(dt float64) {
	if len(ss.States) == 0 {
		return
	}
	//The most important state is at the top and it needs updating first.
	//Each state can return a value, stored in the OK variable. If OK is false
	//then the loop breaks and no subsequent states are updated.
	ss.States[ss.GetLastIndex()].Update(dt)

	//this duplicate is needed, after user interaction,
	//user does Pop() hence empty check
	if len(ss.States) == 0 {
		return
	}
	top := ss.States[ss.GetLastIndex()]

	top.HandleInput(ss.Win)
}

func (ss StateStack) Top() *StackInterface {
	return &ss.States[ss.GetLastIndex()]
}

func (ss StateStack) GetLastIndex() int {
	if len(ss.States) == 1 {
		return 0
	}

	return len(ss.States) - 1
}

//Render only last item in array
//unless the last item gets Pop() next would show.
func (ss StateStack) Render(renderer *pixelgl.Window) {
	if len(ss.States) == 0 {
		return
	}
	//ss.States[ss.GetLastIndex()].Render(renderer) //<-- this would render only 1 stack at a time

	//But we want all of them render together
	for _, v := range ss.States {
		v.Render(renderer)
	}
}

func (ss *StateStack) PushSelectionMenu(x, y, width, height float64, txt string, choices []string, onSelection func(int, interface{}), showColumns bool) {
	textBoxMenu := TextboxWithMenuCreate(ss, txt, pixel.V(x, y), width, height, choices, onSelection, showColumns)
	textBoxMenu.AppearTween = animation.TweenCreate(0.9, 1, 0.2)
	ss.States = append(ss.States, textBoxMenu)
}

func (ss *StateStack) PushFixed(x, y, width, height float64, txt, avatarName string, avatarPng pixel.Picture) {
	fixed := TextboxCreateFixed(ss, txt, pixel.V(x, y), width, height, avatarName, avatarPng, false)
	ss.States = append(ss.States, &fixed)
}

//PushFITMenu PENDING not getting correct height and width
func (ss *StateStack) PushFITMenu(x, y float64, txt string, choices []string, onSelection func(int, interface{})) {
	fitMenu := TextboxFITMenuCreate(ss, x, y, txt, choices, onSelection)
	ss.States = append(ss.States, fitMenu)
}

func (ss *StateStack) PushFitted(x, y float64, txt string) *Textbox {
	fitted := TextboxCreateFitted(ss, txt, pixel.V(x, y), false)
	ss.States = append(ss.States, &fitted)
	return &fitted
}
