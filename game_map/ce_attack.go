package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"math"
)

type CEAttack struct {
	name       string
	countDown  float64
	owner      *combat.Actor
	Targets    []*combat.Actor
	Scene      *CombatState
	Finished   bool
	Character  *Character
	Storyboard *Storyboard
}

func CEAttackCreate(scene *CombatState, owner *combat.Actor, targets []*combat.Actor) *CEAttack {
	c := &CEAttack{
		Scene:     scene,
		owner:     owner,
		Targets:   targets,
		Character: scene.ActorCharMap[owner],
		name:      fmt.Sprintf("Attack for %s ->)", owner.Name),
	}
	c.Character.Controller.Change(csRunanim, csProne, true) //CombatState, CombatAnimationID

	storyboardEvents := []interface{}{
		//stateMachine, stateID, ...animID, additionalParams
		RunState(c.Character.Controller, csMove, Direction{1, 0}),
		RunState(c.Character.Controller, csRunanim, csAttack, false),
		RunFunction(c.DoAttack),
		RunState(c.Character.Controller, csMove, Direction{-1, 0}),
		RunFunction(c.onFinished),
		RunState(c.Character.Controller, csRunanim, csStandby, false),
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
		c.Targets = CombatSelector.WeakestEnemy(c.Scene)
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

func (c *CEAttack) DoAttack() {
	for _, v := range c.Targets {
		c.attackTarget(v)
	}
}
func (c *CEAttack) attackTarget(target *combat.Actor) {

	stats := c.owner.Stats
	enemyStats := target.Stats

	// Simple attack get
	attack := stats.Get("Attack")
	attack = attack + stats.Get("Strength")
	defense := enemyStats.Get("Defense")

	damage := math.Max(0, attack-defense)
	fmt.Println("Attacked for ", damage, attack, defense)

	hp := enemyStats.Get("HpNow")
	hp = hp - damage

	enemyStats.Set("HpNow", math.Max(0, hp))
	fmt.Println("HpNow :", enemyStats.Get("HpNow"))

	// the enemy needs stats
	// the player needs a weapon

	//Change actor's Character to hurt state
	character := c.Scene.ActorCharMap[target]
	if damage > 0 {
		state := character.Controller.Current

		//check if its NOT csHurt then change it to csHurt
		switch state.(type) {
		case *CSHurt:
			//fmt.Println("already in Hurt state, do nothing")
		default:
			character.Controller.Change(csHurt, state)
		}
	}

	c.Scene.HandleDeath()
}
