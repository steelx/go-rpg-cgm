package game_map

import (
	"fmt"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/world"
	"math"
	"reflect"
)

var CombatActions = map[world.Action]func(state *CombatState, owner *combat.Actor, targets []*combat.Actor, defI interface{}){
	world.HpRestore:    HpRestore,
	world.MpRestore:    MpRestore,
	world.Revive:       Revive,
	world.ElementSpell: elementSpell,
}

func HpRestore(state *CombatState, owner *combat.Actor, targets []*combat.Actor, defI interface{}) {
	def := reflect.ValueOf(defI).Interface().(world.Item)
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

func MpRestore(state *CombatState, owner *combat.Actor, targets []*combat.Actor, defI interface{}) {
	def := reflect.ValueOf(defI).Interface().(world.Item)
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

func Revive(state *CombatState, owner *combat.Actor, targets []*combat.Actor, defI interface{}) {
	def := reflect.ValueOf(defI).Interface().(world.Item)
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
	pos := entity.GetSelectPosition()
	x := pos.X
	y := pos.Y

	effect := AnimEntityFxCreate(x, y, fxEntityDef, fxEntityDef.Frames, spf)
	state.AddEffect(effect)
}

func AddTextNumberEffect(state *CombatState, entity *Entity, num float64, hexColor string) {
	pos := entity.GetSelectPosition()
	x, y := pos.X, pos.Y-entity.Height/2

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

func elementSpell(state *CombatState, owner *combat.Actor, targets []*combat.Actor, defI interface{}) {
	def := reflect.ValueOf(defI).Interface().(world.SpecialItem)

	for _, v := range targets {
		_, _, entity := StatsCharEntity(state, v)
		damage, hitResult := MagicAttack(state, owner, v, def)
		if hitResult == HitResultHit {
			state.ApplyDamage(v, damage, true)
		}

		if def.Element == world.SpellFire {
			AddAnimEffect(state, entity, Entities["fx_fire"], 0.06)
		} else if def.Element == world.SpellBolt {
			AddAnimEffect(state, entity, Entities["fx_electric"], 0.12)
		} else if def.Element == world.SpellIce {
			AddAnimEffect(state, entity, Entities["fx_ice_1"], 0.1)
			pos := entity.GetSelectPosition()
			x := pos.X
			y := pos.Y

			spark := Entities["fx_ice_spark"]
			effect := AnimEntityFxCreate(x, y, spark, spark.Frames, 0.12)
			state.AddEffect(effect)

			x2 := x + entity.Width*0.8
			ice2 := Entities["fx_ice_2"]
			effect = AnimEntityFxCreate(x2, y, ice2, ice2.Frames, 0.1)
			state.AddEffect(effect)

			x3 := x - entity.Width*0.8
			y3 := y - entity.Height*0.6
			ice3 := Entities["fx_ice_3"]
			effect = AnimEntityFxCreate(x3, y3, ice3, ice3.Frames, 0.1)
			state.AddEffect(effect)
		}

	}
}
