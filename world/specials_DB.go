package world

const (
	SpecialSlash = "Slash"
	SpecialSteal = "Steal"
)

var SpecialsDB = map[string]SpecialItem{
	SpecialSlash: {
		Name:       "Slash",
		MpCost:     15,
		Action:     ElementSlash,
		TimePoints: 10,
		Target: ItemTarget{
			Selector:    SideEnemy,
			SwitchSides: false,
			Type:        CombatTargetTypeSIDE,
		},
	},

	SpecialSteal: {
		Name:       "Steal",
		MpCost:     0,
		Action:     ElementSteal,
		TimePoints: 10,
		Target: ItemTarget{
			Selector:    WeakestEnemy,
			SwitchSides: false,
			Type:        CombatTargetTypeONE,
		},
	},
}
