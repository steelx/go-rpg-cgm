package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/world"
	"reflect"
)

type CEUseItem struct {
	Scene       *CombatState
	Character   *Character
	Targets     []*combat.Actor
	ItemDef     world.Item
	owner       *combat.Actor
	name        string
	countDown   float64
	mIsFinished bool
	Storyboard  *Storyboard
}

func CEUseItemCreate(scene *CombatState, owner *combat.Actor, item world.Item, targets []*combat.Actor) *CEUseItem {
	c := &CEUseItem{
		Scene:     scene,
		Character: scene.ActorCharMap[owner],
		owner:     owner,
		Targets:   targets,
		ItemDef:   item,
		name:      fmt.Sprintf("%s is using item '%s'", owner.Name, item.Name),
	}

	gWorld := reflect.ValueOf(scene.GameState.Globals["world"]).Interface().(*combat.WorldExtended)
	// Remove item here, otherwise 2 people could try and use the 1 potion
	gWorld.RemoveItem(item.Id, 1)

	c.Character.Controller.Change(csRunanim, csProne, false)
	storyboardEvents := []interface{}{
		//stateMachine, stateID, ...animID, additionalParams
		RunFunction(c.ShowItemNotice),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: 1}),
		RunState(c.Character.Controller, csRunanim, csUse, false),
		RunFunction(c.DoUseItem),
		Wait(1.3),
		RunState(c.Character.Controller, csMove, CSMoveParams{Dir: -1}),
		RunFunction(func() {
			c.DoFinish()
		}),
	}

	c.Storyboard = StoryboardCreate(scene.InternalStack, scene.win, storyboardEvents, false)
	return c
}

func (c *CEUseItem) Name() string {
	return c.name
}

func (c *CEUseItem) CountDown() float64 {
	return c.countDown
}

func (c *CEUseItem) CountDownSet(t float64) {
	c.countDown = t
}

func (c *CEUseItem) Owner() *combat.Actor {
	return c.owner
}

func (c *CEUseItem) Update() {

}

func (c *CEUseItem) IsFinished() bool {
	return c.mIsFinished
}

func (c *CEUseItem) DoFinish() {
	c.mIsFinished = true
}

func (c *CEUseItem) Execute(queue *EventQueue) {
	c.Scene.InternalStack.Push(c.Storyboard)
}

func (c CEUseItem) TimePoints(queue *EventQueue) float64 {
	speed := c.owner.Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}

func (c *CEUseItem) ShowItemNotice() {
	str := fmt.Sprintf("Item: %s", c.ItemDef.Name)
	c.Scene.ShowNotice(str)
}

func (c *CEUseItem) DoUseItem() {
	c.Scene.HideNotice()
	entity := Entities["fx_use_item"]
	pos := c.Character.Entity.GetSelectPosition()
	effect := AnimEntityFxCreate(pos.X, pos.Y, entity, entity.Frames, 0.1)
	c.Scene.AddEffect(effect)

	action := c.ItemDef.Use.Action
	CombatActions[action](c.Scene, c.owner, c.Targets, c.ItemDef)
}
