package combat

import "fmt"

type Scene struct {
	PartyActors, EnemyActors []*Actor
	EventQueue               *EventQueue
}

func SceneCreate(partyMembers []*Actor, enemies []*Actor) *Scene {
	s := &Scene{
		PartyActors: partyMembers,
		EnemyActors: enemies,
		EventQueue:  EventsQueueCreate(),
	}
	s.AddTurns(s.EnemyActors)
	s.AddTurns(s.PartyActors)
	return s
}

func (s Scene) Update() {
	s.EventQueue.Update()

	if s.IsPartyDefeated() || s.IsEnemyDefeated() {
		//END GAME DETECTED
		//s.EventQueue.Clear() could be used
		s.EventQueue.Queue = make([]Event, 0)
		return
	}

	//keep the queue pumping
	s.AddTurns(s.PartyActors)
	s.AddTurns(s.EnemyActors)
}

func (s *Scene) AddTurns(actors []*Actor) {
	for _, v := range actors {
		if !s.EventQueue.ActorHasEvent(v) {
			event := CETurnCreate(s, v)
			tp := event.TimePoints(*s.EventQueue)
			s.EventQueue.Add(event, tp)
		}
	}
}

func (s *Scene) GetTarget(owner *Actor) *Actor {
	if owner.IsPlayer() {
		return s.EnemyActors[len(s.EnemyActors)-1]
	}

	return s.PartyActors[len(s.EnemyActors)-1]
}

func (s Scene) GetAlivePartyActors() []*Actor {
	var alive []*Actor
	for _, a := range s.PartyActors {
		if !a.IsKOed() {
			alive = append(alive, a)
		}
	}
	return alive
}

//OnDead makes Actor KnockOut
func (s *Scene) OnDead(actor *Actor) {
	if actor.IsPlayer() {
		actor.KO()
	} else {
		for i := len(s.EnemyActors) - 1; i >= 0; i-- {
			if actor == s.EnemyActors[i] {
				s.EnemyActors = s.removeAtIndex(s.EnemyActors, i)
			}
		}
	}

	//Remove owned events
	s.EventQueue.RemoveEventsOwnedBy(actor)

	if s.IsPartyDefeated() {
		fmt.Println("Scene OnDead: Party loses")
	} else if s.IsEnemyDefeated() {
		fmt.Println("Scene OnDead: Enemy loses")
	}
}

func (s Scene) removeAtIndex(arr []*Actor, i int) []*Actor {
	return append(arr[:i], arr[i+1:]...)
}

//IsPartyDefeated check's at least 1 Actor is standing return false
func (s Scene) IsPartyDefeated() bool {
	for _, actor := range s.PartyActors {
		if !actor.IsKOed() {
			return false
		}
	}
	return true
}

func (s Scene) IsEnemyDefeated() bool {
	return len(s.EnemyActors) == 0
}
