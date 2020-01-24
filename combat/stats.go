package combat

import (
	"fmt"
	"github.com/fatih/structs"
	"reflect"
)

/*
example:

m := combat.StatsCreate(combat.BaseStats{300, 300, 300, 300, 10, 10, 10})
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
	Hp           int
	HpMax        int
	Mp           int
	MpMax        int
	Strength     int
	Speed        int
	Intelligence int
}

type Modifier struct {
	UniqueId int
	Mod      Mod
}
type Mod struct {
	Add  BaseStats
	Mult BaseStats
}

type Stats struct {
	Base      map[string]int
	Modifiers map[int]Mod
}

func StatsCreate(stats BaseStats) Stats {
	s := Stats{
		Base:      make(map[string]int),
		Modifiers: make(map[int]Mod),
	}

	//Iterate through the fields of a struct
	v := reflect.ValueOf(stats)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i).Interface().(int)
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

//Get id = BaseStats.KEY
func (s Stats) Get(id string) int {
	total := s.Base[id] //10
	multiplier := 0

	for uniqueId, modifier := range s.Modifiers {
		fmt.Println("uniqueId", uniqueId)
		add := structs.Map(modifier.Add)
		addVal := reflect.ValueOf(add[id]) //e.g. modifier.Add.Strength if id = Strength
		total += addVal.Interface().(int)  //+ 5

		mult := structs.Map(modifier.Mult)
		multVal := reflect.ValueOf(mult[id])    //e.g. modifier.Mult.Strength if id = Strength
		multiplier += multVal.Interface().(int) //+ 2
	}

	return total + (total * multiplier) //15 + (15*2) == 45
}

func (s Stats) GetBaseStat(id string) int {
	return s.Base[id]
}
