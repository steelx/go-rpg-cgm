package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
)

//CombatEventTurn
type CETurn struct {
	Scene     *CombatState
	owner     *combat.Actor
	name      string
	countDown float64
	finished  bool
}

func CETurnCreate(scene *CombatState, owner *combat.Actor) *CETurn {
	return &CETurn{
		Scene: scene,
		owner: owner,
		name:  fmt.Sprintf("Turn for (%s)", owner.Name),
	}
}

func (c CETurn) Name() string {
	return c.name
}

func (c CETurn) CountDown() float64 {
	return c.countDown
}

func (c *CETurn) CountDownSet(t float64) {
	c.countDown = t
}

func (c CETurn) Owner() *combat.Actor {
	return c.owner
}

func (c *CETurn) Update() {

}

func (c CETurn) IsFinished() bool {
	return c.finished
}

func (c *CETurn) Execute(queue *EventQueue) {

	// 1. Player
	if c.Scene.IsPartyMember(c.owner) {
		state := CombatChoiceStateCreate(c.Scene, c.owner)
		c.Scene.InternalStack.Push(state)
	} else {
		// 2. an Enemy
		// do a dumb attack
		targets := CombatSelector.RandomAlivePlayer(c.Scene)
		queue := c.Scene.EventQueue
		event := CEAttackCreate(c.Scene, c.owner, targets)
		tp := event.TimePoints(*queue)
		queue.Add(event, tp)
	}

	c.finished = true

}

func (c CETurn) TimePoints(queue *EventQueue) float64 {
	speed := c.Owner().Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}
