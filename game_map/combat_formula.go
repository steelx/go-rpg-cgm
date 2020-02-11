package game_map

import (
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
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

func IsHitMagic(state *CombatState, attacker, target *combat.Actor, spell world.SpecialItem) HitResult {
	// Spell hit information determined by the spell
	hitChance := spell.BaseHitChance
	if utilz.RandFloat(0, 1) <= hitChance {
		return HitResultHit
	}

	return HitResultMiss
}

func CalcSpellDamage(state *CombatState, attacker, target *combat.Actor, spell world.SpecialItem) (damage float64) {
	// Find the basic damage
	base := utilz.RandFloat(spell.BaseDamage[0], spell.BaseDamage[1])
	damage = base * 4
	// Increase power of spell by caster
	level := attacker.Level
	stats := attacker.Stats
	bonus := float64(level) * stats.Get("Intelligence") * (base / 32)

	damage += bonus

	// Apply elemental weakness / strength modifications
	if spell.Element != "" {
		modifier := target.Stats.Get(spell.Element)
		damage += damage * modifier
	}
	// Handle resistance [0..255]
	resist := math.Min(255, target.Stats.Get("Resist"))
	resist01 := 1 - (resist / 255)
	damage = damage * resist01
	return damage
}

func MagicAttack(state *CombatState, attacker, target *combat.Actor, spell world.SpecialItem) (float64, HitResult) {
	damage := 0.0
	hitResult := IsHitMagic(state, attacker, target, spell)
	if hitResult == HitResultMiss {
		return damage, HitResultMiss
	}

	// Dodging spells not allowed.
	damage = CalcSpellDamage(state, attacker, target, spell)
	return math.Floor(damage), HitResultHit
}

func Steal(state *CombatState, attacker, target *combat.Actor) bool {
	cts := 0.05 // 5%

	if attacker.Level > target.Level {
		cts = float64(50+attacker.Level-target.Level) / 128
		cts = utilz.Clamp(cts, 0.05, 0.95)
	}

	randN := utilz.RandFloat(0, 1) //wondering if should be 0 to 1 or higher
	return randN <= cts
}
