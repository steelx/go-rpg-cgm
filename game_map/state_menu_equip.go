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
	inInventoryList               bool
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

	slotMenu := gui.SelectionMenuCreate(26, 80,
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
	var topMargin, leftMargin float64 = 25, 20
	basicAtlasAscii := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Title
	titleX := e.Layout.MidX("title")
	titleY := e.Layout.MidY("title")
	pos := pixel.V(titleX, titleY)
	textBase := text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, frontMenuOrder[equip])
	textBase.Draw(renderer, pixel.IM)

	// Char summary - Top Left
	_, titleHeight := e.Layout.Panels["title"].GetSize()
	avatarX := e.Layout.Left("top")
	avatarY := e.Layout.Top("top")
	avatarY = avatarY - titleHeight - 10
	e.actorSummary.SetPosition(avatarX, avatarY)
	e.actorSummary.Render(renderer)

	// Slots selection - Top Right
	equipX := e.Layout.MidX("top") + leftMargin
	equipY := e.Layout.Top("top") - titleHeight - leftMargin
	e.SlotMenu.SetPosition(equipX, equipY)
	e.SlotMenu.Render(renderer)

	// Inventory list - Bottom Right
	listX := e.Layout.Left("list") + 6
	listY := e.Layout.Top("list") - topMargin
	menu := e.FilterMenus[e.menuIndex]
	menu.SetPosition(listX, listY)
	menu.Render(renderer)

	// Char stat panel - Bottom Left
	slot := e.GetSelectedSlot() //Accessory2
	itemId := e.GetSelectedItem()

	item := world.ItemsDB[itemId]
	diffs := e.actorSummary.Actor.PredictStats(slot, item)
	x := e.Layout.Left("stats") + leftMargin
	y := e.Layout.Top("stats") - topMargin

	statList := e.actorSummary.Actor.CreateStatNameList()
	statLabels := e.actorSummary.Actor.CreateStatLabelList()
	for k, v := range statList {
		e.DrawStat(renderer, x, y, statLabels[k], v, diffs[v])
		y = y - 15
	}

	// Description panel
	descX := e.Layout.Left("desc") + leftMargin
	descY := e.Layout.MidY("desc") - 5
	pos = pixel.V(descX, descY)
	textBase = text.New(pos, basicAtlasAscii)
	fmt.Fprintln(textBase, item.Description)
	textBase.Draw(renderer, pixel.IM)
}

func (e EquipMenuState) Exit() {
}

func (e *EquipMenuState) Update(dt float64) {
	if e.inInventoryList {
		menu := e.FilterMenus[e.menuIndex]
		if !menu.IsDataSourceEmpty() {
			menu.ShowCursor()
			menu.HandleInput(e.win)
		}
		if e.win.JustReleased(pixelgl.KeyEscape) {
			e.FocusSlotMenu()
		}

	} else {
		e.SlotMenu.HandleInput(e.win)
		e.OnEquipMenuChanged()

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

	for i, slot := range e.actorSummary.Actor.ActiveEquipSlots {
		slotType := e.actorSummary.Actor.GetItemTypeBySlotPos(slot)

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

	e.FilterMenus = make([]*gui.SelectionMenu, len(filterList))
	for index, f := range filterList {
		menu := gui.SelectionMenuCreate(20, 60,
			f.list,
			false,
			pixel.V(0, 0),
			e.OnDoEquip,
			e.parent.World.DrawItem,
		)
		e.FilterMenus[index] = &menu
	}

}

func (e *EquipMenuState) OnDoEquip(index int, itemIdxI interface{}) {
	itemIdxV := reflect.ValueOf(itemIdxI)
	itemIdx := itemIdxV.Interface().(world.ItemIndex)
	item := e.parent.World.Get(itemIdx)

	//equipSlotId := e.actorSummary.Actor.GetEquipSlotIdByItemType(item.ItemType)
	e.actorSummary.Actor.Equip(e.GetSelectedSlot(), item)

	e.RefreshFilteredMenus()
	e.FocusSlotMenu()
}

//FocusSlotMenu - show Cursor on Item Slots (Top Right)
func (e *EquipMenuState) FocusSlotMenu() {
	e.inInventoryList = false
	e.SlotMenu.ShowCursor()
	e.menuIndex = e.SlotMenu.GetIndex()
	e.FilterMenus[e.menuIndex].HideCursor()
}

//OnSelectMenu get trigger when user selects a Item Slot &
// then cursor should be visible in Inventory List (Bottom Right)
func (e *EquipMenuState) OnSelectMenu(i int, wItemTypeI interface{}) {
	e.inInventoryList = true
	e.SlotMenu.HideCursor()
	//e.menuIndex = e.SlotMenu.GetIndex()
	//e.FilterMenus[e.menuIndex].ShowCursor()
}

func (e *EquipMenuState) OnEquipMenuChanged() {
	e.menuIndex = e.SlotMenu.GetIndex()
	e.FilterMenus[e.menuIndex].HideCursor()
}

//GetSelectedSlot takes index e.g. 3, returns "Accessory2"
func (e EquipMenuState) GetSelectedSlot() string {
	i := e.SlotMenu.GetIndex()
	return combat.ActorLabels.EquipSlotId[i]
}

func (e *EquipMenuState) GetSelectedItem() int {
	if e.inInventoryList {
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
	pos := pixel.V(x, y)
	basicAtlasAscii := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textBase := text.New(pos, basicAtlasAscii)

	current := e.actorSummary.Actor.Stats.Get(statId)
	changed := current + diff

	var textFormatted string
	if changed == current {
		textFormatted = fmt.Sprintf(`%-14s: %-4v`, label, current)
	} else {
		textFormatted = fmt.Sprintf(`%-14s: %-4v -> %+4v`, label, current, changed)
	}
	fmt.Fprintln(textBase, textFormatted)
	textBase.Draw(renderer, pixel.IM)

	textWidth := textBase.BoundsOf(textFormatted).W() + 20
	pos = pixel.V(x+textWidth, y+4)
	if diff > 0 {
		e.betterStatsIcon.Draw(renderer, pixel.IM.Moved(pos))
	} else if diff < 0 {
		e.badStatsIcon.Draw(renderer, pixel.IM.Moved(pos))
	}

}
