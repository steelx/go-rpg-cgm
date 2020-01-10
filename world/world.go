package world

import (
	"fmt"
	"log"
	"math"
)

type World struct {
	Time, Gold      float64
	Items, KeyItems map[int]ItemIndex
}

type ItemIndex struct {
	id, count int
}

func WorldCreate() World {
	w := World{
		Time:     0,
		Gold:     0,
		Items:    make(map[int]ItemIndex, 0),
		KeyItems: make(map[int]ItemIndex, 0),
	}

	w.Items[2] = ItemIndex{id: 2, count: 1}
	return w
}

func (w *World) AddItem(itemId, count int) {
	if _, ok := ItemsDB[itemId]; !ok {
		log.Fatal(fmt.Sprintf("Item ID {%v} does not exists in DB", itemId))
	}

	//Does it already exist in World
	if item, ok := w.Items[itemId]; ok {
		item.count += count
		return
	}

	w.Items[itemId] = ItemIndex{
		id:    itemId,
		count: count,
	}
}

func (w *World) RemoveItem(itemId, count int) {
	if _, ok := ItemsDB[itemId]; !ok {
		log.Fatal(fmt.Sprintf("Item ID {%v} does not exists in DB", itemId))
	}

	//if it doesnt exists
	if _, ok := w.Items[itemId]; !ok {
		return
	}

	item := w.Items[itemId]
	if item.count-count <= 0 {
		delete(w.Items, itemId)
	} else {
		item.count = item.count + count
	}
}

func (w World) hasKeyItem(itemId int) bool {
	_, ok := w.KeyItems[itemId]
	return ok
}

func (w *World) AddKeyItem(itemId int) {
	if w.hasKeyItem(itemId) {
		return
	}

	w.KeyItems[itemId] = ItemIndex{id: itemId, count: 1}
}
func (w *World) RemoveKeyItem(itemId int) {
	if !w.hasKeyItem(itemId) {
		return
	}

	delete(w.KeyItems, itemId)
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
