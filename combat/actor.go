package combat

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
	"golang.org/x/image/font/basicfont"
	"reflect"
)

type ActorDropItem struct {
	XP     float64
	Gold   float64
	Always []int //ActionItem ids that are guaranteed to drop
	Chance *OddmentTable
}

// Actor is any creature or character that participates in combat
// and therefore requires stats, equipment, etc
type Actor struct {
	Id, Name   string
	Stats      world.Stats
	StatGrowth map[string]func() int

	PortraitTexture  pixel.Picture
	Portrait         *pixel.Sprite
	Level            int
	XP, NextLevelXP  float64
	Actions          []string
	Magic            []string
	Special          []string
	StealItem        int //Item ID only for Enemy actors
	ActiveEquipSlots []int
	Equipped         map[string]int //int is ItemsDB Id
	worldRef         *WorldExtended
	isPlayer         bool
	Drop             ActorDropItem
}

// ActorCreate
func ActorCreate(def ActorDef, randName ...interface{}) Actor {
	randNameV := ""
	if len(randName) > 0 {
		randNameV = reflect.ValueOf(randName[0]).Interface().(string)
	}
	actorAvatar, err := utilz.LoadPicture(def.Portrait)
	utilz.PanicIfErr(err)

	a := Actor{
		Id:               def.Id,
		isPlayer:         def.IsPlayer,
		Name:             fmt.Sprintf("%s%s", def.Name, randNameV),
		StatGrowth:       def.StatGrowth,
		Stats:            world.StatsCreate(def.Stats),
		XP:               0,
		Level:            1,
		PortraitTexture:  actorAvatar,
		Portrait:         pixel.NewSprite(actorAvatar, actorAvatar.Bounds()),
		Actions:          def.Actions,
		Magic:            def.Magic,
		Special:          def.Special,
		StealItem:        def.StealItem,
		ActiveEquipSlots: def.ActiveEquipSlots,
		Equipped: map[string]int{
			ActorLabels.EquipSlotId[0]: def.Weapon,
			ActorLabels.EquipSlotId[1]: def.Armor,
			ActorLabels.EquipSlotId[2]: def.Access1,
			ActorLabels.EquipSlotId[3]: def.Access2,
		},
	}

	if !def.IsPlayer {
		gold := utilz.RandInt(def.Drop.Gold[0], def.Drop.Gold[1])
		a.Drop.XP = def.Drop.XP
		a.Drop.Gold = float64(gold)
		a.Drop.Chance = OddmentTableCreate(def.Drop.Chance)
	}

	a.NextLevelXP = NextLevel(a.Level)
	return a
}

func (a *Actor) RenderEquipment(args ...interface{}) {
	//renderer pixel.Target, x, y float64, index int
	rendererV := reflect.ValueOf(args[0])
	renderer := rendererV.Interface().(pixel.Target)
	xV := reflect.ValueOf(args[1])
	x := xV.Interface().(float64)
	yV := reflect.ValueOf(args[2])
	y := yV.Interface().(float64)
	itemV := reflect.ValueOf(args[3])
	slot := itemV.Interface().(int)

	label := ActorLabels.EquipSlotId[slot]

	itemId := a.Equipped[label]
	item := world.ItemsDB[itemId]
	equipmentText := item.Name

	basicAtlasAscii := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	pos := pixel.V(x, y)
	textBase := text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, fmt.Sprintf("%-12s: %v", label, equipmentText))
	textBase.Draw(renderer, pixel.IM)
}

func (a Actor) ReadyToLevelUp() bool {
	return a.XP >= a.NextLevelXP
}

func (a *Actor) AddXP(xp float64) bool {
	a.XP += xp
	return a.ReadyToLevelUp()
}

func (a Actor) CreateLevelUp() LevelUp {
	levelUp := LevelUp{
		XP:        -a.NextLevelXP,
		Level:     1,
		BaseStats: make(map[string]float64),
	}

	for id, diceRoll := range a.StatGrowth {
		levelUp.BaseStats[id] = float64(diceRoll())
	}

	//Pending feature
	// Additional level up code
	// e.g. if you want to apply
	// a bonus every 4 levels
	// or heal the players MP/HP

	return levelUp
}

