package world

var magicHat = Modifier{
	UniqueId: 1,
	Mod: Mod{
		Add: BaseStats{
			Strength: 5,
		},
	},
}

var magicSword = Modifier{
	UniqueId: 2,
	Mod: Mod{
		Add: BaseStats{
			Strength: 5,
		},
	},
}

var curse = Modifier{
	UniqueId: 3,
	Mod: Mod{
		Add: BaseStats{
			Strength:     -0.5,
			Speed:        -0.5,
			Intelligence: -0.5,
		},
	},
}

var spellBravery = Modifier{
	UniqueId: 4,
	Mod: Mod{
		Add: BaseStats{
			Strength:     0.5,
			Speed:        0.5,
			Intelligence: 0.5,
		},
	},
}

//weapon
var halfSword = Modifier{
	Name:     "Half Sword",
	UniqueId: 5,
	Mod: Mod{
		Add: BaseStats{
			Attack: 5,
		},
	},
}

//weapon
var spikedHelmOfNecromancy = Modifier{
	Name:     "Spiked helm of Necromancy",
	UniqueId: 6,
	Mod: Mod{
		Add: BaseStats{
			Attack:  1,
			Defense: 3,
			Magic:   1,
			Resist:  5,
		},
	},
}
