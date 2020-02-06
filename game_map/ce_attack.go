package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
)

type CEAttack struct {
	name            string
	countDown       float64
	owner           *combat.Actor
	Targets         []*combat.Actor
	Scene           *CombatState
	Finished        bool
	Character       *Character
	Storyboard      *Storyboard
	AttackEntityDef EntityDefinition
	DefaultTargeter func(state *CombatState) []*combat.Actor
	options         AttackOptions
}

type AttackOptions struct {
	Counter bool
}

func CEAttackCreate(scene *CombatState, owner *combat.Actor, targets []*combat.Actor, options AttackOptions) *CEAttack {
	c := &CEAttack{
		options:   options,
		Scene:     scene,
		owner:     owner,
		Targets:   targets,
		Character: scene.ActorCharMap[owner],
		name:      fmt.Sprintf("Attack for %s ->)", owner.Name),
	}
	c.Character.Controller.Change(csRunanim, csProne, true) //CombatState, CombatAnimationID
	var storyboardEvents []interface{}
	if owner.IsPlayer() {
		c.DefaultTargeter = CombatSelector.WeakestEnemy
		c.AttackEntityDef = Entities["slash"]

		storyboardEvents = []interface{}{
			//stateMachine, stateID, ...animID, additionalParams
			RunState(c.Character.Controller, csMove, CSMoveParams{Dir: 1}),
			RunState(c.Character.Controller, csRunanim, csAttack, false),
			RunFunction(c.DoAttack),
			RunState(c.Character.Controller, csMove, CSMoveParams{Dir: -1}),
			RunFunction(c.onFinished),
			RunState(c.Character.Controller, csRunanim, csStandby, false),
		}
	} else {
		c.DefaultTargeter = CombatSelector.RandomAlivePlayer
		c.AttackEntityDef = Entities["claw"]

		storyboardEvents = []interface{}{
			RunState(c.Character.Controller, csMove, CSMoveParams{Dir: -1, Distance: 10, Time: 0.2}),
			RunFunction(c.DoAttack),
			RunState(c.Character.Controller, csMove, CSMoveParams{Dir: 1, Distance: 10, Time: 0.4}),
			RunFunction(c.onFinished),
			RunState(c.Character.Controller, csRunanim, csStandby, false),
		}
	}

	c.Storyboard = StoryboardCreate(scene.InternalStack, scene.win, storyboardEvents, false)

	return c
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
	return c.Finished
}

func (c *CEAttack) Execute(queue *EventQueue) {
	c.Scene.InternalStack.Push(c.Storyboard)

	for i := len(c.Targets) - 1; i >= 0; i-- {
		v := c.Targets[i]
		hpNow := v.Stats.Get("HpNow")
		if hpNow <= 0 {
			c.Targets = c.removeAtIndex(c.Targets, i)
		}
	}

	//find next Target!
	if len(c.Targets) == 0 {
		c.Targets = c.DefaultTargeter(c.Scene)
	}
}

func (c CEAttack) removeAtIndex(arr []*combat.Actor, i int) []*combat.Actor {
	return append(arr[:i], arr[i+1:]...)
}

func (c CEAttack) TimePoints(queue EventQueue) float64 {
	speed := c.Owner().Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}

func (c *CEAttack) onFinished() {
	c.Finished = true
}

//CounterTarget - Decide if the attack is countered.
func (c *CEAttack) CounterTarget(target *combat.Actor) {
	countered := Formula.IsCountered(c.Scene, c.owner, target)
	if countered {
		c.Scene.ApplyCounter(target, c.owner)
	}
}

func (c *CEAttack) DoAttack() {
	for _, v := range c.Targets {
		c.attackTarget(v)

		if !c.options.Counter {
			c.CounterTarget(v)
		}
	}
}
func (c *CEAttack) attackTarget(target *combat.Actor) {

	//hit result lets us know the status of this attack
	damage, hitResult := Formula.MeleeAttack(c.Scene, c.owner, target)
	entity := c.Scene.ActorCharMap[target].Entity

	if hitResult == HitResultMiss {
		c.Scene.ApplyMiss(target)
		return
	} else if hitResult == HitResultDodge {
		c.Scene.ApplyDodge(target)
	}

	var isCrit bool
	if hitResult == HitResultCritical {
		isCrit = true
	}
	c.Scene.ApplyDamage(target, damage, isCrit)

	//FX
	x, y := entity.X, entity.Y
	effect := AnimEntityFxCreate(x, y, c.AttackEntityDef, c.AttackEntityDef.Frames)
	c.Scene.AddEffect(effect)
}
