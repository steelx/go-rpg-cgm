package game_map

import "github.com/steelx/go-rpg-cgm/combat"

//CombatEvent (CE)
type CombatEvent interface {
	Name() string
	CountDown() float64
	CountDownSet(t float64)
	Owner() *combat.Actor
	Update()
	IsFinished() bool
	Execute(queue *EventQueue)
	TimePoints(queue *EventQueue) float64
}
