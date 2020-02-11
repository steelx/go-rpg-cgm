package world

/*
	• name - The name displayed in the browse spell box.
	• action - Name of a combat action that applies the spell effect to the combat
	state.
	• element - Extra data for the element_spell action. Describes the element of
	the spell. Optional.
	• mp_cost - How much mana is required to cast the spell.
	• cast_time - The time points required to cast the spell.
	• base_damage - The basic range of damage to feed into the spell calculation.
	This can also be a single number.
	• base_hit_chance - Spell’s basic chance to hit; 1 here mean 100% chance.
	• target - The targeting information, much the same as in the item definitions.
*/

const (
	SpellFire = "Fire"
	SpellBurn = "Burn"
	SpellIce  = "Ice"
	SpellBolt = "Bolt"
)

type SpecialItem struct {
	Name    string
	Action  Action
	Element string
	MpCost,
	CastTime,
	BaseHitChance,
	TimePoints float64
	BaseDamage [2]float64 // multiplied by level
	Target     ItemTarget
}

// spell cast time 1 is base, 2 is twice as long etc
var SpellsDB = map[string]SpecialItem{
	SpellFire: {
		Name:          "Fire",
		Action:        ElementSpell,
		Element:       SpellFire,
		MpCost:        8,
		CastTime:      0.7,
		BaseDamage:    [2]float64{3, 5}, // multiplied by level
		BaseHitChance: 1,
		TimePoints:    10,
		Target: ItemTarget{
			Selector:    WeakestEnemy,
			SwitchSides: true,
			Type:        CombatTargetTypeONE,
		},
		//Damage = Spell Power * 4 + (Level * Magic Power * Spell Power / 32)
	},
	SpellBurn: {
		Name:          "Burn",
		Action:        ElementSpell,
		Element:       SpellFire,
		MpCost:        16,
		CastTime:      0.9,
		BaseDamage:    [2]float64{3, 6},
		BaseHitChance: 1,
		TimePoints:    20,
		Target: ItemTarget{
			Selector:    SideEnemy,
			SwitchSides: true,
			Type:        CombatTargetTypeSIDE,
		},
	},

	SpellIce: {
		Name:          "Ice",
		Action:        ElementSpell,
		Element:       SpellIce,
		MpCost:        8,
		CastTime:      1,
		BaseDamage:    [2]float64{7, 17},
		BaseHitChance: 1,
		TimePoints:    10,
		Target: ItemTarget{
			Selector:    WeakestEnemy,
			SwitchSides: true,
			Type:        CombatTargetTypeONE,
		},
	},

	SpellBolt: {
		Name:          "Electric bolt",
		Action:        ElementSpell,
		Element:       SpellBolt,
		MpCost:        8,
		CastTime:      0.5,
		BaseDamage:    [2]float64{4, 14},
		BaseHitChance: 1,
		TimePoints:    10,
		Target: ItemTarget{
			Selector:    WeakestEnemy,
			SwitchSides: true,
			Type:        CombatTargetTypeONE,
		},
	},
}
