package main

type UV struct {
	ux, uy, vx, vy float64
}

type State interface {
	Enter(data Direction)
	Render()
	Exit()
	Update(dt float64)
}

type Character struct {
	mAnimUp     []int
	mAnimRight  []int
	mAnimDown   []int
	mAnimLeft   []int
	mEntity     *Entity
	mController *StateMachine //[name] -> [function that returns state]
}

type Direction struct {
	x, y float64
}
