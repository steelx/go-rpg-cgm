package combat

import "fmt"

//CombatEventTurn
type CETurn struct {
	Scene     *Scene
	owner     *Actor
	name      string
	countDown float64
}

func CETurnCreate(scene *Scene, owner *Actor) *CETurn {
	return &CETurn{
		Scene: scene,
		owner: owner,
		name:  fmt.Sprintf("CETurn(_, %s)", owner.Name),
	}
}

func (c CETurn) Name() string {
	return c.name
}

func (c CETurn) CountDown() float64 {
	return c.countDown
}

func (c *CETurn) CountDownSet(t float64) {
	c.countDown = t
}

func (c CETurn) Owner() *Actor {
	return c.owner
}

func (c *CETurn) Update() {

}

func (c CETurn) IsFinished() bool {
	return true
}

func (c *CETurn) Execute(queue *EventQueue) {
	target := c.Scene.GetTarget(c.owner)
	msg := fmt.Sprintf("%s decides to attack %s", c.owner.Name, target.Name)
	fmt.Println(msg)

	event := CEAttackCreate(c.Scene, c.owner, target)
	tp := event.TimePoints(*queue)
	queue.Add(event, tp)
}

func (c CETurn) TimePoints(queue *EventQueue) float64 {
	speed := c.Owner().Stats.Get("Speed")
	return queue.SpeedToTimePoints(speed)
}
