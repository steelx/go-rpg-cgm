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
