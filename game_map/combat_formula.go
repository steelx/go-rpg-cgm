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
	MeleeAttack func(state *CombatState, attacker, target *combat.Actor) (dmg float64, hit HitResult)
	BaseAttack  func(state *CombatState, attacker, target *combat.Actor) (dmg float64)
	CalcDamage  func(state *CombatState, attacker, target *combat.Actor) (dmg float64)
	IsHit       func(state *CombatState, attacker, target *combat.Actor) HitResult
	IsDodged    func(state *CombatState, attacker, target *combat.Actor) bool
	IsCountered func(state *CombatState, attacker, target *combat.Actor) bool
}

var Formula = FormulaT{
	MeleeAttack: meleeAttack,
	BaseAttack:  baseAttack,
	CalcDamage:  calcDamage,
	IsHit:       isHit,
	IsDodged:    isDodged,
	IsCountered: isCountered,
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
