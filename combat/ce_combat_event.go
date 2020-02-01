package combat

//CombatEvent (CE)
type Event interface {
	Name() string
	CountDown() float64
	CountDownSet(t float64)
	Owner() *Actor
	Update()
	IsFinished() bool
	Execute(queue *EventQueue)
}
