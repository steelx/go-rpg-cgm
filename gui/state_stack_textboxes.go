package gui

import (
	"github.com/faiface/pixel"
)

type StateStack struct {
	States []*Textbox
}

func StateStackCreate() StateStack {
	//call AddFixed, AddFitted after
	return StateStack{}
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

func (ss *StateStack) AddFixed(x, y, width, height float64, txt, avatarName string, avatarPng pixel.Picture) {
	fixed := TextboxCreateFixed(txt, pixel.V(x, y), width, height, "Ajinkya", avatarPng,
		false,
	)
	ss.States = append(ss.States, &fixed)
}

func (ss *StateStack) AddFitted(x, y float64, txt string) {
	fitted := TextboxCreateFitted(txt, pixel.V(x, y), false)
	ss.States = append(ss.States, &fitted)
}
