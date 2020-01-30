package combat

import (
	"fmt"
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
		q.Queue = append([]Event{eventI}, q.Queue...)
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

func (q *EventQueue) removeAtIndex(i int) {
	// Remove the element at index i from a. option 1
	/*
		copy(q.Queue[i:], q.Queue[i+1:])// Shift a[i+1:] left one index.
		q.Queue[len(q.Queue)-1] = nil// Erase last element (write zero value).
		q.Queue = q.Queue[:len(q.Queue)-1]// Truncate slice.
	*/
	//option 2
	q.Queue = append(q.Queue[:i], q.Queue[i+1:]...)
}

func (q *EventQueue) Clear() {
	q.Queue = make([]Event, 0)
	q.CurrentEvent = nil
}
func (q EventQueue) IsEmpty() bool {
	return len(q.Queue) == 0
}

func (q EventQueue) ActorHasEvent(actor *Actor) bool {
	if q.CurrentEvent != nil {
		a := q.CurrentEvent.Owner()
		if &a == &actor {
			return true
		}
	}

	for _, v := range q.Queue {
		a := v.Owner()
		if &a == &actor {
			return true
		}
	}

	return false
}

func (q *EventQueue) RemoveEventsOwnedBy(actor *Actor) {
	for i := len(q.Queue) - 1; i >= 0; i-- {
		v := q.Queue[i]
		if actor == v.Owner() {
			q.removeAtIndex(i)
			//q.Clear()
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
		fmt.Println("Event Queue is empty.")
		return
	}

	fmt.Println("Event Queue:")
	if q.CurrentEvent != nil {
		fmt.Println("Current event:", q.CurrentEvent.Name)
	}

	for k, v := range q.Queue {
		msg := fmt.Sprintf("[%d] Event: [%v][%s]", k, v.CountDown, v.Name)
		fmt.Println(msg)
	}
}

func (q *EventQueue) Update() {
	if q.CurrentEvent != nil {
		q.CurrentEvent.Update()

		if q.CurrentEvent.IsFinished() {
			q.CurrentEvent = nil
			//once finished we go to next event from here
		} else {
			return //Only one event is executed at a time
		}
	} else if q.IsEmpty() {
		return
	} else {
		// Need to chose an event
		front := q.Queue[0]
		q.removeAtIndex(0)
		front.Execute(q)
		q.CurrentEvent = front
	}

	//all the other events countdown reduced by one.
	for _, v := range q.Queue {
		//ensure countdown doesnt drop below 0
		v.CountDownSet(math.Max(0, v.CountDown()-1))
	}
}
