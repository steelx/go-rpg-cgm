package gui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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
	States []StackInterface
	win    *pixelgl.Window
}

func StateStackCreate(win *pixelgl.Window) *StateStack {
	//call PushFixed, PushFitted after
	return &StateStack{win: win}
}

func (ss *StateStack) Push(state StackInterface) {
	ss.States = append(ss.States, state)
	state.Enter()
}
func (ss *StateStack) Pop() *StackInterface {
	top := ss.States[ss.getLastIndex()]
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
	for last := len(ss.States) - 1; last >= 0; last-- {
		v := ss.States[last]
		if OK := v.Update(dt); !OK {
			break
		}
	}
	//this duplicate is needed, after user interaction,
	//user does Pop() hence empty check
	if len(ss.States) == 0 {
		return
	}
	top := ss.States[ss.getLastIndex()]

	top.HandleInput(ss.win)
}

func (ss StateStack) getLastIndex() int {
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
	ss.States[ss.getLastIndex()].Render(renderer)
	//for _, v := range ss.States {
	//	v.Render(renderer)
	//}
}

func (ss *StateStack) PushSelectionMenu(x, y, width, height float64, txt string, choices []string, onSelection func(int, string)) {
	textBoxMenu := TextboxWithMenuCreate(ss, txt, pixel.V(x, y), width, height, choices, onSelection)
	ss.States = append(ss.States, textBoxMenu)
}

func (ss *StateStack) PushFixed(
	x, y, width, height float64, txt, avatarName string, avatarPng pixel.Picture) {
	fixed := TextboxCreateFixed(ss, txt, pixel.V(x, y), width, height, "Ajinkya", avatarPng, false)
	ss.States = append(ss.States, &fixed)
}

func (ss *StateStack) PushFitted(x, y float64, txt string) {
	fitted := TextboxCreateFitted(ss, txt, pixel.V(x, y), false)
	ss.States = append(ss.States, &fitted)
}
