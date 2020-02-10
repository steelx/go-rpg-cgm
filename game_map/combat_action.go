package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/world"
	"math"
)

var CombatActions = map[world.Action]func(state *CombatState, owner *combat.Actor, targets []*combat.Actor, def world.Item){
	world.HpRestore: HpRestore,
	world.MpRestore: MpRestore,
	world.Revive:    Revive,
}

func HpRestore(state *CombatState, owner *combat.Actor, targets []*combat.Actor, def world.Item) {
	restoreAmount := def.Use.Restore
	animEffect := Entities["fx_restore_hp"]
	restoreColor := "#00ff45"

	for _, v := range targets {
		stats, _, entity := StatsCharEntity(state, v)
		maxHP := stats.Get("HpMax")
		nowHP := stats.Get("HpNow")

		if nowHP > 0 {
			AddTextNumberEffect(state, entity, restoreAmount, restoreColor)
			nowHP = math.Min(maxHP, nowHP+restoreAmount)
			stats.Set("HpNow", nowHP)
		}

		AddAnimEffect(state, entity, animEffect, 0.1)
	}
}

func MpRestore(state *CombatState, owner *combat.Actor, targets []*combat.Actor, def world.Item) {
	restoreAmount := def.Use.Restore
	animEffect := Entities["fx_restore_mp"]
	restoreColor := "#00ffff"

	for _, v := range targets {
		stats, _, entity := StatsCharEntity(state, v)
		maxMP := stats.Get("MpMax")
		nowMP := stats.Get("MpNow")
		nowHP := stats.Get("HpNow")

		if nowHP > 0 {
			AddTextNumberEffect(state, entity, restoreAmount, restoreColor)
			nowMP = math.Min(maxMP, nowMP+restoreAmount)
			stats.Set("MpNow", nowHP)
		}

		AddAnimEffect(state, entity, animEffect, 0.1)
	}
}

func Revive(state *CombatState, owner *combat.Actor, targets []*combat.Actor, def world.Item) {
	restoreAmount := def.Use.Restore
	animEffect := Entities["fx_revive"]
	restoreColor := "#00ff00"

	for _, v := range targets {
		stats, character, entity := StatsCharEntity(state, v)
		maxHP := stats.Get("HpMax")
		nowHP := stats.Get("HpNow")

		if nowHP <= 0 {
			nowHP = math.Min(maxHP, nowHP+restoreAmount)

			// the character will get a CETurn event automatically
			// assigned next update
			character.Controller.Change(csStandby, csStandby)

			stats.Set("HpNow", nowHP)
			AddTextNumberEffect(state, entity, restoreAmount, restoreColor)
		}

		AddAnimEffect(state, entity, animEffect, 0.1)
	}
}

func AddAnimEffect(state *CombatState, entity *Entity, fxEntityDef EntityDefinition, spf float64) {
	x := entity.X
	y := entity.Y + (entity.Height * 0.75)

	effect := AnimEntityFxCreate(x, y, fxEntityDef, fxEntityDef.Frames, spf)
	state.AddEffect(effect)
}

func AddTextNumberEffect(state *CombatState, entity *Entity, num float64, hexColor string) {
	x := entity.X
	y := entity.Y

	fxText := fmt.Sprintf("+%v", num)
	textEffect := CombatTextFXCreate(x, y, fxText, hexColor)
	state.AddEffect(textEffect)
}

func StatsCharEntity(state *CombatState, actor *combat.Actor) (world.Stats, *Character, *Entity) {
	stats := actor.Stats
	character := state.ActorCharMap[actor]
	entity := character.Entity
	return stats, character, entity
}
