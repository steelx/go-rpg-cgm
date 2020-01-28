package combat

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/utilz"
	"github.com/steelx/go-rpg-cgm/world"
	"golang.org/x/image/font/basicfont"
)

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
	ActiveEquipSlots []int
	Equipment        map[string]int
	worldRef         *WorldExtended
}

/* example: ActorCreate(HeroDef)
var HeroDef = combat.ActorDef{
		Stats: combat.DefaultStats,
		StatGrowth: map[string]func() int{
			"HpMax":        dice.Create("4d50+100"),
			"MpMax":        dice.Create("2d50+100"),
			"Strength":     combat.StatsGrowth.Fast,
			"Speed":        combat.StatsGrowth.Fast,
			"Intelligence": combat.StatsGrowth.Med,
		},
	}
*/
// ActorCreate
func ActorCreate(def ActorDef) Actor {
	actorAvatar, err := utilz.LoadPicture(def.Portrait)
	utilz.PanicIfErr(err)

	a := Actor{
		Id:               def.Id,
		Name:             def.Name,
		StatGrowth:       def.StatGrowth,
		Stats:            world.StatsCreate(def.Stats),
		XP:               0,
		Level:            1,
		PortraitTexture:  actorAvatar,
		Portrait:         pixel.NewSprite(actorAvatar, actorAvatar.Bounds()),
		Actions:          def.Actions,
		ActiveEquipSlots: def.ActiveEquipSlots,
		Equipment: map[string]int{
			"Weapon":  def.Weapon,
			"Armor":   def.Armor,
			"Access1": def.Access1,
			"Access2": def.Access2,
		},
	}

	a.NextLevelXP = NextLevel(a.Level)
	return a
}

func (a *Actor) RenderEquipment(renderer pixel.Target, x, y float64, index int) {
	label := ActorLabels.EquipSlotLabels[index]

	equipmentText := "none"
	if index < len(a.Equipment) {
		slotId := ActorLabels.EquipSlotId[index]
		itemId := a.Equipment[slotId]
		item := world.ItemsDB[itemId]
		equipmentText = item.Name
	}

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

func (a *Actor) Equip(equipSlotId string, item world.Item) {
	prevItemId, ok := a.Equipment[equipSlotId]
	if ok {
		delete(a.Equipment, equipSlotId)
		a.Stats.RemoveModifier(prevItemId)
		a.worldRef.AddItem(prevItemId, 1) //return back to World
	}

	//UnEquip
	if item.Id == -1 {
		return
	}

	a.worldRef.RemoveItem(item.Id, 1) //remove from World
	a.Equipment[equipSlotId] = item.Id

	modifier := item.Stats
	a.Stats.AddModifier(item.Id, modifier)
}

func (a *Actor) UnEquip(equipSlotId string) {
	a.Equip(equipSlotId, world.Item{Id: -1})
}

func (a Actor) createStatNameList() (statsIDs []string) {
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

// PredictStats
//compare Equipped Item Stats with given Item
// returns -> BaseStats after comparison
func (a Actor) PredictStats(equipSlotId string, item world.Item) map[string]float64 {
	statsIDs := a.createStatNameList()

	currentStats := make(map[string]float64)
	for _, key := range statsIDs {
		currentStats[key] = a.Stats.Get(key)
	}

	// Replace item
	prevItemId, ok := a.Equipment[equipSlotId]
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
