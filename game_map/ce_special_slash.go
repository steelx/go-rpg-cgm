package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/world"
	"math"
	"reflect"
)

type CESlash struct {
	mOwner      *combat.Actor
	mName       string
	mCountDown  float64
	mIsFinished bool

	SpecialItem     world.SpecialItem
	Targets         []*combat.Actor
	Scene           *CombatState
	Character       *Character
	Storyboard      *Storyboard
	AttackEntityDef EntityDefinition
	DefaultTargeter func(state *CombatState) []*combat.Actor
}

func CESlashCreate(scene *CombatState, owner *combat.Actor, targets []*combat.Actor, specialI interface{}) CombatEvent {
	special := reflect.ValueOf(specialI).Interface().(world.SpecialItem)
	c := &CESlash{
		mOwner: owner,
		mName:  fmt.Sprintf("%s is using Slash : %s", owner.Name, special.Name),

		SpecialItem: special,
		Targets:     targets,
		Scene:       scene,
		Character:   scene.ActorCharMap[owner],
	}

	c.Character.Controller.Change(csRunanim, csProne, true)
	c.DefaultTargeter = CombatSelector.SideEnemy
	c.AttackEntityDef = Entities["slash"]

	storyboardEvents := []interface{}{
		RunFunction(c.ShowNotice),
		Wait(0.5),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: 3}),
		RunState(c.Character.Controller, csRunanim, csSpecial, false),
		Wait(0.5),
		RunState(c.Character.Controller, csRunanim, csProne, false),
		RunFunction(c.DoAttack),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: -3}),
		Wait(0.5),
		RunState(c.Character.Controller, csRunanim, csProne, false),
		RunFunction(c.OnFinish),
	}

	c.Storyboard = StoryboardCreate(scene.InternalStack, scene.win, storyboardEvents, false)

	return c
}

func (c *CESlash) Name() string {
	return c.mName
}

func (c *CESlash) CountDown() float64 {
	return c.mCountDown
}

func (c *CESlash) CountDownSet(t float64) {
	c.mCountDown = t
}

func (c *CESlash) Owner() *combat.Actor {
	return c.mOwner
}

func (c *CESlash) Update() {

}

func (c *CESlash) IsFinished() bool {
	return c.mIsFinished
}
func (c *CESlash) OnFinish() {
	c.mIsFinished = true
}

func (c *CESlash) Execute(queue *EventQueue) {
	c.Scene.InternalStack.Push(c.Storyboard)
	for i := len(c.Targets) - 1; i >= 0; i-- {
		v := c.Targets[i]
		hp := v.Stats.Get("HpNow")
		if hp <= 0 {
			c.Targets = removeActorAtIndex(c.Targets, i)
		}
	}

	if len(c.Targets) == 0 {
		//Find another enemy
		c.Targets = c.DefaultTargeter(c.Scene)
	}
}

func (c *CESlash) TimePoints(queue *EventQueue) float64 {
	speed := c.mOwner.Stats.Get("Speed")
	tp := queue.SpeedToTimePoints(speed)
	return tp + c.SpecialItem.TimePoints
}

func (c *CESlash) ShowNotice() {
	c.Scene.ShowNotice(c.SpecialItem.Name)
}

func (c *CESlash) DoAttack() {
	c.Scene.HideNotice()

	mp := c.mOwner.Stats.Get("MpNow")
	cost := c.SpecialItem.MpCost
	mp = math.Max(mp-cost, 0)
	c.mOwner.Stats.Set("MpNow", mp)
	for _, target := range c.Targets {
		c.AttackTarget(target)
		if !c.SpecialItem.Counter {
			//Decide if the attack is countered.
			c.CounterTarget(target)
		}
	}
}

func (c *CESlash) CounterTarget(target *combat.Actor) {
	countered := Formula.IsCountered(c.Scene, c.mOwner, target)
	if countered {
		c.Scene.ApplyCounter(target, c.mOwner)
	}
}

func (c *CESlash) AttackTarget(target *combat.Actor) {
	damage, hitResult := Formula.MeleeAttack(c.Scene, c.mOwner, target)
	entity := c.Scene.ActorCharMap[target].Entity

	if hitResult == HitResultMiss {
		c.Scene.ApplyMiss(target)
		return
	} else if hitResult == HitResultDodge {
		c.Scene.ApplyDodge(target)
	} else {
		c.Scene.ApplyDamage(target, damage, hitResult == HitResultCritical)
	}

	pos := entity.GetSelectPosition()
	x := pos.X
	y := pos.Y
	slashEffect := AnimEntityFxCreate(x, y, c.AttackEntityDef, c.AttackEntityDef.Frames)
	c.Scene.AddEffect(slashEffect)
}
