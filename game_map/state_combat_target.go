package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
)

type TargetType int

const (
	CombatTargetTypeONE TargetType = iota
	CombatTargetTypeSIDE
	CombatTargetTypeALL
)

type CombatSelectorFunc struct {
	WeakestEnemy,
	SideEnemy,
	SelectAll func(state *CombatState) []*combat.Actor
}

var CombatSelector = CombatSelectorFunc{
	WeakestEnemy: func(state *CombatState) []*combat.Actor {
		enemyList := state.Actors[enemies]
		health := 99999.9

		var target *combat.Actor
		for _, v := range enemyList {
			hpNow := v.Stats.Get("HpNow")
			if hpNow < health {
				health = hpNow
				target = v
			}
		}
		return []*combat.Actor{target}
	},
	SideEnemy: func(state *CombatState) []*combat.Actor {
		return state.Actors[enemies]
	},
	SelectAll: func(state *CombatState) []*combat.Actor {
		return append(state.Actors[enemies], state.Actors[party]...)
	},
}

type CombatTargetState struct {
	CombatState     *CombatState
	Stack           *gui.StateStack                          //The internal stack of states from the CombatState object.
	DefaultSelector func(state *CombatState) []*combat.Actor //The function that chooses which characters are targeted
	//when the state begins.
	CanSwitchSide bool
	SelectType    TargetType
	OnSelect      func(targets []*combat.Actor)
	OnExit        func()
	Targets,
	Enemies, Party []*combat.Actor
	MarkerPNG      pixel.Picture
	Marker         *pixel.Sprite
	MarkerPosition pixel.Vec
}

type CombatChoiceParams struct {
	OnSelect        func(targets []*combat.Actor)
	OnExit          func()
	SwitchSides     bool
	DefaultSelector func(state *CombatState) []*combat.Actor
	TargetType      TargetType
}

func CombatTargetStateCreate(state *CombatState, choiceParams CombatChoiceParams) *CombatTargetState {
	t := &CombatTargetState{
		CombatState:     state,
		Stack:           state.InternalStack,
		DefaultSelector: choiceParams.DefaultSelector,
		CanSwitchSide:   choiceParams.SwitchSides,
		SelectType:      choiceParams.TargetType,
		OnSelect:        choiceParams.OnSelect,
		OnExit:          choiceParams.OnExit,
		MarkerPNG:       gui.CursorPng,
		Marker:          pixel.NewSprite(gui.CursorPng, gui.CursorPng.Bounds()),
	}

	if t.DefaultSelector == nil {
		if t.SelectType == CombatTargetTypeONE {
			t.DefaultSelector = CombatSelector.WeakestEnemy
		} else if t.SelectType == CombatTargetTypeSIDE {
			t.DefaultSelector = CombatSelector.SideEnemy
		} else if t.SelectType == CombatTargetTypeALL {
			t.DefaultSelector = CombatSelector.SelectAll
		}
	}

	return t
}

func (t *CombatTargetState) Enter() {
	t.Enemies = t.CombatState.Actors[enemies]
	t.Party = t.CombatState.Actors[party]
	t.Targets = t.DefaultSelector(t.CombatState)
}

func (t *CombatTargetState) Exit() {
	t.Enemies = nil
	t.Party = nil
	t.Targets = nil
	t.OnExit()
}

func (t *CombatTargetState) Update(dt float64) bool {
	return true
}

func (t *CombatTargetState) Render(renderer *pixelgl.Window) {

	for _, v := range t.Targets {
		char := t.CombatState.ActorCharMap[v]
		pos := char.Entity.GetTargetPosition()
		pos = pos.Add(pixel.V(0, t.MarkerPNG.Bounds().W()/2))
		t.MarkerPosition = pos
		t.Marker.Draw(renderer, pixel.IM.Moved(pos))
	}

}

func (t *CombatTargetState) HandleInput(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyBackspace) {
		t.Stack.Pop()
	} else if win.JustPressed(pixelgl.KeyUp) {
		t.Up()
	} else if win.JustPressed(pixelgl.KeyDown) {
		t.Down()
	} else if win.JustPressed(pixelgl.KeyLeft) {
		t.Left()
	} else if win.JustPressed(pixelgl.KeyRight) {
		t.Right()
	} else if win.JustPressed(pixelgl.KeySpace) {
		t.OnSelect(t.Targets)
	}
}

/////////////////////////////////
// CombatTargetState additional methods below
/////////////////////////////////

func (t CombatTargetState) GetActorList(actor *combat.Actor) []*combat.Actor {
	if isParty := t.CombatState.IsPartyMember(actor); isParty {
		return t.Party
	}

	return t.Enemies
}

//GetIndex finds Actors based on uniq names, since we might copy paste Enemies of same entity
//e.g. combat.ActorCreate(enemyDef, "1")
func (t CombatTargetState) GetIndex(listI []*combat.Actor, item *combat.Actor) int {
	for k, v := range listI {
		if v.Name == item.Name {
			return k
		}
	}
	return 0
}

func (t *CombatTargetState) Left() {
	if !t.CanSwitchSide || !t.CombatState.IsPartyMember(t.Targets[0]) {
		return
	}
	if t.SelectType == CombatTargetTypeONE {
		t.Targets = []*combat.Actor{t.Enemies[0]}
	}
	if t.SelectType == CombatTargetTypeSIDE {
		t.Targets = t.Enemies
	}
}

func (t *CombatTargetState) Right() {
	if !t.CanSwitchSide || !t.CombatState.IsPartyMember(t.Targets[0]) {
		return
	}
	if t.SelectType == CombatTargetTypeONE {
		t.Targets = []*combat.Actor{t.Party[0]}
	}
	if t.SelectType == CombatTargetTypeSIDE {
		t.Targets = t.Party
	}
}

func (t *CombatTargetState) Up() {
	if t.SelectType != CombatTargetTypeONE {
		return
	}

	selected := t.Targets[0]
	side := t.GetActorList(selected)
	index := t.GetIndex(side, selected)

	index = index + 1
	if index >= len(side) {
		index = 0
	}
	t.Targets = []*combat.Actor{side[index]}
}

func (t *CombatTargetState) Down() {
	if t.SelectType != CombatTargetTypeONE {
		return
	}

	selected := t.Targets[0]
	side := t.GetActorList(selected)
	index := t.GetIndex(side, selected)

	index = index - 1
	if index == -1 {
		index = len(side) - 1
	}
	t.Targets = []*combat.Actor{side[index]}
}
