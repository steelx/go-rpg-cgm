package world

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"log"
	"math"
)

type World struct {
	Time, Gold      float64
	Items, KeyItems []ItemIndex
	//Party check world_extended.go
	Icons Icons
}

type ItemIndex struct {
	Id, Count int
}

func Create() *World {
	w := &World{
		Time:     0,
		Gold:     0,
		Items:    make([]ItemIndex, 0),
		KeyItems: make([]ItemIndex, 0),
		Icons:    IconsCreate(),
	}

	//temp user items in inventory
	//w.Items = append(w.Items, ItemIndex{Id: 1, Count: 2})
	//w.Items = append(w.Items, ItemIndex{Id: 2, Count: 1})
	//w.Items = append(w.Items, ItemIndex{Id: 3, Count: 1})
	//w.KeyItems = append(w.KeyItems, ItemIndex{Id: 4, Count: 1})

	return w
}

func (w *World) AddItem(itemId, count int) {
	if _, ok := ItemsDB[itemId]; !ok {
		log.Fatal(fmt.Sprintf("Item ID {%v} does not exists in DB", itemId))
	}

	for i := range w.Items {
		//Does it already exist in World
		if w.Items[i].Id == itemId {
			w.Items[i].Count += count
			return
		}
	}

	//Add new
	w.Items = append(w.Items, ItemIndex{
		Id:    itemId,
		Count: count,
	})
}

func (w *World) RemoveItem(itemId, count int) {
	if _, ok := ItemsDB[itemId]; !ok {
		log.Fatal(fmt.Sprintf("Item ID {%v} does not exists in DB", itemId))
	}

	for i := len(w.Items) - 1; i <= 0; i-- {
		//Does it already exist in World
		if w.Items[i].Id == itemId {
			w.Items[i].Count -= count
		}

		if w.Items[i].Count <= 0 {
			w.removeItemFromArray(i)
		}
	}
}

func (w *World) removeItemFromArray(index int) {
	if len(w.Items) == 1 {
		w.Items = make([]ItemIndex, 0)
		return
	}
	w.Items[index], w.Items[0] = w.Items[0], w.Items[index]
	w.Items = w.Items[1 : len(w.Items)-1]
}

func (w World) hasKeyItem(itemId int) bool {
	for _, v := range w.KeyItems {
		if v.Id == itemId {
			return true
		}
	}
	return false
}

func (w *World) AddKeyItem(itemId int) {
	if w.hasKeyItem(itemId) {
		//if already exists we dont add again
		return
	}

	w.KeyItems = append(w.KeyItems, ItemIndex{Id: itemId, Count: 1})
}
func (w *World) RemoveKeyItem(itemId int) {
	if !w.hasKeyItem(itemId) {
		return
	}

	w.removeKeyItemFromArray(itemId)
}
func (w *World) removeKeyItemFromArray(index int) {
	if len(w.KeyItems) == 1 {
		w.Items = make([]ItemIndex, 0)
		return
	}
	w.KeyItems[index], w.KeyItems[0] = w.KeyItems[0], w.KeyItems[index]
	w.KeyItems = w.KeyItems[1 : len(w.KeyItems)-1]
}

func (w *World) Update(dt float64) {
	w.Time = w.Time + dt
}

func (w World) TimeAsString() string {
	time := w.Time
	hours := math.Floor(time / 3600)
	minutes := math.Ceil(math.Mod(time, 3600)/60) - 1
	seconds := int(time) % 60
	return fmt.Sprintf("%v:%v:%v", hours, minutes, seconds)
}

func (w World) GoldAsString() string {
	return fmt.Sprintf("%v", w.Gold)
}

func (w World) GetItemsAsStrings() []string {
	var items []string
	for _, item := range w.Items {
		items = append(items, fmt.Sprintf("%s, (%v)", ItemsDB[item.Id].Name, item.Count))
	}
	return items
}

func (w World) GetKeyItemsAsStrings() []string {
	var items []string
	for _, item := range w.KeyItems {
		items = append(items, fmt.Sprintf("%s, (%v)", ItemsDB[item.Id].Name, item.Count))
	}
	return items
}

//pending: use inside SelectionMenu renderItem pending
func (w World) DrawItem(renderer pixel.Target, x, y float64, itemIdx ItemIndex) {
	itemDef := ItemsDB[itemIdx.Id]

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	pos1 := pixel.V(x+40, y+(18/2))
	pos2 := pixel.V(x+40+18, y)
	textBase := text.New(pos2, basicAtlas)
	fmt.Fprintln(textBase, fmt.Sprintf("%-6s (%v)", itemDef.Name, itemIdx.Count))
	textBase.Draw(renderer, pixel.IM)

	iconSprite := w.Icons.Get(itemIdx.Id)
	iconSprite.Draw(renderer, pixel.IM.Moved(pos1))
}

func (w *World) HasKey(id int) bool {
	for _, v := range w.KeyItems {
		if v.Id == id {
			return true
		}
	}
	return false
}
