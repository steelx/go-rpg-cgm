package world

//• Key Item - a special item,
// usually required to progress past a certain part of the game.
//• Weapon - an item that can be equipped
// and has properties related to inflicting damage.
type Item struct {
	ItemType          ItemType
	Name, Description string
	Special           bool
	Stats             ItemStats
}
type ItemType int

const (
	empty ItemType = iota
	accessory
	usable
	weapon
	key
)

type ItemStats struct {
	strength, speed, intelligence, attack, defense, magic, resist int
}

var ItemsDB = make(map[int]Item)

func init() {
	//Item index key and ItemType are connected
	// e.g. ItemsDB[0] & empty = 0
	ItemsDB[0] = Item{
		ItemType:    empty,
		Name:        "",
		Description: "",
		Special:     false,
		Stats: ItemStats{
			strength:     0,
			speed:        0,
			intelligence: 0,
			attack:       0,
			defense:      0,
			magic:        0,
			resist:       0,
		},
	}

	ItemsDB[1] = Item{
		ItemType:    accessory,
		Name:        "Mysterious Torque",
		Description: "A golden torque that glitters.",
		Stats: ItemStats{
			strength: 10,
			speed:    10,
		},
	}

	ItemsDB[2] = Item{
		ItemType:    usable,
		Name:        "Heal Potion",
		Description: "Heals a little HP",
	}

	ItemsDB[3] = Item{
		ItemType:    weapon,
		Name:        "Bronze Sword",
		Description: "A short sword with weak blade",
		Stats: ItemStats{
			attack: 10,
		},
	}

	ItemsDB[4] = Item{
		ItemType:    key,
		Name:        "Old Bone",
		Description: "A human leg bone, open's a hidden room",
	}
}
