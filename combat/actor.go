package combat

import (
	"github.com/steelx/go-rpg-cgm/dice"
)

type GrowthT struct {
	Fast func() int
	Med  func() int
	Slow func() int
}

var Growth = GrowthT{
	Fast: dice.Create("3d2"),
	Med:  dice.Create("1d3"),
	Slow: dice.Create("1d2"),
}
var DefaultStats = BaseStats{
	HpNow:    300,
	HpMax:    300,
	MpNow:    300,
	MpMax:    300,
	Strength: 10, Speed: 10, Intelligence: 10,
}

type Actor struct {
	Stats      Stats
	StatGrowth map[string]func() int

	Level           int
	XP, NextLevelXP float64
	Def             ActorDef
}

type ActorDef struct {
	Stats      BaseStats
	StatGrowth map[string]func() int
}

type LevelUp struct {
	XP        float64
	Level     int
	BaseStats map[string]float64
}

/* example: ActorCreate(HeroDef)
var HeroDef = combat.ActorDef{
		Stats: combat.DefaultStats,
		StatGrowth: map[string]func() int{
			"HpMax":        dice.Create("4d50+100"),
			"MpMax":        dice.Create("2d50+100"),
			"Strength":     combat.Growth.Fast,
			"Speed":        combat.Growth.Fast,
			"Intelligence": combat.Growth.Med,
		},
	}
*/
// ActorCreate
func ActorCreate(def ActorDef) Actor {
	a := Actor{
		Def:        def,
		StatGrowth: def.StatGrowth,
		Stats:      StatsCreate(def.Stats),
		XP:         0,
		Level:      1,
	}

	a.NextLevelXP = NextLevel(a.Level)
	return a
}

func (a Actor) ReadyToLevelUp() bool {
	return a.XP >= a.NextLevelXP
}

func (a *Actor) AddXP(xp float64) bool {
	a.XP += xp
	return a.ReadyToLevelUp()
}

func (a Actor) CreateLevelUp() LevelUp {
	levelUp := LevelUp{
		XP:        -a.NextLevelXP,
		Level:     1,
		BaseStats: make(map[string]float64),
	}

	for id, diceRoll := range a.StatGrowth {
		levelUp.BaseStats[id] = float64(diceRoll())
	}

	//Pending feature
	// Additional level up code
	// e.g. if you want to apply
	// a bonus every 4 levels
	// or heal the players MP/HP

	return levelUp
}

func (a *Actor) ApplyLevel(levelUp LevelUp) {
	a.XP += levelUp.XP
	a.Level += levelUp.Level
	a.NextLevelXP = NextLevel(a.Level)

	for k, v := range levelUp.BaseStats {
		a.Stats.Base[k] += v
	}

	//Pending feature
	// Unlock any special abilities etc.
}
