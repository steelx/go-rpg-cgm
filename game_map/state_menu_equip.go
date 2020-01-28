package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/world"
	"golang.org/x/image/font/basicfont"
	"reflect"
)

type FilterList struct {
	slotType world.ItemType
	list     []world.ItemIndex
}

type EquipMenuState struct {
	parent *InGameMenuState
	win    *pixelgl.Window

	Panels                        []gui.Panel
	Layout                        gui.Layout
	betterStatsIcon, badStatsIcon *pixel.Sprite
	inList                        bool
	equipment                     map[string]int
	menuIndex                     int
	actorSummary                  gui.ActorSummary
	FilterMenus                   []*gui.SelectionMenu
	SlotMenu                      *gui.SelectionMenu
}

func EquipMenuStateCreate(parent *InGameMenuState, win *pixelgl.Window) *EquipMenuState {
	// Create panel layout
	layout := gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.12, 2)
	layout.SplitVert("top", "title", "category", 0.75, 2)
	titlePanel := layout.Panels["title"]

	layout = gui.LayoutCreate(0, 0, win)
	layout.Contract("screen", 120, 40)
	layout.SplitHorz("screen", "top", "bottom", 0.42, 2)
	layout.SplitHorz("bottom", "desc", "bottom", 0.2, 2)
	layout.SplitVert("bottom", "stats", "list", 0.6, 2)
	layout.Panels["title"] = titlePanel

	e := &EquipMenuState{
		win:             win,
		parent:          parent,
		betterStatsIcon: world.IconsDB.Get(10),
		badStatsIcon:    world.IconsDB.Get(11),
		Layout:          layout,
		Panels: []gui.Panel{
			layout.CreatePanel("top"),
			layout.CreatePanel("desc"),
			layout.CreatePanel("stats"),
			layout.CreatePanel("list"),
			layout.CreatePanel("title"),
		},
	}

	return e
}

func (e *EquipMenuState) Enter(actorSummaryI interface{}) {
	actorSummaryV := reflect.ValueOf(actorSummaryI)
	actorSummary := actorSummaryV.Interface().(gui.ActorSummary)
	e.actorSummary = actorSummary
	e.actorSummary.HideXP()
	e.equipment = actorSummary.Actor.Equipped

	e.RefreshFilteredMenus()
	e.menuIndex = 0

	slotMenu := gui.SelectionMenuCreate(26, 0,
		e.actorSummary.Actor.ActiveEquipSlots,
		false,
		pixel.V(0, 0),
		e.OnSelectMenu,
		e.actorSummary.Actor.RenderEquipment,
	)

	e.SlotMenu = &slotMenu

}

func (e EquipMenuState) Render(renderer *pixelgl.Window) {
	for _, v := range e.Panels {
		v.Draw(renderer)
	}

	basicAtlasAscii := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Title
	titleX := e.Layout.MidX("title")
	titleY := e.Layout.MidY("title")
	pos := pixel.V(titleX, titleY)
	textBase := text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, frontMenuOrder[equip])
	textBase.Draw(renderer, pixel.IM)

	// Char summary
	_, titleHeight := e.Layout.Panels["title"].GetSize()
	avatarX := e.Layout.Left("top")
	avatarY := e.Layout.Top("top")
	avatarX = avatarX + 10
	avatarY = avatarY - titleHeight - 10
	e.actorSummary.SetPosition(avatarX, avatarY)
	e.actorSummary.Render(renderer)

	// Slots selection
	equipX := e.Layout.MidX("top") + 40
	equipY := e.Layout.Top("top") - titleHeight - 20
	e.SlotMenu.SetPosition(equipX, equipY)
	e.SlotMenu.Render(renderer)

	// Inventory list
	listX := e.Layout.Left("list") + 6
	listY := e.Layout.Top("list") - 20
	menu := e.FilterMenus[e.menuIndex]
	menu.SetPosition(listX, listY)
	menu.Render(renderer)

	// Char stat panel
	slot := e.GetSelectedSlot()
	itemId := e.GetSelectedItem()

	item := world.ItemsDB[itemId]
	diffs := e.actorSummary.Actor.PredictStats(slot, item)
	x := e.Layout.MidX("stats") - 10
	y := e.Layout.Top("stats") - 14

	statList := e.actorSummary.Actor.CreateStatNameList()
	statLabels := e.actorSummary.Actor.CreateStatLabelList()
	for k, v := range statList {
		e.DrawStat(renderer, x, y, statLabels[k], v, diffs[v])
		y = y - 14
	}

	// Description panel
	descX := e.Layout.Left("desc") + 10
	descY := e.Layout.MidY("desc")
	pos = pixel.V(descX, descY)
	textBase = text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, item.Description)
	textBase.Draw(renderer, pixel.IM)
}

