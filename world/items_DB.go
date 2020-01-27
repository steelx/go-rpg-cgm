package world

//• Key Item - a special item,
// usually required to progress past a certain part of the game.
//• Weapon - an item that can be equipped
// and has properties related to inflicting damage.
type Item struct {
	ItemType          ItemType
	Name, Description string
	Special           bool
	Stats             Mod
	Restrictions      []string //e.g. {"hero","mage",}
	Use               UseAction
}

type Action int

const (
	Revive Action = iota
	Heal
)

type ActionTarget int

const (
	Any ActionTarget = iota
	FriendlyDead
	Enemy
)

type UseAction struct {
	Action        Action
	Target        ActionTarget //character type
	TargetDefault ActionTarget
	Hint          string
}

type ItemType int

const (
	Empty ItemType = iota
	Usable
	Accessory
	Weapon
	Sword
	Dagger
	Stave
	Armor
	Plate
	Leather
	Robe
	UpArrow
	DownArrow
)

var ItemsDB = make(map[int]Item)

func init() {

	ItemsDB[0] = Item{
		ItemType:    Empty,
		Name:        "empty",
		Description: "",
		Special:     false,
		Stats: Mod{
			Add: BaseStats{
				Strength:     0,
				Speed:        0,
				Intelligence: 0,
				Attack:       0,
				Defense:      0,
				Magic:        0,
				Resist:       0,
			},
			Mult: BaseStats{},
		},
	}

	ItemsDB[1] = Item{
		ItemType:     Weapon,
		Name:         "Bone Blade",
		Description:  "A wicked sword made from bone.",
		Restrictions: []string{"hero"},
		Stats: Mod{
			Add: BaseStats{
				Attack: 5,
			},
		},
	}

	ItemsDB[2] = Item{
		ItemType:     Armor,
		Name:         "Bone Armor",
		Description:  "Armor made from plates of blackened bone.",
		Restrictions: []string{"hero"},
		Stats: Mod{
			Add: BaseStats{
				Defense: 5,
				Resist:  1,
			},
		},
	}

	ItemsDB[3] = Item{
		ItemType:    Accessory,
		Name:        "Ring of Titan",
		Description: "Grants the strength of the Titan.",
		Stats: Mod{
			Add: BaseStats{
				Strength: 10,
			},
		},
	}

	ItemsDB[4] = Item{
		ItemType:    Usable,
		Name:        "Old Bone",
		Description: "A human leg bone, open's a hidden room",
	}

	ItemsDB[5] = Item{
		ItemType:     Weapon,
		Name:         "World Tree Branch",
		Description:  "A hard wood branch.",
		Restrictions: []string{"mage"},
		Stats: Mod{
			Add: BaseStats{
				Attack: 2,
				Magic:  5,
			},
		},
	}

	ItemsDB[6] = Item{
		ItemType:     Armor,
		Name:         "Dragon's Cloak",
		Description:  "A cloak of dragon scales.",
		Restrictions: []string{"mage"},
		Stats: Mod{
			Add: BaseStats{
				Defense: 3,
				Resist:  10,
			},
		},
	}

	ItemsDB[7] = Item{
		ItemType:    Accessory,
		Name:        "Singer's Stone",
		Description: "The stone's song resists magical attacks.",
		Stats: Mod{
			Add: BaseStats{
				Resist: 10,
			},
		},
	}

	ItemsDB[8] = Item{
		ItemType:     Weapon,
		Name:         "Black Dagger",
		Description:  "A dagger made out of an unknown material.",
		Restrictions: []string{"thief"},
		Stats: Mod{
			Add: BaseStats{
				Attack: 4,
			},
		},
	}

	ItemsDB[9] = Item{
		ItemType:     Armor,
		Name:         "Footpad Leathers",
		Description:  "Light Armor for silent movement.",
		Restrictions: []string{"thief"},
		Stats: Mod{
			Add: BaseStats{
				Defense: 4,
			},
		},
	}

	ItemsDB[10] = Item{
		ItemType:    Accessory,
		Name:        "Swift Boots",
		Description: "Increases speed by 25%",
		Stats: Mod{
			Mult: BaseStats{
				Speed: 0.25,
			},
		},
	}

	ItemsDB[11] = Item{
		ItemType:    Usable,
		Name:        "Heal Potion",
		Description: "Heal a small amount of HP.",
		Use: UseAction{
			Action:        Revive,
			Target:        Any,
			TargetDefault: FriendlyDead,
			Hint:          "Choose target to revive.",
		},
	}
}
