package game_map

import "github.com/steelx/go-rpg-cgm/combat"

//CombatEvent (CE)
type Event interface {
	Name() string
	CountDown() float64
	CountDownSet(t float64)
	Owner() *combat.Actor
	Update()
	IsFinished() bool
	Execute(queue *EventQueue)
}
