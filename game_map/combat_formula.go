package game_map

import (
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/utilz"
	"math"
)

type HitResult int

const (
	HitResultMiss HitResult = iota
	HitResultDodge
	HitResultHit
	HitResultCritical
)

type FormulaT struct {
	MeleeAttack      func(state *CombatState, attacker, target *combat.Actor) (dmg float64, hit HitResult)
	BaseAttack       func(state *CombatState, attacker, target *combat.Actor) (dmg float64)
	CalcDamage       func(state *CombatState, attacker, target *combat.Actor) (dmg float64)
	IsHit            func(state *CombatState, attacker, target *combat.Actor) HitResult
	IsDodged         func(state *CombatState, attacker, target *combat.Actor) bool
	IsCountered      func(state *CombatState, attacker, target *combat.Actor) bool
	CanFlee          func(state *CombatState, target *combat.Actor) bool
	MostHurtEnemy    func(state *CombatState) []*combat.Actor
	MostHurtParty    func(state *CombatState) []*combat.Actor
	MostDrainedParty func(state *CombatState) []*combat.Actor
	DeadParty        func(state *CombatState) []*combat.Actor
}

var Formula = FormulaT{
	MeleeAttack: meleeAttack,
	BaseAttack:  baseAttack,
	CalcDamage:  calcDamage,
	IsHit:       isHit,
	IsDodged:    isDodged,
	IsCountered: isCountered,
	CanFlee:     canFlee,
	MostHurtEnemy: func(state *CombatState) []*combat.Actor {
		return WeakestActor(state.Actors[enemies], true)
	},
	MostHurtParty: func(state *CombatState) []*combat.Actor {
		return WeakestActor(state.Actors[party], true)
	},
	MostDrainedParty: func(state *CombatState) []*combat.Actor {
		return MostDrainedActor(state.Actors[party], true)
	},
	DeadParty: func(state *CombatState) []*combat.Actor {
		return DeadActors(state.Actors[party])
	},
}

func meleeAttack(state *CombatState, attacker, target *combat.Actor) (dmg float64, hit HitResult) {
	//stats := attacker.Stats
	//enemyStats := target.Stats

	var damage float64
	hitResult := isHit(state, attacker, target)

	if hitResult == HitResultMiss {
		return math.Floor(damage), HitResultMiss
	}

	if isDodged(state, attacker, target) {
		return math.Floor(damage), HitResultDodge
	}

	damage = calcDamage(state, attacker, target)

	if hitResult == HitResultHit {
		return math.Floor(damage), HitResultHit
	}

	// Critical
	damage = damage + baseAttack(state, attacker, target)
	return math.Floor(damage), HitResultCritical
}

func isHit(state *CombatState, attacker, target *combat.Actor) HitResult {
	stats := attacker.Stats
	speed := stats.Get("Speed")
	intelligence := stats.Get("Intelligence")

	cth := 0.8 //Chance to Hit
	ctc := 0.1 //Chance to Crit

	bonus := ((speed + intelligence) / 2) / 255 //divide by max value
	cth = cth + (bonus / 2)

	rand := utilz.RandFloat(0, 1)
	isHit := rand <= cth
	isCrit := rand <= ctc

	if isCrit {
		return HitResultCritical
	} else if isHit {
		return HitResultHit
	} else {
		return HitResultMiss
	}
}

func isDodged(state *CombatState, attacker, target *combat.Actor) bool {
	stats := attacker.Stats
	enemyStats := target.Stats

	speed := stats.Get("Speed")
	enemySpeed := enemyStats.Get("Speed")

	ctd := 0.03 //Chance to Dodge
	speedDiff := speed - enemySpeed
	// clamp speed diff to plus or minus 10%
	speedDiff = utilz.Clamp(speedDiff, -10, 10) * 0.01

	ctd = math.Max(0, ctd+speedDiff)

	return utilz.RandFloat(0, 1) <= ctd
}

func isCountered(state *CombatState, attacker, target *combat.Actor) bool {
	// if not assigned 0 is returned, which will mean no chance of countering
	counter := target.Stats.Get("Counter")

	// I want random to be between 0 and under 1
	// This means 1 always counters and 0 it never happens
	return utilz.RandFloat(0, 1)*0.99999 < counter
}

func baseAttack(state *CombatState, attacker, target *combat.Actor) (dmg float64) {
	stats := attacker.Stats
	strength := stats.Get("Strength")
	attackStat := stats.Get("Attack")

	attack := (strength / 2) + attackStat
	return utilz.RandFloat(attack, attack*2)
}

func calcDamage(state *CombatState, attacker, target *combat.Actor) (dmg float64) {
	targetStats := target.Stats
	defense := targetStats.Get("Defense")

	attack := baseAttack(state, attacker, target)
	return math.Floor(math.Max(0, attack-defense))
}

func canFlee(state *CombatState, target *combat.Actor) bool {
	fc := 0.35 // flee chance
	stats := target.Stats
	speed := stats.Get("Speed")

	// Get the average speed of the enemies
	var enemyCount, totalSpeed float64
	for _, v := range state.Actors[enemies] {
		speed := v.Stats.Get("Speed")
		totalSpeed += speed
		enemyCount += 1
	}
	avgSpeed := totalSpeed / enemyCount
	if speed > avgSpeed {
		fc = fc + 0.15
	} else {
		fc = fc - 0.15
	}
	return utilz.RandFloat(0, 1) <= fc
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
