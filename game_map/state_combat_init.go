package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/utilz"
	"image/color"
	"reflect"
)

//CS -> CombatState
const (
	csNpcStand = "cs_npc_stand"
	csStandby  = "cs_standby"  // The character is waiting to be told what action to do by the player or AI
	csProne    = "cs_prone"    // The character is waiting and ready to perform a given combat action
	csAttack   = "cs_attack"   // The character will run an attack animation and attack an enemy
	csCast     = "cs_cast"     // The character will run a cast-spell animation and a special effect will play
	csUse      = "cs_use"      // The character uses some item with a use-item animation
	csHurt     = "cd_hurt"     // The character takes some damage. Animation and numbers
	csDie      = "cs_die"      // The character dies and the sprite is changed to the death sprite
	csMove     = "cs_move"     // The character moves toward or away from the enemy, in order to perform an action
	csVictory  = "cs_victory"  // The character dances around and combat ends
	csRunanim  = "cs_run_anim" // plays common animations states
)

type CombatState struct {
	GameState     *gui.StateStack
	InternalStack *gui.StateStack
	win           *pixelgl.Window
	Background    *pixel.Sprite
	Pos           pixel.Vec
	Layout        map[string][][]pixel.Vec
	Actors        map[string][]*combat.Actor
	Characters    map[string][]*Character
	ActorCharMap  map[*combat.Actor]*Character
	SelectedActor *combat.Actor

	Panels []gui.Panel
	TipPanel,
	NoticePanel gui.Panel
	PanelTitles []PanelTitle
	PartyList,
	StatsList *gui.SelectionMenu
	//HP and MP columns in the bottom right panel
	StatsYCol,
	marginLeft,
	marginTop float64
	Bars       map[string]BarStats //actor ID = BarStats
	imd        *imdraw.IMDraw
	EventQueue *EventQueue
}

type PanelTitle struct {
	text string
	x, y float64
}
type BarStats struct {
	HP, MP gui.ProgressBarIMD
}

func CombatStateCreate(state *gui.StateStack, win *pixelgl.Window, def CombatDef) *CombatState {
	backgroundImg, err := utilz.LoadPicture(def.Background)
	utilz.PanicIfErr(err)

	// Setup layout panel
	pos := pixel.V(0, 0)
	layout := gui.LayoutCreate(pos.X, pos.Y, win)
	layout.SplitHorz("screen", "top", "bottom", 0.72, 0)
	layout.SplitHorz("top", "notice", "top", 0.25, 0)
	layout.Contract("notice", 75, 25)
	layout.SplitHorz("bottom", "tip", "bottom", 0.24, 0)
	layout.SplitVert("bottom", "left", "right", 0.67, 0)

	c := &CombatState{
		win:           win,
		GameState:     state,
		InternalStack: gui.StateStackCreate(win),
		Background:    pixel.NewSprite(backgroundImg, backgroundImg.Bounds()),
		Pos:           pos,
		Actors: map[string][]*combat.Actor{
			party:   def.Actors.Party,
			enemies: def.Actors.Enemies,
		},
		Characters:   make(map[string][]*Character),
		ActorCharMap: make(map[*combat.Actor]*Character),
		StatsYCol:    208,
		marginLeft:   18,
		marginTop:    20,
		imd:          imdraw.New(nil),
		EventQueue:   EventsQueueCreate(),
	}

	c.Layout = combatLayout
	c.CreateCombatCharacters(party)
	c.CreateCombatCharacters(enemies)

	c.Panels = []gui.Panel{
		layout.CreatePanel("left"),
		layout.CreatePanel("right"),
	}
	c.TipPanel = layout.CreatePanel("tip")
	c.NoticePanel = layout.CreatePanel("notice")

	//Set up player list
	partyListMenu := gui.SelectionMenuCreate(19, 0, 100,
		c.Actors[party],
		false,
		pixel.ZV,
		c.OnPartyMemberSelect,
		c.RenderPartyNames,
	)
	c.PartyList = &partyListMenu

	//title
	x := layout.Left("left")
	y := layout.Top("left")

	marginTop := c.marginTop
	marginLeft := c.marginLeft
	c.PanelTitles = []PanelTitle{
		{"NAME", x + marginLeft, y - marginTop + 2},
		{"HP", layout.Left("right") + marginLeft, y - marginTop + 2},
		{"MP", layout.Left("right") + marginLeft + c.StatsYCol, y - marginTop + 2},
	}

	y = y - 35 // - margin top
	c.PartyList.SetPosition(x+marginLeft, y)
	c.PartyList.HideCursor()

	c.Bars = make(map[string]BarStats)
	for _, v := range c.Actors[party] {

		hpBar := gui.ProgressBarIMDCreate(
			0, 0,
			v.Stats.Get("HpNow"),
			v.Stats.Get("HpMax"),
			"#FF001E",
			"#15FF00",
			3, 100,
			c.imd,
		)
		mpBar := gui.ProgressBarIMDCreate(
			0, 0,
			v.Stats.Get("MpNow"),
			v.Stats.Get("MpMax"),
			"#A48B2C",
			"#00E7DA",
			3, 100,
			c.imd,
		)

		c.Bars[v.Id] = BarStats{
			HP: hpBar,
			MP: mpBar,
		}
	}

	statsListMenu := gui.SelectionMenuCreate(19, 0, 100,
		c.Actors[party],
		false,
		pixel.ZV,
		c.OnPartyMemberSelect,
		c.RenderPartyStats,
	)
	c.StatsList = &statsListMenu

	x = layout.Left("right") - 8
	c.StatsList.SetPosition(x, y)
	c.StatsList.HideCursor()

	return c
}

