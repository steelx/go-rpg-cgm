package combat

import (
	"github.com/fatih/structs"
	"reflect"
)

/*
example: https://goplay.space/#xJX_ZyzORdZ

heroStats := combat.BaseStats{
	HpNow           : 300
	HpMax        : 300
	MpNow           : 300
	MpMax        : 300
	Strength: 10, Speed: 10, Intelligence: 10,
}

m := combat.StatsCreate(heroStats)
fmt.Println("Base Strength", m.GetBaseStat("Strength")) //10

magicSword := combat.Modifier{
	UniqueId: 1,
	Mod: combat.Mod{
		Add: combat.BaseStats{
			Strength: 5,
			Speed: 5,
		},
		Mult: combat.BaseStats{
			Strength: 2,
		},
	},
}
m.AddModifier(magicSword.UniqueId, magicSword.Mod)

fmt.Println("with Mod", m.Get("Strength")) //45 ==> 10 + 5 + (15*2)
*/

type BaseStats struct {
	HpNow, HpMax                   float64
	MpNow, MpMax                   float64
	Strength, Speed, Intelligence  float64 //Hero Stats
	Attack, Defense, Magic, Resist float64 //Equipment Stats
}

type Modifier struct {
	Name     string
	UniqueId int
	Mod      Mod
}
type Mod struct {
	Add  BaseStats
	Mult BaseStats
}

type Stats struct {
	Base      map[string]float64
	Modifiers map[int]Mod
}

func StatsCreate(stats BaseStats) Stats {
	s := Stats{
		Base:      make(map[string]float64),
		Modifiers: make(map[int]Mod),
	}

	//Iterate through the fields of a struct
	v := reflect.ValueOf(stats)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i).Interface().(float64)
		s.Base[typeOfS.Field(i).Name] = val
	}
	return s
}

//AddModifier
/*
magic_sword := Modifier{
		UniqueId: 1,
		Mod: Mod{
			Add:  BaseStats{
				Strength: 5,
				Speed: 5,
			},
		},
	}
*/
func (s *Stats) AddModifier(uniqueId int, modifier Mod) {
	s.Modifiers[uniqueId] = modifier
}
func (s *Stats) RemoveModifier(uniqueId int) {
	delete(s.Modifiers, uniqueId)
}

//Get id = BaseStats.KEY e.g. Get("Strength")
func (s Stats) Get(id string) float64 {
	total := s.Base[id] //10
	multiplier := 0.0

	for _, modifier := range s.Modifiers {
		add := structs.Map(modifier.Add)
		addVal := reflect.ValueOf(add[id])    //e.g. modifier.Add.Strength if id = Strength
		total += addVal.Interface().(float64) //+ 5

		mult := structs.Map(modifier.Mult)
		multVal := reflect.ValueOf(mult[id])        //e.g. modifier.Mult.Strength if id = Strength
		multiplier += multVal.Interface().(float64) //+ 2
	}

	return total + (total * multiplier) //15 + (15*2) == 45
}

func (s Stats) GetBaseStat(id string) float64 {
	return s.Base[id]
}

//Set e.g. Set("HpNow", 50)
//In combat, the HpNow and MpNow stats often change
func (s *Stats) Set(baseStatId string, val float64) {
	s.Base[baseStatId] = val
}
