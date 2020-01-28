package combat

import "github.com/steelx/go-rpg-cgm/world"

type ActorDef struct {
	Id               string //must match entityDef
	Stats            world.BaseStats
	StatGrowth       map[string]func() int
	Portrait         string
	Name             string
	Actions          []string
	ActiveEquipSlots []world.ItemType
	Equipment
}

type LevelUp struct {
	XP        float64
	Level     int
	BaseStats map[string]float64
}

//Must match to ItemsDB ID
type Equipment struct {
	Weapon,
	Armor,
	Access1,
	Access2 int //ItemsDB.Id
}
