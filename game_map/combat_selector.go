package game_map

import (
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
)

type CombatSelectorFunc struct {
	RandomAlivePlayer,
	WeakestEnemy,
	SideEnemy,
	SelectAll func(state *CombatState) []*combat.Actor
}

var CombatSelectorMap = map[string]func(state *CombatState) []*combat.Actor{
	world.RandomAlivePlayer: RandomAlivePlayer,
	world.WeakestEnemy:      WeakestEnemy,
	world.SideEnemy:         SideEnemy,
	world.SelectAll:         SelectAll,

	world.MostHurtEnemy: func(state *CombatState) []*combat.Actor {
		return WeakestActor(state.Actors[enemies], true)
	},
	world.MostHurtParty: func(state *CombatState) []*combat.Actor {
		return WeakestActor(state.Actors[party], true)
	},
	world.MostDrainedParty: func(state *CombatState) []*combat.Actor {
		return MostDrainedActor(state.Actors[party], true)
	},
	world.DeadParty: func(state *CombatState) []*combat.Actor {
		return DeadActors(state.Actors[party])
	},
}

var CombatSelector = CombatSelectorFunc{
	RandomAlivePlayer: RandomAlivePlayer,
	WeakestEnemy:      WeakestEnemy,
	SideEnemy:         SideEnemy,
	SelectAll:         SelectAll,
}

func RandomAlivePlayer(state *CombatState) []*combat.Actor {
	aliveList := make([]*combat.Actor, 0)
	for _, v := range state.Actors[party] {
		if v.Stats.Get("HpNow") > 0 {
			aliveList = append(aliveList, v)
		}
	}
	if len(aliveList) == 1 {
		return []*combat.Actor{aliveList[0]}
	}
	randIndex := utilz.RandInt(0, len(aliveList)-1)
	return []*combat.Actor{aliveList[randIndex]}
}

func WeakestEnemy(state *CombatState) []*combat.Actor {
	enemyList := state.Actors[enemies]
	health := 99999.9

	var target *combat.Actor
	for _, v := range enemyList {
		hpNow := v.Stats.Get("HpNow")
		if hpNow < health {
			health = hpNow
			target = v
		}
	}
	return []*combat.Actor{target}
}

func SideEnemy(state *CombatState) []*combat.Actor {
	return state.Actors[enemies]
}

func SelectAll(state *CombatState) []*combat.Actor {
	return append(state.Actors[enemies], state.Actors[party]...)
}

func WeakestActor(actors []*combat.Actor, onlyCheckHurt bool) []*combat.Actor {
	var target *combat.Actor = nil
	health := 99999.9

	for _, v := range actors {
		hp := v.Stats.Get("HpNow")
		isHurt := false
		if hp < v.Stats.Get("HpMax") {
			isHurt = true
		}
		skip := false
		if onlyCheckHurt && !isHurt {
			skip = true
		}
		if hp < health && !skip {
			health = hp
			target = v
		}
	}

	if target != nil {
		return []*combat.Actor{target}
	}

	return []*combat.Actor{actors[0]}
}

func MostDrainedActor(actors []*combat.Actor, onlyCheckDrained bool) []*combat.Actor {
	var target *combat.Actor = nil
	magic := 99999.9

	for _, v := range actors {
		mp := v.Stats.Get("MpNow")
		isDrained := false
		if mp < v.Stats.Get("MpMax") {
			isDrained = true
		}
		skip := false
		if onlyCheckDrained && !isDrained {
			skip = true
		}
		if mp < magic && !skip {
			magic = mp
			target = v
		}
	}
	if target != nil {
		return []*combat.Actor{target}
	}

	return []*combat.Actor{actors[0]}
}

func DeadActors(actors []*combat.Actor) []*combat.Actor {
	for _, v := range actors {
		hp := v.Stats.Get("HpNow")
		if hp <= 0 {
			return []*combat.Actor{v}
		}
	}

	return []*combat.Actor{actors[0]}
}
