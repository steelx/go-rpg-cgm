package main

type Trigger struct {
	OnEnter func()
	OnExit  func()
	OnUse   func()
}
