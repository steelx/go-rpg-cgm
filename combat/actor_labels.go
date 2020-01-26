package combat

type ActorLabel struct {
	EquipSlotLabels []string
	EquipSlotId     []string
	ActorStats      []string
	ItemStats       []string
	ActorStatLabels []string
	ItemStatLabels  []string
	ActionLabels    ActionLabels
}
type ActionLabels struct {
	Attack, Item string
}

var ActorLabels = ActorLabel{
	EquipSlotLabels: []string{
		`Weapon `,
		`Armor `,
		`Accessory 1 `,
		`Accessory 2 `,
	},
	EquipSlotId: []string{
		"Weapon",
		"Armor",
		"Access1",
		"Access2",
	},
	//BaseStats stats.go
	ActorStats: []string{
		"Strength",
		"Speed",
		"Intelligence",
	},
	ItemStats: []string{
		"Attack",
		"Defense",
		"Magic",
		"Resist",
	},
	ActorStatLabels: []string{
		`Strength `,
		`Speed `,
		`Intelligence `,
	},
	ItemStatLabels: []string{
		`Attack `,
		`Defense `,
		`Magic `,
		`Resist `,
	},
	ActionLabels: ActionLabels{
		Attack: "Attack",
		Item:   "Item",
	},
}