func (e EquipMenuState) Exit() {
}

func (e *EquipMenuState) Update(dt float64) {
	if e.inList {
		menu := e.FilterMenus[e.menuIndex]
		menu.HandleInput(e.win)
		if e.win.JustReleased(pixelgl.KeyEscape) {
			e.FocusSlotMenu()
		}

	} else {
		prevEquipIndex := e.SlotMenu.GetIndex()
		e.SlotMenu.HandleInput(e.win)
		if prevEquipIndex == e.SlotMenu.GetIndex() {
			e.OnEquipMenuChanged()
		}
		if e.win.JustPressed(pixelgl.KeyEscape) {
			e.parent.StateMachine.Change("frontmenu", nil)
			return
		}
	}
}

func (e *EquipMenuState) RefreshFilteredMenus() {
	// Get a list of filters by slot type
	// Items will be sorted into these lists
	slotCount := len(e.actorSummary.Actor.ActiveEquipSlots)
	filterList := make([]*FilterList, slotCount)

	for i := 0; i < slotCount; i++ {
		slotType := e.actorSummary.Actor.ActiveEquipSlots[i]
		filterList[i] = &FilterList{
			slotType: slotType,
			list:     make([]world.ItemIndex, 0),
		}
	}

	//Actually sort inventory items into array
	gWorld := e.parent.World
	for _, v := range gWorld.Items {
		item := world.ItemsDB[v.Id]

		for _, f := range filterList {
			if item.ItemType == f.slotType && e.actorSummary.Actor.CanUse(item) {
				f.list = append(f.list, v)
			}
		}
	}

	for _, f := range filterList {
		menu := gui.SelectionMenuCreate(26, 256,
			f.list,
			false,
			pixel.V(0, 0),
			e.OnDoEquip,
			e.parent.World.DrawItem,
		)
		e.FilterMenus = append(e.FilterMenus, &menu)
	}

}

func (e *EquipMenuState) OnDoEquip(index int, itemIdxI interface{}) {
	itemIdxV := reflect.ValueOf(itemIdxI)
	itemIdx := itemIdxV.Interface().(world.ItemIndex)
	item := e.parent.World.Get(itemIdx)
	fmt.Println("item", item)
	equipSlotId := e.actorSummary.Actor.GetEquipSlotIdByItemType(item.ItemType)
	e.actorSummary.Actor.Equip(equipSlotId, item)

	e.RefreshFilteredMenus()
	e.FocusSlotMenu()
}

func (e *EquipMenuState) FocusSlotMenu() {
	e.inList = false
	e.SlotMenu.ShowCursor()
	e.menuIndex = e.SlotMenu.GetIndex()
	e.FilterMenus[e.menuIndex].HideCursor()
}

func (e *EquipMenuState) OnSelectMenu(i int, wItemTypeI interface{}) {
	e.inList = true
	e.SlotMenu.HideCursor()
	e.menuIndex = e.SlotMenu.GetIndex()
	e.FilterMenus[e.menuIndex].ShowCursor()
}

func (e *EquipMenuState) OnEquipMenuChanged() {
	e.menuIndex = e.SlotMenu.GetIndex()
	e.FilterMenus[e.menuIndex].HideCursor()
}
func (e EquipMenuState) GetSelectedSlot() string {
	i := e.SlotMenu.GetIndex()
	return combat.ActorLabels.EquipSlotId[i]
}

func (e *EquipMenuState) GetSelectedItem() int {
	if e.inList {
		menu := e.FilterMenus[e.menuIndex]
		if menu.DataI != nil {
			itemIndexI := menu.SelectedItem()
			itemIndexV := reflect.ValueOf(itemIndexI)
			itemIndex := itemIndexV.Interface().(world.ItemIndex)
			return itemIndex.Id
		}
	}

	slot := e.GetSelectedSlot()
	return e.actorSummary.Actor.Equipped[slot]
}

func (e *EquipMenuState) DrawStat(renderer pixel.Target, x, y float64, label, statId string, diff float64) {
	basicAtlasAscii := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	pos := pixel.V(x, y)
	textBase := text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, label)
	textBase.Draw(renderer, pixel.IM)

	current := e.actorSummary.Actor.Stats.Get(statId)
	pos = pixel.V(x+15, y)
	textBase = text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, current)
	textBase.Draw(renderer, pixel.IM)

	changed := current + diff
	pos = pixel.V(x+60, y)
	textBase = text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, changed)
	textBase.Draw(renderer, pixel.IM)

	if diff > 0 {
		e.betterStatsIcon.Draw(renderer, pixel.IM.Moved(pixel.V(x+80, y)))
	} else if diff < 0 {
		e.badStatsIcon.Draw(renderer, pixel.IM.Moved(pixel.V(x+80, y)))
	}

}
