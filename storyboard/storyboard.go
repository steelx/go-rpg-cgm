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
	States        map[string]*gui.StackInterface
	Events        []interface{} //always keep as last args
}

func Create(stack *gui.StateStack, win *pixelgl.Window, eventsI interface{}) Storyboard {
	sb := Storyboard{
		Stack:         stack,
		InternalStack: gui.StateStackCreate(win),
		States:        make(map[string]*gui.StackInterface),
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
	s.States[identifier] = &state //identifier e.g. blackscreen
	s.InternalStack.Push(state)
}

//func (s *Storyboard) RemoveState(identifier string) {
//	stateV := s.States[identifier]
//	delete(s.States, identifier)
//	for _, v := range s.InternalStack.States {
//		if reflect.DeepEqual(v, stateV) {
//			fmt.Println("found")
//		}
//	}
//}

/*
	StateStack interface implemented below
*/
func (s Storyboard) Enter() {

}

func (s Storyboard) Exit() {

}

func (s *Storyboard) Update(dt float64) bool {
	s.InternalStack.Update(dt)

	if len(s.Events) == 0 {
		s.Stack.Pop()
	}
	deleteIndex := -1
Loop:
	for k, v := range s.Events {

		switch x := v.(type) {
		case *WaitEvent:
			x.Update(dt)
			if x.IsFinished() {
				deleteIndex = k
				break Loop
			}
			if x.IsBlocking() {
				break Loop
			}

		case func(storyboard *Storyboard) *WaitEvent:
			xv := x(s)
			xv.Update(dt)
			if xv.IsFinished() {
				deleteIndex = k
				break Loop
			}
			if xv.IsBlocking() {
				break Loop
			}

		case func(storyboard *Storyboard, dt float64) TweenEvent:
			xv := x(s, dt)
			xv.Update(dt)

			if xv.IsFinished() {
				deleteIndex = k
				break Loop
			}
			if xv.IsBlocking() {
				break Loop
			}

		default:
			fmt.Printf("Unsupported type: %T\n", x)
		}

	}
	//Loop END

	if deleteIndex != -1 {
		s.Events[deleteIndex], s.Events[0] = s.Events[0], s.Events[deleteIndex]
		s.Events = s.Events[1:]
	}

	return true
}

func (s Storyboard) Render(win *pixelgl.Window) {
	debugText := fmt.Sprintf("Storyboard Events # %v", len(s.Events))
	fmt.Println(debugText)

	s.InternalStack.Render(win)
}

func (s Storyboard) HandleInput(win *pixelgl.Window) {

}
