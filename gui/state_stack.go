package gui

import (
	"github.com/faiface/pixel"
)

//A state stack is, simply, a stack of states. Every time you add a new state it goes on
//top of the stack and is rendered last.
type StateStack struct {
	States []*Textbox
}

func StateStackCreate() *StateStack {
	//call AddFixed, AddFitted after
	return &StateStack{}
}

func (ss *StateStack) Push(state *Textbox) {
	ss.States = append(ss.States, state)
	state.Enter()
}
func (ss *StateStack) Pop() *Textbox {
	top := ss.States[ss.getLastIndex()]
	ss.States = ss.States[:len(ss.States)-1] //remove
	top.Exit()
	return top
}

func (ss *StateStack) Update(dt float64) {
	if len(ss.States) == 0 {
		return
	}
	//The most important state is at the top and it needs updating first.
	//Each state can return a value, stored in the OK variable. If OK is false
	//then the loop breaks and no subsequent states are updated.
	if ss.getLastIndex() == 0 {
		v := ss.States[0]
		v.Update(dt)
	} else {
		for last := len(ss.States) - 1; last >= 0; last-- {
			v := ss.States[last]
			if OK := v.Update(dt); !OK {
				break
			}
		}
	}

	top := ss.States[ss.getLastIndex()]

	if top.IsDead() {
		ss.States = ss.States[:len(ss.States)-1]
	}

	top.HandleInput()
}

func (ss StateStack) Render(renderer pixel.Target) {
	for _, v := range ss.States {
		v.Render(renderer)
	}
}

func (ss *StateStack) AddSelectionMenu(x, y, width, height float64, txt string, choices []string, onSelection func(int, string)) {
	textBoxMenu := TextboxWithMenuCreate(ss, txt, pixel.V(x, y), width, height, choices, onSelection)
	ss.States = append(ss.States, textBoxMenu)
}

func (ss *StateStack) AddFixed(
	x, y, width, height float64, txt, avatarName string, avatarPng pixel.Picture) {
	fixed := TextboxCreateFixed(ss, txt, pixel.V(x, y), width, height, "Ajinkya", avatarPng, false)
	ss.States = append(ss.States, &fixed)
}

func (ss *StateStack) AddFitted(x, y float64, txt string) {
	fitted := TextboxCreateFitted(ss, txt, pixel.V(x, y), false)
	ss.States = append(ss.States, &fitted)
}

func (ss StateStack) getLastIndex() int {
	if len(ss.States) == 1 {
		return 0
	}

	return len(ss.States) - 1
}
