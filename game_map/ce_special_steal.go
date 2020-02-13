package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/world"
	"reflect"
)

type CESteal struct {
	mOwner               *combat.Actor
	mName                string
	mCountDown           float64
	mIsFinished, Success bool

	SpecialItem     world.SpecialItem
	Targets         []*combat.Actor
	Scene           *CombatState
	Character       *Character
	Storyboard      *Storyboard
	AttackEntityDef EntityDefinition
	DefaultTargeter func(state *CombatState) []*combat.Actor
	OriginalPos     pixel.Vec
}

func CEStealCreate(scene *CombatState, owner *combat.Actor, targets []*combat.Actor, specialI interface{}) CombatEvent {
	special := reflect.ValueOf(specialI).Interface().(world.SpecialItem)
	c := &CESteal{
		mOwner: owner,
		mName:  fmt.Sprintf("%s is using Steal: %s", owner.Name, special.Name),

		SpecialItem: special,
		Targets:     targets,
		Scene:       scene,
		Character:   scene.ActorCharMap[owner],
	}

	c.OriginalPos = pixel.V(c.Character.Entity.X, c.Character.Entity.Y)
	c.Character.Controller.Change(csRunanim, csProne, true)
	c.DefaultTargeter = CombatSelector.WeakestEnemy
	c.AttackEntityDef = Entities["slash"]

	storyboardEvents := []interface{}{
		RunFunction(c.ShowNotice),
		Wait(0.5),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: 3}),
		RunState(c.Character.Controller, csRunanim, "steal_1", false),
		RunFunction(c.TeleportOut),
		RunState(c.Character.Controller, csRunanim, "steal_2", false),
		Wait(1.2),
		RunFunction(c.DoSteal),
		RunState(c.Character.Controller, csRunanim, "steal_3", false),
		RunFunction(c.TeleportIn),
		RunState(c.Character.Controller, csRunanim, "steal_4", false),
		RunFunction(c.ShowResult),
		Wait(1.0),
		RunFunction(c.HideNotice),
		Wait(0.5),
		RunState(c.Character.Controller, csRunanim, csStandby, false),
		RunFunction(c.OnFinish),
	}

	c.Storyboard = StoryboardCreate(scene.InternalStack, scene.win, storyboardEvents, false)

	return c
}

func (c *CESteal) Name() string {
	return c.mName
}

func (c *CESteal) CountDown() float64 {
	return c.mCountDown
}

func (c *CESteal) CountDownSet(t float64) {
	c.mCountDown = t
}

func (c *CESteal) Owner() *combat.Actor {
	return c.mOwner
}

func (c *CESteal) Update() {

}

func (c *CESteal) IsFinished() bool {
	return c.mIsFinished
}

func (c *CESteal) OnFinish() {
	c.mIsFinished = true
}

func (c *CESteal) Execute(queue *EventQueue) {
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

func (c *CESteal) TimePoints(queue *EventQueue) float64 {
	speed := c.mOwner.Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}

func (c *CESteal) HideNotice() {
	c.Scene.HideNotice()
}
func (c *CESteal) ShowNotice() {
	c.Scene.ShowNotice(c.SpecialItem.Name)
}

func (c *CESteal) TeleportOut() {
	target := c.Targets[0]
	entity := c.Scene.ActorCharMap[target].Entity
	width := entity.Texture.Bounds().W()
	height := entity.Texture.Bounds().H()
	pos := entity.GetSelectPosition()

	c.Character.Entity.X = pos.X - width/2
	c.Character.Entity.Y = pos.Y - height/2
}

func (c *CESteal) TeleportIn() {
	c.Character.Entity.X = c.OriginalPos.X
	c.Character.Entity.Y = c.OriginalPos.Y
}

func (c *CESteal) ShowResult() {
	animId := "steal_failure"
	if c.Success {
		animId = "steal_success"
	}
	c.Character.Controller.Change(csRunanim, animId, false)
}
func (c *CESteal) DoSteal() {
	target := c.Targets[0]
	c.Scene.HideNotice()

	if target.StealItem == 0 {
		c.Scene.ShowNotice("Nothing to steal")
		return
	}

	c.Success = c.StealFrom(target)

	if c.Success {
		id := target.StealItem
		def := world.ItemsDB[id]

		gWorld := reflect.ValueOf(c.Scene.GameState.Globals["world"]).Interface().(*combat.WorldExtended)
		gWorld.AddItem(id, 1)
		target.StealItem = 0 //remove StealItem from enemy
		notice := fmt.Sprintf("Stolen: %s", def.Name)
		c.Scene.ShowNotice(notice)
	} else {
		c.Scene.ShowNotice("Steal failed !")
	}
}

func (c *CESteal) StealFrom(target *combat.Actor) bool {
	success := Formula.Steal(c.Scene, c.mOwner, target)

	entity := c.Scene.ActorCharMap[target].Entity
	pos := entity.GetSelectPosition()
	x := pos.X
	y := pos.Y
	effect := AnimEntityFxCreate(x, y, c.AttackEntityDef, c.AttackEntityDef.Frames)
	c.Scene.AddEffect(effect)

	return success
}
