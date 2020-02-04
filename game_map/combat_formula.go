package game_map

import (
	"github.com/steelx/go-rpg-cgm/combat"
	"math"
)

type FormulaT struct {
	MeleeAttack func(state *CombatState, attacker, target *combat.Actor) (dmg float64)
}

var Formula = FormulaT{
	MeleeAttack: func(state *CombatState, attacker, target *combat.Actor) (dmg float64) {
		stats := attacker.Stats
		enemyStats := target.Stats
		// Simple attack get
		attack := stats.Get("Attack")
		attack = attack + stats.Get("Strength")
		defense := enemyStats.Get("Defense")
		dmg = math.Max(0, attack-defense)
		return
	},
}
