package combat

import (
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
)

type OddmentTable struct {
	Items   []DropChanceItem
	Oddment float64
}

func OddmentTableCreate(items []DropChanceItem) *OddmentTable {
	o := &OddmentTable{
		Items: items,
	}
	o.Oddment = o.CalcOddment()
	return o
}

func (o OddmentTable) CalcOddment() (total float64) {
	for _, v := range o.Items {
		total += v.Oddment
	}
	return total
}

func (o OddmentTable) Pick() world.ItemIndex {
	n := utilz.RandFloat(0, o.Oddment)
	var total float64
	for _, v := range o.Items {
		total += v.Oddment
		if total >= n {
			return world.ItemIndex{
				Id:    v.ItemId,
				Count: 1,
			}
		}
	}

	//Otherwise return the last item
	last := o.Items[len(o.Items)-1]
	return world.ItemIndex{
		Id:    last.ItemId,
		Count: 1,
	}
}
