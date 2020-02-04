package game_map

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/sirupsen/logrus"
	"github.com/steelx/go-rpg-cgm/combat"
	"github.com/steelx/go-rpg-cgm/gui"
	"math"
)

type EventQueue struct {
	Queue        []Event
	CurrentEvent Event
}

func EventsQueueCreate() *EventQueue {

	return &EventQueue{
		Queue:        make([]Event, 0),
		CurrentEvent: nil,
	}
}

func (q *EventQueue) Add(eventI Event, timePoints float64) {
	//Instant event
	eventI.CountDownSet(timePoints)
	if timePoints == -1 {
		//push the event to top
		q.insertAtIndex(0, eventI)
		return
	}

	for i := 0; i < len(q.Queue); i++ {
		count := q.Queue[i].CountDown()
		if count > eventI.CountDown() {
			q.insertAtIndex(i, eventI)
			return
		}
	}

	q.Queue = append(q.Queue, eventI)
}

func (q *EventQueue) insertAtIndex(index int, eventI Event) {
	temp := append([]Event{}, q.Queue[index:]...)
	q.Queue = append(q.Queue[0:index], eventI)
	q.Queue = append(q.Queue, temp...)
}

func (q *EventQueue) removeQueAtIndex(i int) {
	// Remove the element at index i from Queue
	q.Queue = append(q.Queue[:i], q.Queue[i+1:]...)
}

func (q *EventQueue) Clear() {
	q.Queue = make([]Event, 0)
	q.CurrentEvent = nil
}
func (q EventQueue) IsEmpty() bool {
	return len(q.Queue) == 0
}

func (q EventQueue) ActorHasEvent(actor *combat.Actor) bool {
	if q.CurrentEvent != nil && q.CurrentEvent.Owner() == actor {
		return true
	}

	for _, v := range q.Queue {
		if v.Owner() == actor {
			return true
		}
	}

	return false
}

func (q *EventQueue) RemoveEventsOwnedBy(actor *combat.Actor) {
	for i := len(q.Queue) - 1; i >= 0; i-- {
		v := q.Queue[i]
		if actor == v.Owner() {
			q.removeQueAtIndex(i)
		}
	}
}

func (q EventQueue) SpeedToTimePoints(speed float64) float64 {
	maxSpeed := 255.0
	speed = math.Min(speed, 255)
	points := maxSpeed - speed
	return math.Floor(points)
}

//Print just for debug
func (q EventQueue) Print() {
	if q.IsEmpty() {
		logrus.Info("Event Queue is empty.")
		return
	}

	logrus.Info("Event Queue:")
	if q.CurrentEvent != nil {
		logrus.Info("Current event:", q.CurrentEvent.Name())
	}

	for k, v := range q.Queue {
		msg := fmt.Sprintf("[%d] Event: [%v][%s]", k, v.CountDown(), v.Name())
		logrus.Info(msg)
	}
}

func (q *EventQueue) Update() {
	if q.CurrentEvent != nil {
		q.CurrentEvent.Update()

		if !q.CurrentEvent.IsFinished() {
			return //Only one event is executed at a time
		}
		q.CurrentEvent = nil
		//once finished we go to update countdown
		// which helps in going to next event
	} else if q.IsEmpty() {
		return
	} else {
		// Need to chose an event
		front := q.Queue[0]
		q.removeQueAtIndex(0)
		front.Execute(q)
		q.CurrentEvent = front
	}

	//all the other events countdown reduced by one.
	for _, v := range q.Queue {
		//ensure countdown doesnt drop below 0
		v.CountDownSet(math.Max(0, v.CountDown()-1))
	}
}

func (q *EventQueue) Render(win *pixelgl.Window) {
	yInc := 15.5
	var width, height float64
	if win.Monitor() != nil {
		width, height = win.Monitor().Size()
	} else {
		width, height = win.Bounds().W(), win.Bounds().H()
	}
	x := -width / 2
	y := height / 2

	textBase := text.New(pixel.V(0, 0), gui.BasicAtlasAscii)
	if q.CurrentEvent != nil {
		fmt.Fprintln(textBase, fmt.Sprintf("CURRENT: %s", q.CurrentEvent.Name()))
	}

	y = y - yInc

	if q.IsEmpty() {
		fmt.Fprintln(textBase, "EMPTY !")
	}

	for k, v := range q.Queue {
		out := fmt.Sprintf("[%d] Event: [%v][%v]", k, v.CountDown(), v.Name())
		fmt.Fprintln(textBase, out)
		y = y - yInc
	}

	textBase.Draw(win, pixel.IM.Moved(pixel.V(x, y)))
}
