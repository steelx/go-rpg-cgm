package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
)

type CEFlee struct {
	Scene      *CombatState
	Character  *Character
	owner      *combat.Actor
	name       string
	countDown  float64
	finished   bool
	FleeParams CSMoveParams
	CanFlee    bool
	Storyboard *Storyboard
}

func CEFleeCreate(scene *CombatState, owner *combat.Actor, fleeParams CSMoveParams) *CEFlee {
	//CSMoveParams{Dir: -1, Distance: 180, Time: 0.6}
	c := &CEFlee{
		Scene:      scene,
		owner:      owner,
		Character:  scene.ActorCharMap[owner],
		FleeParams: fleeParams,
		name:       fmt.Sprintf("Flee for %s", owner.Name),
	}

	c.Character.Facing = CharacterFacingDirection[1] //right
	c.Character.Controller.Change(csRunanim, csProne, false)
	var storyboardEvents []interface{}

	//Scene CanFlee override
	if c.Scene.CanFlee {
		c.CanFlee = Formula.CanFlee(scene, owner)
	} else {
		c.CanFlee = false
	}

	if c.CanFlee {
		storyboardEvents = []interface{}{
			//stateMachine, stateID, ...animID, additionalParams
			RunFunction(func() {
				c.Scene.ShowNotice("Attempting to Flee...")
			}),
			Wait(1),
			RunFunction(func() {
				c.Scene.ShowNotice("Success")
				c.Character.Controller.Change(csMove, c.FleeParams)
			}),
			Wait(1),
			RunFunction(c.DoFleeSuccess),
			Wait(0.6),
		}
	} else {
		storyboardEvents = []interface{}{
			RunFunction(func() {
				c.Scene.ShowNotice("Attempting to Flee...")
			}),
			Wait(1),
			RunFunction(func() {
				c.Scene.ShowNotice("Failed !")
			}),
			Wait(1),
			RunFunction(c.OnFleeFail),
		}
	}

	c.Storyboard = StoryboardCreate(scene.InternalStack, scene.win, storyboardEvents, false)

	return c
}

func (c CEFlee) Name() string {
	return c.name
}

func (c CEFlee) CountDown() float64 {
	return c.countDown
}

func (c *CEFlee) CountDownSet(t float64) {
	c.countDown = t
}

func (c CEFlee) Owner() *combat.Actor {
	return c.owner
}

func (c *CEFlee) Update() {

}

func (c CEFlee) IsFinished() bool {
	return c.finished
}

func (c *CEFlee) Execute(queue *EventQueue) {
	c.Scene.InternalStack.Push(c.Storyboard)
}

func (c CEFlee) TimePoints(queue *EventQueue) float64 {
	speed := c.owner.Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}

func (c *CEFlee) OnFleeFail() {
	c.Character.Facing = CharacterFacingDirection[3]    //left
	c.Character.Controller.Change(csStandby, csStandby) //animId
	c.finished = true
	c.Scene.HideNotice()
}

func (c *CEFlee) DoFleeSuccess() {
	for _, v := range c.Scene.Actors[party] {
		alive := v.Stats.Get("HpNow") > 0
		var isFleer bool
		if v == c.owner {
			isFleer = true
		}

		if alive && !isFleer {
			char := c.Scene.ActorCharMap[v]
			char.Facing = CharacterFacingDirection[1]
			char.Controller.Change(csMove, c.FleeParams)
		}
	}

	c.Scene.OnFlee()
	c.Scene.HideNotice()
}
