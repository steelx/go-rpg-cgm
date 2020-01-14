package storyboard

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
	"reflect"
)

type Storyboard struct {
	Stack         *gui.StateStack
	InternalStack *gui.StateStack
	States        map[string]gui.StackInterface
	Events        []interface{} //always keep as last args
}

func Create(stack *gui.StateStack, win *pixelgl.Window, eventsI interface{}) *Storyboard {
	sb := &Storyboard{
		Stack:         stack,
		InternalStack: gui.StateStackCreate(win),
		States:        make(map[string]gui.StackInterface),
	}

	if eventsI != nil {
		events := reflect.ValueOf(eventsI)
		if events.Len() > 0 {
			sb.Events = make([]interface{}, events.Len())
			for i := 0; i < events.Len(); i++ {
				sb.Events[i] = events.Index(i).Interface()
			}
		}
	}

	return sb
}

func (s Storyboard) CleanUp() {

}

func (s *Storyboard) PushState(identifier string, state gui.StackInterface) {
	//push a State on the stack but keep a reference here
	if _, ok := s.States[identifier]; ok {
		//found already
		return
	}
	fmt.Println("Adding", identifier)
	s.States[identifier] = state
	s.InternalStack.Push(state)
}

func (s *Storyboard) RemoveState(identifier string) {
	stateV := s.States[identifier]
	delete(s.States, identifier)
	for i, v := range s.InternalStack.States {
		if reflect.DeepEqual(v, stateV) {
			fmt.Println("removing", identifier)
			s.removeSliceItem(i)
			break
		}
	}
}

func (s *Storyboard) removeSliceItem(i int) {
	s.InternalStack.States[i] = s.InternalStack.States[len(s.InternalStack.States)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	s.InternalStack.States = s.InternalStack.States[:len(s.InternalStack.States)-1]
}

/*
	StateStack interface implemented below
*/
func (s Storyboard) Enter() {

}

func (s Storyboard) Exit() {

}

func (s *Storyboard) Update(dt float64) bool {
	s.InternalStack.Update(dt)

	deleteIndex := -1
Loop:
	for k, v := range s.Events {

		switch x := v.(type) {
		case func(storyboard *Storyboard):
			x(s)

		case *WaitEvent:
			s.Events[k] = x
		case *NonBlockEvent:
			s.Events[k] = x
		case *TweenEvent:
			s.Events[k] = x
		case *BlockUntilEvent:
			s.Events[k] = x
		case *TimedTextboxEvent:
			s.Events[k] = x

		case func(storyboard *Storyboard) *WaitEvent:
			s.Events[k] = x(s)

		case func(storyboard *Storyboard) *NonBlockEvent:
			s.Events[k] = x(s)

		case func(storyboard *Storyboard) *TweenEvent:
			s.Events[k] = x(s)

		case func(storyboard *Storyboard) *BlockUntilEvent:
			s.Events[k] = x(s)

		case func(storyboard *Storyboard) *TimedTextboxEvent:
			s.Events[k] = x(s)

		default:
			fmt.Printf("Unsupported type: %T\n", x)
		}

		valV := reflect.ValueOf(s.Events[k])
		valI := valV.Interface().(SBEvent)
		s.Events[k] = valI
		valI.Update(dt)
		if valI.IsFinished() {
			deleteIndex = k
			break Loop
		}
		if valI.IsBlocking() {
			break Loop
		}

	}
	//Loop END

	if deleteIndex != -1 {
		s.Events[deleteIndex], s.Events[0] = s.Events[0], s.Events[deleteIndex]
		s.Events = s.Events[1:]
	}

	if len(s.Events) == 0 {
		s.Stack.Pop()
		return true
	}

	return true
}

func (s Storyboard) Render(win *pixelgl.Window) {
	//debugText := fmt.Sprintf("Storyboard Events # %v", len(s.Events))
	//fmt.Println(debugText)

	s.InternalStack.Render(win)
}

func (s Storyboard) HandleInput(win *pixelgl.Window) {

}
