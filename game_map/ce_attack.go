package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
)

type CEAttack struct {
	name      string
	countDown float64
	owner,
	Target *combat.Actor
	Scene *CombatState
}

func CEAttackCreate(scene *CombatState, owner, target *combat.Actor) *CEAttack {
	return &CEAttack{
		Scene:  scene,
		owner:  owner,
		Target: target,
		name:   fmt.Sprintf("CEAttack(_, %s -> %s)", owner.Name, target.Name),
	}
}

func (c CEAttack) Name() string {
	return c.name
}

func (c CEAttack) CountDown() float64 {
	return c.countDown
}

func (c *CEAttack) CountDownSet(t float64) {
	c.countDown = t
}

func (c CEAttack) Owner() *combat.Actor {
	return c.owner
}

func (c CEAttack) Update() {
}

func (c CEAttack) IsFinished() bool {
	return true
}

func (c *CEAttack) Execute(queue *EventQueue) {
	target := c.Target
	targetHP := target.Stats.Get("HpNow")
	// has Already killed!
	if targetHP <= 0 {
		//Get a new random target
		target = c.Scene.GetTarget(c.owner)
	}

	damage := c.owner.Stats.Get("Attack")
	targetHP = targetHP - damage
	target.Stats.Set("HpNow", targetHP)

	dmgMsg := fmt.Sprintf("%s hit for %v damage", target.Name, damage)
	fmt.Println(dmgMsg)

	if targetHP <= 0 {
		msg := fmt.Sprintf("%s is killed by %s [%v]", target.Name, c.owner.Name, c.owner.Stats.Get("HpNow"))
		fmt.Println(msg)
		c.Scene.OnDead(target)
	}
}

func (c CEAttack) TimePoints(queue EventQueue) float64 {
	speed := c.Owner().Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}
