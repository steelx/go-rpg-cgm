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

type Direction struct {
	x, y float64
}
