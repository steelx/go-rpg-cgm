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
		`Weapon    :`,
		`Armor     :`,
		`Accessory :`,
		`Accessory :`,
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
		`Strength     :`,
		`Speed        :`,
		`Intelligence :`,
	},
	ItemStatLabels: []string{
		`Attack  :`,
		`Defense :`,
		`Magic   :`,
		`Resist  :`,
	},
	ActionLabels: ActionLabels{
		Attack: "Attack",
		Item:   "Item",
	},
}

//maps to EquipSlotId
var equipment = map[string]int{
	"Weapon":  0,
	"Armor":   1,
	"Access1": 2,
	"Access2": 3,
}
