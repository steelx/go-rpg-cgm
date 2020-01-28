package game_map

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"github.com/steelx/go-rpg-cgm/state_machine"
	"reflect"
)

//igm.StateMachine
const (
	status int = iota
	items
	equip
)

var frontMenuOrder = []string{
	"Status",
	"Items",
	"Equipment",
}

//parent
type InGameMenuState struct {
	Stack        *gui.StateStack
	StateMachine *state_machine.StateMachine
	World        *combat.WorldExtended
}

func InGameMenuStateCreate(stack *gui.StateStack, win *pixelgl.Window) *InGameMenuState {
	worldV := reflect.ValueOf(stack.Globals["world"])
	igm := &InGameMenuState{
		Stack: stack,
		World: worldV.Interface().(*combat.WorldExtended),
	}

	igm.StateMachine = state_machine.Create(map[string]func() state_machine.State{
		"frontmenu": func() state_machine.State {
			return FrontMenuStateCreate(igm, win)
		},
		frontMenuOrder[items]: func() state_machine.State {
			return ItemsMenuStateCreate(igm, win)
		},
		"magic": func() state_machine.State {
			//return MagicMenuStateCreate(this)
			return state_machine.Create(map[string]func() state_machine.State{})
		},
		frontMenuOrder[equip]: func() state_machine.State {
			return EquipMenuStateCreate(igm, win)
		},
		frontMenuOrder[status]: func() state_machine.State {
			//return StatusMenuStateCreate(this)
			return StatusMenuStateCreate(igm, win)
		},
	})

	igm.StateMachine.Change("frontmenu", nil)

	return igm
}

func (igm *InGameMenuState) Update(dt float64) bool {
	igm.StateMachine.Update(dt)
	//fmt.Println("ingame_menu_state", reflect.DeepEqual(igm.Stack.Top(), igm)) // temp
	//if reflect.DeepEqual(igm.Stack.Top(), igm) {
	//	igm.StateMachine.Update(dt)
	//}
	return true
}
func (igm InGameMenuState) Render(win *pixelgl.Window) {
	igm.StateMachine.Render(win)

	//temp camera matrix
	cam := pixel.IM.Scaled(win.Bounds().Center(), 1.0).Moved(win.Bounds().Center())
	win.SetMatrix(cam)
}

func (igm InGameMenuState) Enter()                          {}
func (igm InGameMenuState) Exit()                           {}
func (igm InGameMenuState) HandleInput(win *pixelgl.Window) {}