func (a *Actor) ApplyLevel(levelUp LevelUp) {
	a.XP += levelUp.XP
	a.Level += levelUp.Level
	a.NextLevelXP = NextLevel(a.Level)

	for k, v := range levelUp.BaseStats {
		a.Stats.Base[k] += v
	}

	//Pending feature
	// Unlock any special abilities etc.
}

//GetEquipSlotIdByItemType takes in Item Type INT return Type string e.g. Weapon
func (a Actor) GetEquipSlotIdByItemType(itemT world.ItemType) string {
	return ActorLabels.EquipSlotTypes[itemT]
}

func (a Actor) GetItemTypeBySlotPos(slot int) world.ItemType {
	switch slot {
	case 0:
		return world.Weapon
	case 1:
		return world.Armor
	case 2:
		return world.Accessory
	case 3:
		return world.Accessory
	}

	return 0
}

func (a *Actor) Equip(equipSlotId string, item world.Item) {
	prevItemId, ok := a.Equipped[equipSlotId]
	if ok && prevItemId != 0 {
		delete(a.Equipped, equipSlotId)
		a.Stats.RemoveModifier(prevItemId)
		a.worldRef.AddItem(prevItemId, 1) //return back to World
	}

	//UnEquip
	if item.Id == -1 {
		return
	}

	a.worldRef.RemoveItem(item.Id, 1) //remove from World
	a.Equipped[equipSlotId] = item.Id

	modifier := item.Stats
	a.Stats.AddModifier(item.Id, modifier)
}

func (a *Actor) UnEquip(equipSlotId string) {
	a.Equip(equipSlotId, world.Item{Id: -1})
}

func (a Actor) CreateStatNameList() (statsIDs []string) {
	for _, v := range ActorLabels.ActorStats {
		statsIDs = append(statsIDs, v)
	}

	for _, v := range ActorLabels.ItemStats {
		statsIDs = append(statsIDs, v)
	}

	statsIDs = append(statsIDs, "HpMax")
	statsIDs = append(statsIDs, "MpMax")

	return
}

func (a Actor) CreateStatLabelList() (statsLabels []string) {
	for _, v := range ActorLabels.ActorStatLabels {
		statsLabels = append(statsLabels, v)
	}

	for _, v := range ActorLabels.ItemStatLabels {
		statsLabels = append(statsLabels, v)
	}

	statsLabels = append(statsLabels, "HP:")
	statsLabels = append(statsLabels, "MP:")

	return
}

// PredictStats
//compare Equipped Item Stats with given Item
// returns -> BaseStats after comparison
func (a Actor) PredictStats(equipSlotId string, item world.Item) map[string]float64 {
	statsIDs := a.CreateStatNameList()

	currentStats := make(map[string]float64)
	for _, key := range statsIDs {
		currentStats[key] = a.Stats.Get(key)
	}

	// Replace item
	prevItemId, ok := a.Equipped[equipSlotId]
	if ok {
		a.Stats.RemoveModifier(prevItemId)
	}
	a.Stats.AddModifier(item.Id, item.Stats)

	// Get values for modified stats
	modifiedStats := make(map[string]float64)
	for _, key := range statsIDs {
		modifiedStats[key] = a.Stats.Get(key)
	}

	diffStats := make(map[string]float64)
	for _, key := range statsIDs {
		diffStats[key] = modifiedStats[key] - currentStats[key]
	}

	// Undo replace item
	a.Stats.RemoveModifier(item.Id)
	if ok {
		a.Stats.AddModifier(prevItemId, world.ItemsDB[prevItemId].Stats)
	}

	return diffStats
}

func (a Actor) CanUse(item world.Item) bool {
	if len(item.Restrictions) == 0 {
		return true
	}

	for _, v := range item.Restrictions {
		if v == a.Id {
			return true
		}
	}

	return false
}

//IsPlayer tell's if Actor is player controlled e.g. Hero, Mage, ..
func (a *Actor) IsPlayer() bool {
	return a.isPlayer
}

//has Knocked Out?
func (a Actor) IsKOed() bool {
	return a.Stats.Get("HpNow") <= 0
}

//Knock Out
func (a *Actor) KO() {
	//TODO : impl Actor KO
}
