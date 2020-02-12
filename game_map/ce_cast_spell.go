package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/world"
	"math"
	"reflect"
)

type CECastSpell struct {
	name        string
	countDown   float64
	owner       *combat.Actor
	Targets     []*combat.Actor
	Scene       *CombatState
	mIsFinished bool
	Character   *Character
	Storyboard  *Storyboard
	Spell       world.SpecialItem
}

func CECastSpellCreate(scene *CombatState, owner *combat.Actor, targets []*combat.Actor, spellI interface{}) CombatEvent {
	spell := reflect.ValueOf(spellI).Interface().(world.SpecialItem)
	c := &CECastSpell{
		name:      fmt.Sprintf("%s is casting spell: %s", owner.Name, spell.Name),
		owner:     owner,
		Targets:   targets,
		Scene:     scene,
		Character: scene.ActorCharMap[owner],
		Spell:     spell,
	}

	c.Character.Controller.Change(csRunanim, csProne, true)
	storyboardEvents := []interface{}{
		RunFunction(c.ShowSpellNotice),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: 3}),
		Wait(0.5),
		RunState(c.Character.Controller, csRunanim, csCast, false),
		Wait(0.20),
		RunState(c.Character.Controller, csRunanim, csProne, false),
		RunFunction(c.DoCast),
		Wait(1),
		RunFunction(c.HideSpellNotice),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: -3}),
		RunFunction(c.DoFinish),
	}

	c.Storyboard = StoryboardCreate(scene.InternalStack, scene.win, storyboardEvents, false)

	return c
}

func (c *CECastSpell) Name() string {
	return c.name
}

func (c *CECastSpell) CountDown() float64 {
	return c.countDown
}

func (c *CECastSpell) CountDownSet(t float64) {
	c.countDown = t
}

func (c *CECastSpell) Owner() *combat.Actor {
	return c.owner
}

func (c *CECastSpell) Update() {
	c.Scene.InternalStack.Push(c.Storyboard)
	for i := len(c.Targets) - 1; i >= 0; i-- {
		v := c.Targets[i]
		hp := v.Stats.Get("HpNow")
		isEnemy := false
		if !c.Scene.IsPartyMember(v) {
			isEnemy = true
		}
		if isEnemy && hp <= 0 {
			c.Targets = removeActorAtIndex(c.Targets, i)
		}
	}

	if len(c.Targets) == 0 {
		selectorF := CombatSelectorMap[c.Spell.Target.Selector]
		c.Targets = selectorF(c.Scene)
	}
}

func (c *CECastSpell) IsFinished() bool {
	return c.mIsFinished
}

func (c *CECastSpell) DoFinish() {
	c.mIsFinished = true
}

func (c *CECastSpell) Execute(queue *EventQueue) {

}

func (c CECastSpell) TimePoints(queue *EventQueue) float64 {
	speed := c.owner.Stats.Get("Speed")
	tp := queue.SpeedToTimePoints(speed)
	return tp + c.Spell.TimePoints
}

func (c *CECastSpell) ShowSpellNotice() {
	c.Scene.ShowNotice(c.Spell.Name)
}
func (c *CECastSpell) HideSpellNotice() {
	c.Scene.HideNotice()
}

func (c *CECastSpell) DoCast() {
	pos := c.Character.Entity.GetSelectPosition()
	fxEntity := Entities["fx_use_item"]
	effect := AnimEntityFxCreate(pos.X, pos.Y, fxEntity, fxEntity.Frames, 0.1)
	c.Scene.AddEffect(effect)

	mpNow := c.owner.Stats.Get("MpNow")
	cost := c.Spell.MpCost
	mp := math.Max(mpNow-cost, 0)

	c.owner.Stats.Set("MpNow", mp)

	action := c.Spell.Action
	CombatActions[action](c.Scene, c.owner, c.Targets, c.Spell)
}