func (c *CombatState) Enter() {
}

func (c *CombatState) Exit() {
}

func (c *CombatState) Update(dt float64) bool {
	for _, v := range c.Characters[party] {
		v.Controller.Update(dt)
	}
	for _, v := range c.Characters[enemies] {
		v.Controller.Update(dt)
	}

	if len(c.InternalStack.States) != 0 && c.InternalStack.Top() != nil {
		c.InternalStack.Update(dt)
		return true
	}

	c.EventQueue.Update()
	c.AddTurns(c.Actors[party])
	c.AddTurns(c.Actors[enemies])

	if c.PartyWins() {
		c.EventQueue.Clear()
		//deal with win
	} else if c.EnemyWins() {
		c.EventQueue.Clear()
		//deal with lost
	}

	return true
}

func (c CombatState) Render(renderer *pixelgl.Window) {
	c.Background.Draw(renderer, pixel.IM.Scaled(c.Pos, 1).Moved(c.Pos))

	for _, v := range c.Characters[party] {
		pos := pixel.V(v.Entity.X, v.Entity.Y)
		v.Entity.Render(nil, renderer, pos)
	}
	for _, v := range c.Characters[enemies] {
		pos := pixel.V(v.Entity.X, v.Entity.Y)
		v.Entity.Render(nil, renderer, pos)
	}

	for _, v := range c.Panels {
		v.Draw(renderer)
	}
	//c.TipPanel.Draw(renderer)
	//c.NoticePanel.Draw(renderer) //pending

	textBase := text.New(pixel.V(0, 0), gui.BasicAtlas12)
	//textBase.Color = txtColor
	for _, v := range c.PanelTitles {
		textBase.Clear()
		fmt.Fprintln(textBase, v.text)
		textBase.Draw(renderer, pixel.IM.Moved(pixel.V(v.x, v.y)))
	}

	c.PartyList.Render(renderer)
	c.StatsList.Render(renderer)

	c.InternalStack.Render(renderer)
	c.EventQueue.Render(renderer)

	camera := pixel.IM.Scaled(c.Pos, 1.0).Moved(c.win.Bounds().Center().Sub(c.Pos))
	c.win.SetMatrix(camera)
}

func (c *CombatState) HandleInput(win *pixelgl.Window) {
}

