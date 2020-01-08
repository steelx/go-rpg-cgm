package state_stacks

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/gui"
)

//A state stack is, simply, a stack of states. Every time you add a new state it goes on
//top of the stack and is rendered last.
type StateStack struct {
	States []*gui.Textbox
}

func StateStackCreate() *StateStack {
	//call AddFixed, AddFitted after
	return &StateStack{}
}

func (ss *StateStack) Update(dt float64) {
	if len(ss.States) == 0 {
		return
	}
	// update them and check input
	for _, v := range ss.States {
		v.Update(dt)
	}

	top := ss.States[len(ss.States)-1]

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
	textBoxMenu := gui.TextboxWithMenuCreate(txt, pixel.V(x, y), width, height, choices, onSelection)
	ss.States = append(ss.States, textBoxMenu)
}

func (ss *StateStack) AddFixed(
	x, y, width, height float64, txt, avatarName string, avatarPng pixel.Picture) {
	fixed := gui.TextboxCreateFixed(txt, pixel.V(x, y), width, height, "Ajinkya", avatarPng, false)
	ss.States = append(ss.States, &fixed)
}

func (ss *StateStack) AddFitted(x, y float64, txt string) {
	fitted := gui.TextboxCreateFitted(txt, pixel.V(x, y), false)
	ss.States = append(ss.States, &fitted)
}
