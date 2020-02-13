package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
)

func (c *CombatState) AddTurns(actorList []*combat.Actor) {
	for _, v := range actorList {
		hpNow := v.Stats.Get("HpNow")
		if hpNow > 0 && !c.EventQueue.ActorHasEvent(v) {
			event := CETurnCreate(c, v)
			tp := event.TimePoints(c.EventQueue)
			c.EventQueue.Add(event, tp)
		}
	}
}

func (c *CombatState) GetTarget(owner *combat.Actor) *combat.Actor {
	if owner.IsPlayer() {
		return c.Actors[enemies][len(c.Actors[enemies])-1]
	}

	return c.Actors[party][len(c.Actors[enemies])-1]
}

func (c CombatState) GetAlivePartyActors() []*combat.Actor {
	var alive []*combat.Actor
	for _, a := range c.Actors[party] {
		if !a.IsKOed() {
			alive = append(alive, a)
		}
	}
	return alive
}

//OnDead makes Actor KnockOut
func (c *CombatState) OnDead(actor *combat.Actor) {
	if actor.IsPlayer() {
		actor.KO()
	} else {
		for i := len(c.Actors[enemies]) - 1; i >= 0; i-- {
			if actor == c.Actors[enemies][i] {
				c.Actors[enemies] = removeActorAtIndex(c.Actors[enemies], i)
			}
		}
	}

	//Remove owned events
	c.EventQueue.RemoveEventsOwnedBy(actor)

	if c.IsPartyDefeated() {
		fmt.Println("CombatState OnDead: Party loses")
	} else if c.IsEnemyDefeated() {
		fmt.Println("CombatState OnDead: Enemy loses")
	}
}

func removeActorAtIndex(arr []*combat.Actor, i int) []*combat.Actor {
	return append(arr[:i], arr[i+1:]...)
}
func (c CombatState) removeCharAtIndex(arr []*Character, i int) []*Character {
	return append(arr[:i], arr[i+1:]...)
}
func (c CombatState) removeFxAtIndex(arr []EffectState, i int) []EffectState {
	return append(arr[:i], arr[i+1:]...)
}
func (c *CombatState) insertFxAtIndex(index int, fxI EffectState) {
	temp := append([]EffectState{}, c.EffectList[index:]...)
	c.EffectList = append(c.EffectList[0:index], fxI)
	c.EffectList = append(c.EffectList, temp...)
}

//IsPartyDefeated check's at least 1 Actor is standing return false
func (c CombatState) IsPartyDefeated() bool {
	for _, actor := range c.Actors[party] {
		if !actor.IsKOed() {
			return false
		}
	}
	return true
}
func (c CombatState) PartyWins() bool {
	return !c.HasLiveActors(c.Actors[enemies])
}

func (c CombatState) IsEnemyDefeated() bool {
	return len(c.Actors[enemies]) == 0
}
func (c CombatState) EnemyWins() bool {
	return !c.HasLiveActors(c.Actors[party])
}

func (c CombatState) HasLiveActors(actorList []*combat.Actor) bool {
	for _, v := range actorList {
		hpNow := v.Stats.Get("HpNow")
		if hpNow > 0 {
			return true
		}
	}
	return false
}

func (c *CombatState) IsPartyMember(owner *combat.Actor) bool {
	for _, v := range c.Actors[party] {
		if v == owner {
			return true
		}
	}
	return false
}
