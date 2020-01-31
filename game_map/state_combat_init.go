package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"github.com/steelx/go-rpg-cgm/utilz"
)

type CombatState struct {
	GameState     *gui.StateStack
	InternalState *gui.StateStack
	win           *pixelgl.Window
	Background    *pixel.Sprite
	Pos           pixel.Vec
	Layout        map[string][][]pixel.Vec
	Actors        map[string][]*combat.Actor
	Characters    map[string][]*Character
	ActorCharMap  map[string]*Character
	SelectedActor *combat.Actor
}

func CombatStateCreate(state *gui.StateStack, win *pixelgl.Window, def CombatDef) *CombatState {
	backgroundImg, err := utilz.LoadPicture(def.Background)
	utilz.PanicIfErr(err)

	c := &CombatState{
		win:           win,
		GameState:     state,
		InternalState: gui.StateStackCreate(win),
		Background:    pixel.NewSprite(backgroundImg, backgroundImg.Bounds()),
		Pos:           pixel.V(0, 0),
		Actors: map[string][]*combat.Actor{
			party:   def.Actors.Party,
			enemies: def.Actors.Enemies,
		},
		Characters:   make(map[string][]*Character),
		ActorCharMap: make(map[string]*Character),
	}

	c.Layout = combatLayout
	c.CreateCombatCharacters(party)
	c.CreateCombatCharacters(enemies)

	return c
}

func (c *CombatState) Enter() {
}

func (c *CombatState) Exit() {
}

func (c *CombatState) Update(dt float64) bool {
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

	camera := pixel.IM.Scaled(c.Pos, 1.0).Moved(c.win.Bounds().Center().Sub(c.Pos))
	c.win.SetMatrix(camera)
}

func (c *CombatState) HandleInput(win *pixelgl.Window) {
}

func (c *CombatState) CreateCombatCharacters(key string) {
	actorsList := c.Actors[key]
	layout := c.Layout[key][len(actorsList)-1]

	var charactersList []*Character
	for k, v := range actorsList {
		charDef := CharacterDefinitions[v.Id]

		if charDef.CombatEntityDef.Texture != nil {
			charDef.EntityDef = charDef.CombatEntityDef
		}

		var char *Character
		charStates := make(map[string]func() state_machine.State)
		for k, v := range charDef.CombatStates {
			charStates[k] = func() state_machine.State {
				return v(char, c)
			}
		}
		char = CharacterCreate(
			charDef,
			charStates,
		)

		charactersList = append(charactersList, char)
		c.ActorCharMap[v.Id] = char

		pos := layout[k]

		// Combat positions are 0 - 1
		// Need scaling to the screen size.
		x := pos.X * c.win.Bounds().W()
		y := pos.Y * c.win.Bounds().H()
		//char.Entity.Sprite:SetPosition(x, y)
		char.Entity.X = x
		char.Entity.Y = y

		// Change to standby because it's combat time
		char.Controller.Change(charDef.DefaultCombatState, nil)

	}

	c.Characters[key] = charactersList

}
