package storyboard

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/steelx/go-rpg-cgm/gui"
	"reflect"
)

type Storyboard struct {
	Stack  *gui.StateStack
	Events []interface{} //always keep as last args
}

func Create(stack *gui.StateStack, eventsI interface{}) Storyboard {
	sb := Storyboard{
		Stack: stack,
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

/*
	StateStack interface implemented below
*/
func (s Storyboard) Enter() {

}

func (s Storyboard) Exit() {

}

func (s *Storyboard) Update(dt float64) bool {
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

		default:
			fmt.Printf("Unsupported type: %T\n", x)
		}

	}
	//Loop END

	if deleteIndex != -1 {
		s.Events[deleteIndex], s.Events[0] = s.Events[0], s.Events[deleteIndex]
		s.Events = s.Events[1:]
	}
	//fmt.Println("AFTER", s.Events)

	return true
}

func (s Storyboard) Render(win *pixelgl.Window) {

}

func (s Storyboard) HandleInput(win *pixelgl.Window) {

}
