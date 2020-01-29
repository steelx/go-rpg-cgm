package combat

import "github.com/steelx/go-rpg-cgm/world"

type ActorLabel struct {
	EquipSlotLabels []string
	EquipSlotId     []string
	ActorStats      []string
	ItemStats       []string
	ActorStatLabels []string
	ItemStatLabels  []string
	ActionLabels    ActionLabels
	EquipSlotTypes  map[world.ItemType]string
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
		"Accessory1",
		"Accessory2",
	},
	EquipSlotTypes: map[world.ItemType]string{
		world.Weapon:    "Weapon",
		world.Armor:     "Armor",
		world.Accessory: "Accessory",
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