func (c *CombatState) CreateCombatCharacters(key string) {
	actorsList := c.Actors[key]
	layout := c.Layout[key][len(actorsList)-1]

	for k, v := range actorsList {
		charDef, ok := CharacterDefinitions[v.Id]
		if !ok {
			panic(fmt.Sprintf("Id '%s' Not found in CharacterDefinitions", v.Id))
		}

		if charDef.CombatEntityDef.Texture != "" {
			charDef.EntityDef = charDef.CombatEntityDef
		}

		var char *Character
		char = CharacterCreate(
			charDef,
			map[string]func() state_machine.State{
				csStandby: func() state_machine.State {
					return CSStandByCreate(char, c)
				},
				csNpcStand: func() state_machine.State {
					return NPCStandCombatStateCreate(char, c)
				},
				csRunanim: func() state_machine.State {
					return CSRunAnimCreate(char, c)
				},
				csHurt: func() state_machine.State {
					return CSHurtCreate(char, c)
				},
				csMove: func() state_machine.State {
					return CSMoveCreate(char, c)
				},
			},
		)

		c.ActorCharMap[v] = char

		pos := layout[k]

		// Combat positions are 0 - 1
		// Need scaling to the screen size.
		x := pos.X * c.win.Bounds().W()
		y := pos.Y * c.win.Bounds().H()

		char.Entity.X = x
		char.Entity.Y = y

		// Change to standby because it's combat time
		animName := csStandby
		char.Controller.Change(csStandby, animName)

		c.Characters[key] = append(c.Characters[key], char)
	}

}

func (c *CombatState) OnPartyMemberSelect(index int, str interface{}) {
	fmt.Println(index, str)
}

func (c *CombatState) RenderPartyNames(args ...interface{}) {
	//renderer pixel.Target, x, y float64, index int
	rendererV := reflect.ValueOf(args[0])
	renderer := rendererV.Interface().(pixel.Target)
	xV := reflect.ValueOf(args[1])
	yV := reflect.ValueOf(args[2])
	x, y := xV.Interface().(float64), yV.Interface().(float64)

	itemV := reflect.ValueOf(args[3])
	actor := itemV.Interface().(*combat.Actor)

	var txtColor color.RGBA
	if c.SelectedActor == actor {
		txtColor = utilz.HexToColor("#00cc00") //green
	} else {
		txtColor = utilz.HexToColor("#addbd8") //light blue
	}

	textBase := text.New(pixel.V(x, y), gui.BasicAtlasAscii)
	textBase.Color = txtColor
	fmt.Fprintln(textBase, actor.Name)
	textBase.Draw(renderer, pixel.IM)

}

func (c *CombatState) RenderPartyStats(args ...interface{}) {
	//renderer pixel.Target, x, y float64, index int
	rendererV := reflect.ValueOf(args[0])
	renderer := rendererV.Interface().(pixel.Target)
	xV := reflect.ValueOf(args[1])
	yV := reflect.ValueOf(args[2])
	x, y := xV.Interface().(float64), yV.Interface().(float64)

	x = x + c.marginLeft + 10
	itemV := reflect.ValueOf(args[3])
	actor := itemV.Interface().(*combat.Actor)

	stats := actor.Stats
	barOffset := 70.0

	bars := c.Bars[actor.Id]
	bars.HP.SetPosition(x+barOffset, y)
	bars.HP.SetValue(stats.Get("HpNow"))
	bars.HP.Render(renderer)

	c.DrawHP(renderer, x, y, actor)

	x = x + c.StatsYCol
	c.DrawMP(renderer, x, y, actor)

	mpNow := stats.Get("MpNow")
	bars.MP.SetPosition(x+barOffset*0.7, y)
	bars.MP.SetValue(mpNow)
	bars.MP.Render(renderer)
}

func (c *CombatState) DrawHP(renderer pixel.Target, x, y float64, actor *combat.Actor) {
	hp, max := actor.Stats.Get("HpNow"), actor.Stats.Get("HpMax")
	percent := hp / max

	txtColor := utilz.HexToColor("#ffffff")
	if percent < 0.2 {
		txtColor = utilz.HexToColor("#00cc00") //green
	} else if percent < 0.45 {
		txtColor = utilz.HexToColor("#ffffa2") //light yellow
	}

	textBase := text.New(pixel.V(x, y), gui.BasicAtlasAscii)
	textBase.Color = txtColor
	fmt.Fprintf(textBase, fmt.Sprintf("%v/%v", hp, max))
	textBase.Draw(renderer, pixel.IM)
}

func (c *CombatState) DrawMP(renderer pixel.Target, x, y float64, actor *combat.Actor) {
	mpNow := actor.Stats.Get("MpNow")
	mpNowStr := fmt.Sprintf("%v", mpNow)
	textBase := text.New(pixel.V(x, y), gui.BasicAtlasAscii)
	fmt.Fprintln(textBase, mpNowStr)
	textBase.Draw(renderer, pixel.IM)
}
